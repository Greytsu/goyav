package playlist

import "context"

type Service struct {
	repository repository
}

type repository interface {
	GetPlaylists(ctx context.Context) (*[]Playlist, error)
	GetPlaylist(ctx context.Context, id string) (*Playlist, error)
	CreatePlaylist(ctx context.Context, playlist *Playlist) (*Playlist, error)
}

func NewService(repo repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) GetPlaylists(ctx context.Context) (*[]Playlist, error) {
	return s.repository.GetPlaylists(ctx)
}

func (s *Service) GetPlaylist(ctx context.Context, id string) (*Playlist, error) {
	return s.repository.GetPlaylist(ctx, id)
}

func (s *Service) CreatePlaylist(ctx context.Context, playlist *Playlist) (*Playlist, error) {
	return s.repository.CreatePlaylist(ctx, playlist)
}
