package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/sqlc"
)

type Song struct {
	repo Repository
}

func (r *Repository) NewSong() *Song {
	return &Song{
		repo: *r,
	}
}

func (s *Song) GetBySpotifyID(ctx context.Context, id string) (*model.Song, error) {
	song, err := s.repo.queries(ctx).GetSongBySpotifyID(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get song by spotify id %s | %w", id, err)
		}
		return nil, nil
	}

	return model.SongModel(song), nil
}

func (s *Song) GetLast(ctx context.Context) (*model.Song, error) {
	song, err := s.repo.queries(ctx).GetLastSongFull(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get last song %w", err)
		}
		return nil, nil
	}

	return model.SongModelHistory(song), nil
}

func (s *Song) GetArtistBySpotifyID(ctx context.Context, id string) (*model.SongArtist, error) {
	artist, err := s.repo.queries(ctx).GetSongArtistBySpotifyID(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get artist by spotify id %s | %w", id, err)
	}

	return &model.SongArtist{
		ID:         int(artist.ID),
		Name:       artist.Name,
		SpotifyID:  artist.SpotifyID,
		Followers:  int(artist.Followers),
		Popularity: int(artist.Popularity),
	}, nil
}

func (s *Song) Create(ctx context.Context, song *model.Song) error {
	id, err := s.repo.queries(ctx).CreateSong(ctx, sqlc.CreateSongParams{
		Title:      song.Title,
		Album:      song.Album,
		SpotifyID:  song.SpotifyID,
		DurationMs: int32(song.DurationMS),
		LyricsType: pgtype.Text{String: song.LyricsType, Valid: song.LyricsType != ""},
		Lyrics:     pgtype.Text{String: song.Lyrics, Valid: song.Lyrics != ""},
	})
	if err != nil {
		return fmt.Errorf("create song %+v | %w", *song, err)
	}

	song.ID = int(id)

	return nil
}

func (s *Song) CreateHistory(ctx context.Context, id int) error {
	if _, err := s.repo.queries(ctx).CreateSongHistory(ctx, int32(id)); err != nil {
		return fmt.Errorf("create song history %d %w", id, err)
	}

	return nil
}

func (s *Song) CreateArtist(ctx context.Context, ) error {
	id, err := s.repo.queries(ctx).CreateSongArtistSong(ctx, sqlc.CreateSongArtistSongParams{
		ArtistID: artist.SpotifyID,
		SongID: ,
	})	
}
