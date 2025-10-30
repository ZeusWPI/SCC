package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/sqlc"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Song struct {
	repo Repository
}

func (r *Repository) NewSong() *Song {
	return &Song{
		repo: *r,
	}
}

func (s *Song) GetBySpotify(ctx context.Context, spotifyID string) (*model.Song, error) {
	song, err := s.repo.queries(ctx).SongGetBySpotify(ctx, spotifyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get song by spotify id %s | %w", spotifyID, err)
	}

	return model.SongModel(song), nil
}

func (s *Song) GetArtistBySpotify(ctx context.Context, spotifyID string) (*model.Artist, error) {
	artist, err := s.repo.queries(ctx).SongArtistGetBySpotify(ctx, spotifyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get song artist by spotify id %s | %w", spotifyID, err)
	}

	return model.ArtistModel(artist), nil
}

func (s *Song) GetGenreByGenre(ctx context.Context, genre string) (*model.Genre, error) {
	genreDB, err := s.repo.queries(ctx).SongGenreGetByGenre(ctx, genre)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get song genre by genre %s | %w", genre, err)
	}

	return model.GenreModel(genreDB), nil
}

func (s *Song) GetArtistsBySpotify(ctx context.Context, artists []model.Artist) ([]*model.Artist, error) {
	artistsDB, err := s.repo.queries(ctx).SongArtistGetBySpotifyIds(ctx, utils.SliceMap(artists, func(a model.Artist) string { return a.SpotifyID }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get artists by ids %+v | %w", artists, err)
	}

	artistMap := make(map[int]*model.Artist)
	for _, artist := range artistsDB {
		a, ok := artistMap[int(artist.SongArtist.ID)]
		if !ok {
			a = model.ArtistModel(artist.SongArtist)
			artistMap[a.ID] = a
		}

		a.Genres = append(a.Genres, *model.GenreModel(artist.SongGenre))
	}

	return utils.MapValues(artistMap), nil
}

func (s *Song) Create(ctx context.Context, song *model.Song) error {
	return s.repo.WithRollback(ctx, func(ctx context.Context) error {
		id, err := s.repo.queries(ctx).SongCreate(ctx, sqlc.SongCreateParams{
			Title:      song.Title,
			Album:      song.Album,
			SpotifyID:  song.SpotifyID,
			DurationMs: int32(song.DurationMS),
			LyricsType: sqlc.LyricsTypeEnum(song.LyricsType),
			Lyrics:     pgtype.Text{String: song.Lyrics, Valid: song.Lyrics != ""},
		})
		if err != nil {
			return fmt.Errorf("create song %+v | %w", *song, err)
		}

		song.ID = int(id)

		for i, artist := range song.Artists {
			artistDB, err := s.GetArtistBySpotify(ctx, artist.SpotifyID)
			if err != nil {
				return err
			}

			if artistDB != nil {
				song.Artists[i] = *artistDB
			} else {
				id, err := s.repo.queries(ctx).SongArtistCreate(ctx, sqlc.SongArtistCreateParams{
					Name:      artist.Name,
					SpotifyID: artist.SpotifyID,
				})
				if err != nil {
					return fmt.Errorf("create song artist %+v | %+v | %w", artist, *song, err)
				}

				song.Artists[i].ID = int(id)
			}

			if _, err := s.repo.queries(ctx).SongArtistSongCreate(ctx, sqlc.SongArtistSongCreateParams{
				ArtistID: int32(song.Artists[i].ID),
				SongID:   int32(song.ID),
			}); err != nil {
				return fmt.Errorf("create song artist song %+v | %+v | %w", artist, *song, err)
			}

			for j, genre := range artist.Genres {
				genreDB, err := s.GetGenreByGenre(ctx, genre.Genre)
				if err != nil {
					return err
				}

				if genreDB != nil {
					song.Artists[i].Genres[j] = *genreDB
				} else {
					id, err := s.repo.queries(ctx).SongGenreCreate(ctx, genre.Genre)
					if err != nil {
						return fmt.Errorf("create song genre %+v | %+v | %+v | %w", genre, artist, *song, err)
					}

					song.Artists[i].Genres[j].ID = int(id)
				}

				if _, err := s.repo.queries(ctx).SongArtistGenreCreate(ctx, sqlc.SongArtistGenreCreateParams{
					ArtistID: int32(song.Artists[i].ID),
					GenreID:  int32(song.Artists[i].Genres[j].ID),
				}); err != nil {
					return fmt.Errorf("create song artist genre %+v | %+v | %+v | %w", genre, artist, *song, err)
				}
			}
		}

		return nil
	})
}

func (s *Song) CreateHistory(ctx context.Context, song model.Song) error {
	if _, err := s.repo.queries(ctx).SongHistoryCreate(ctx, int32(song.ID)); err != nil {
		return fmt.Errorf("create song history %+v | %w", song, err)
	}

	return nil
}
