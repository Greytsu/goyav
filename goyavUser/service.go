package goyavUser

import "context"

type Service struct {
	repository repository
}

type repository interface {
	CreateUser(ctx context.Context, user User) error
}

func NewService(repo repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) CreateUser(ctx context.Context, user User) error {
	return s.repository.CreateUser(ctx, user)
}
