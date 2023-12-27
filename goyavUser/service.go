package goyavUser

import "context"

type Service struct {
	repository repository
}

type repository interface {
	GetUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user User) error
}

func NewService(repo repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repository.GetUser(ctx, id)
}

func (s *Service) CreateUser(ctx context.Context, user User) error {
	return s.repository.CreateUser(ctx, user)
}
