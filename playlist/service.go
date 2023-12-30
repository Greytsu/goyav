package playlist

import (
	"context"
	"errors"
	"slices"

	"github.com/rs/zerolog/log"
	spotify2 "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"

	"goyav/spotify"
)

type Service struct {
	repository  repository
	spotifyServ *spotify.Service
}

var (
	InvalidAttributesError  = errors.New("invalid attributes")
	UserNotContributorError = errors.New("user not allowed to view or modify this playlist, ask owner to give permission")
	UserNotOwnerError       = errors.New("only the owner of this playlist is allowed to modify contributors")
)

type repository interface {
	GetPlaylists(ctx context.Context) (*[]Playlist, error)
	GetPlaylist(ctx context.Context, playlistID string) (*Playlist, error)
	CreatePlaylist(ctx context.Context, playlist *Playlist) (*Playlist, error)
	UpdatePlaylistTracks(ctx context.Context, playlist *Playlist) error
	UpdatePlaylistContributors(ctx context.Context, playlist *Playlist) error
}

func NewService(repo repository, spotifyService *spotify.Service) *Service {
	return &Service{
		repository:  repo,
		spotifyServ: spotifyService,
	}
}

func (s *Service) GetPlaylists(ctx context.Context) (*[]Playlist, error) {
	return s.repository.GetPlaylists(ctx)
}

func (s *Service) GetPlaylist(ctx context.Context, playlistID, userID string) (*Playlist, error) {
	foundPlaylist, err := s.repository.GetPlaylist(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	// returns an error if user is not a contributor
	if !slices.Contains(foundPlaylist.Contributors, userID) && foundPlaylist.Owner != userID {
		return nil, UserNotContributorError
	}

	return foundPlaylist, nil
}

func (s *Service) CreatePlaylist(ctx context.Context, newPlaylist *Playlist, tok *oauth2.Token) (*Playlist, error) {
	// Track attributes check
	err := checkAttributes(newPlaylist)
	if err != nil {
		return nil, err
	}

	// Setup playlist
	setupPlaylist(newPlaylist)

	// Create spotify playlist
	spotifyPlaylistID, err := s.spotifyServ.CreateSpotifyPlayList(tok, newPlaylist.Owner, newPlaylist.Name)
	if err != nil {
		log.Info().Err(err).Msg("Error while creating spotify playlist")
		return nil, err
	}
	newPlaylist.SpotifyID = spotify2.ID(spotifyPlaylistID)

	return s.repository.CreatePlaylist(ctx, newPlaylist)
}

func (s *Service) AddRecommendationsToPlaylist(ctx context.Context, tok *oauth2.Token, playlistID, userID string, recommendationRequest *spotify.RecommendationRequest) error {
	playlist, err := s.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return err
	}

	recommendations, err := s.spotifyServ.GetRecommendations(tok, recommendationRequest)
	if err != nil {
		return err
	}

	// insert recommendations in playlist
	for _, recommendation := range recommendations {
		trackIdStr := string(recommendation.TrackId)
		if playlist.Tracks[trackIdStr].TrackId == "" {
			// add track to playlist
			recommendation.Contributors = []string{userID}
			playlist.Tracks[trackIdStr] = recommendation
		} else {
			// add user to track contributors
			if !slices.Contains(playlist.Tracks[trackIdStr].Contributors, userID) {
				playlist.AddContributorToTrack(trackIdStr, userID)
			}
		}
	}

	// update playlist
	err = s.repository.UpdatePlaylistTracks(ctx, playlist)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateContributors(ctx context.Context, playlistID, userID string, contributors []string) error {
	playlist, err := s.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return err
	}

	// check if user is owner of playlist
	if playlist.Owner != userID {
		return UserNotOwnerError
	}

	// update contributors
	playlist.Contributors = contributors
	err = s.repository.UpdatePlaylistContributors(ctx, playlist)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateSpotifyPlaylistTracks(ctx context.Context, tok *oauth2.Token, playlistID, userID string) error {
	playlist, err := s.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return err
	}

	err = s.spotifyServ.UpdatePlaylistTracks(ctx, tok, playlist.SpotifyID, playlist.Tracks)
	if err != nil {
		return err
	}

	return nil
}
