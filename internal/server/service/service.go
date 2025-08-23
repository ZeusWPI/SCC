package service

import (
	"context"

	"github.com/zeusWPI/scc/internal/database/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) withRollback(ctx context.Context, fn func(context.Context) error) error {
	return s.repo.WithRollback(ctx, fn)
}
