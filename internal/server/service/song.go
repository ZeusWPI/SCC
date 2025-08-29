package service

import (
	"context"

	"github.com/zeusWPI/scc/internal/server/dto"
)

type Song struct{}

func (s *Service) NewSong() *Song {
	return &Song{}
}

func (s *Song) New(_ context.Context, _ dto.Song) error {
	// TODO: Fill in
	return nil
}
