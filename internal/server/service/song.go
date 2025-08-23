package service

import (
	"context"

	"github.com/zeusWPI/scc/internal/server/dto"
)

type Song struct{}

func (s *Service) NewSong() *Song {
	return &Song{}
}

// TODO: Fill in
func (s *Song) New(ctx context.Context, song dto.Song) error {
	return nil
}
