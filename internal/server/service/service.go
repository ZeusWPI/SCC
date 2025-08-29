package service

import (
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
