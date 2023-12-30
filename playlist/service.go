package playlist

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"goyav/spotify"
	"goyav/track"
)

type Service struct {
	repository  repository
	spotifyServ *spotify.Service
}

var (
	InvalidAttributesError = errors.New("invalid attributes")
)

type repository interface {
	GetPlaylists(ctx context.Context) (*[]Playlist, error)
	GetPlaylist(ctx context.Context, id string) (*Playlist, error)
	CreatePlaylist(ctx context.Context, playlist *Playlist) (*Playlist, error)
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

func (s *Service) GetPlaylist(ctx context.Context, id string) (*Playlist, error) {
	return s.repository.GetPlaylist(ctx, id)
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
	newPlaylist.SpotifyID = spotifyPlaylistID

	return s.repository.CreatePlaylist(ctx, newPlaylist)
}

func checkAttributes(newPlaylist *Playlist) error {
	for _, attribute := range newPlaylist.TrackAttributes {
		res, err := IsValidAttributeValue(attribute.AttributeType, attribute.AttributeValue)
		if err != nil {
			log.Info().Err(err).Msg("Error while checking attributes")
			return err
		}
		if !res {
			return InvalidAttributesError
		}
	}
	return nil
}

func setupPlaylist(newPlaylist *Playlist) {
	// Setup playlist contributors
	if newPlaylist.Contributors == nil || len(newPlaylist.Contributors) == 0 {
		newPlaylist.Contributors = []string{newPlaylist.Owner}
	}

	// Setup playlist tracks
	if newPlaylist.Tracks == nil {
		newPlaylist.Tracks = map[string]track.Track{}
	}
}
