package spotify

import (
	"context"
	"golang.org/x/oauth2"

	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"goyav/artist"
	"goyav/goyavUser"
	"goyav/track"
)

type Service struct {
	auth *spotifyAuth.Authenticator
}

func NewService(auth *spotifyAuth.Authenticator) *Service {
	return &Service{auth: auth}
}

func (s *Service) CheckUser(tok *oauth2.Token) bool {
	// new spotify client
	client := spotify.New(s.auth.Client(context.Background(), tok))

	// retrieve current user info
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser")
		return false
	}
	if user.ID == "" {
		return false
	}
	return true
}

func (s *Service) GetCurrentUser(tok *oauth2.Token) (goyavUser.User, error) {
	// new spotify client
	client := spotify.New(s.auth.Client(context.Background(), tok))

	// retrieve current user info
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser")
		return goyavUser.User{}, err
	}
	goyavCurrentUser := goyavUser.MapSpotifyUserToGoyavUser(user)
	return goyavCurrentUser, nil
}

func (s *Service) GetCurrentUserTopArtists(tok *oauth2.Token) ([]artist.Artist, error) {
	// new spotify client
	client := spotify.New(s.auth.Client(context.Background(), tok))

	spotifyArtists, err := client.CurrentUsersTopArtists(context.Background())
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser top artists")
		return []artist.Artist{}, err
	}

	goyavArtists := artist.MapSpotifyArtistsToGoyavArtists(spotifyArtists.Artists)

	return goyavArtists, nil
}

func (s *Service) GetCurrentUserTopTracks(tok *oauth2.Token) (*spotify.FullTrackPage, error) {
	// new spotify client
	client := spotify.New(s.auth.Client(context.Background(), tok))

	tracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Timerange(spotify.ShortTermRange))
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser top tracks")
		return &spotify.FullTrackPage{}, err
	}

	return tracks, nil
}

func (s *Service) GetFullRecommendations(tok *oauth2.Token) ([]track.Track, error) {
	// new spotify client
	client := spotify.New(s.auth.Client(context.Background(), tok))

	// get user top artists
	spotifyArtists, err := client.CurrentUsersTopArtists(context.Background())
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser top artists")
		return nil, err
	}
	goyavArtists := artist.MapSpotifyArtistsToGoyavArtists(spotifyArtists.Artists)

	// get user top tracks
	spotifyTracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Timerange(spotify.ShortTermRange))
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser top tracks")
		return nil, err
	}
	goyavTracks := track.MapSpotifyTracksToGoyavTracks(spotifyTracks.Tracks)

	var recommendedTracks []track.Track

	// iterate by chunk since only 5 seeds can be used at a time
	tracksLenght := len(goyavArtists)
	for i := 0; i < tracksLenght; i += spotify.MaxNumberOfSeeds {
		end := i + spotify.MaxNumberOfSeeds
		if end > tracksLenght {
			end = tracksLenght
		}
		trackChunk := goyavTracks[i:end]

		recommendations, err := getTrackRecommendations(context.Background(), client, trackChunk)
		if err != nil {
			log.Info().Err(err).Msg("Error while retrieving recommendations")
			return nil, err
		}

		recommendedTracks = append(recommendedTracks, track.MapSpotifySimpleTracksToGoyavTracks(recommendations.Tracks)...)
	}

	//TODO: get recommendations from artists
	//TODO: remove duplicates from recommendations

	return recommendedTracks, nil
}

func getRecommendations(ctx context.Context, client *spotify.Client, artists []artist.Artist, tracks []track.Track) (*spotify.Recommendations, error) {
	seeds := spotify.Seeds{
		Artists: artist.MapToTrackSeeds(artists),
		Tracks:  track.MapToTrackSeeds(tracks),
		Genres:  artist.ExtractGenresFromArtists(artists),
	}

	log.Info().Interface("seeds", seeds).Msg("seeds")

	trackAttributes := spotify.NewTrackAttributes()
	trackAttributes.MinDuration(1000)

	spotify.Limit(100)

	return client.GetRecommendations(ctx, seeds, trackAttributes, spotify.Limit(100))
}

func getTrackRecommendations(ctx context.Context, client *spotify.Client, tracks []track.Track) (*spotify.Recommendations, error) {
	seeds := spotify.Seeds{
		Artists: []spotify.ID{},
		Tracks:  track.MapToTrackSeeds(tracks),
		Genres:  []string{},
	}

	log.Info().Interface("seeds", seeds).Msg("seeds")

	trackAttributes := spotify.NewTrackAttributes()
	trackAttributes.MinDuration(1000)

	spotify.Limit(100)

	return client.GetRecommendations(ctx, seeds, trackAttributes, spotify.Limit(100))
}
