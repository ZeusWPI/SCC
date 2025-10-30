package service

import (
	"context"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/server/dto"
	songClient "github.com/zeusWPI/scc/internal/song"
	"go.uber.org/zap"
)

type Song struct {
	song repository.Song
}

func (s *Service) NewSong() *Song {
	return &Song{
		song: *s.repo.NewSong(),
	}
}

func (s *Song) New(_ context.Context, songSave dto.Song) error {
	song := songSave.ToModel()

	// Run in the background as it can take some time
	go func(ctx context.Context, song *model.Song) {
		if err := s.save(ctx, song); err != nil {
			zap.S().Error(err)
			return
		}

		if err := s.saveHistory(ctx, *song); err != nil {
			zap.S().Error(err)
			return
		}
	}(context.Background(), song)

	return nil
}

func (s *Song) save(ctx context.Context, song *model.Song) error {
	songDB, err := s.song.GetBySpotify(ctx, song.SpotifyID)
	if err != nil {
		return err
	}
	if songDB != nil {
		// Song is already in the database
		*song = *songDB
		return nil
	}

	if err := songClient.C.Populate(song); err != nil {
		return err
	}

	if err := s.song.Create(ctx, song); err != nil {
		return err
	}

	return nil
}

func (s *Song) saveHistory(ctx context.Context, song model.Song) error {
	if err := s.song.CreateHistory(ctx, song); err != nil {
		return err
	}

	return nil
}
