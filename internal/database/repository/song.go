package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/sqlc"
)

type Song struct {
	repo Repository
}

func (r *Repository) NewSong() *Song {
	return &Song{
		repo: *r,
	}
}

func (s *Song) GetLastPopulated(ctx context.Context) (*model.Song, error) {
	last, err := s.repo.queries(ctx).SongGetLastPopulated(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get last song populated %w", err)
	}

	song := model.SongModel(last[0].Song)
	song.PlayedAt = last[0].SongHistory.CreatedAt.Time

	for _, s := range last {
		song.Artists = append(song.Artists, *model.ArtistModel(s.SongArtist))
	}

	return song, nil
}

func (s *Song) GetLast50(ctx context.Context) ([]*model.Song, error) {
	lasts, err := s.repo.queries(ctx).SongGetLast50(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get last 50 songs %w", err)
	}

	songs := make([]*model.Song, 0, len(lasts))
	for _, last := range lasts {
		song := model.SongModel(last.Song)
		song.PlayCount = int(last.PlayCount)

		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Song) GetTopSongs(ctx context.Context) ([]*model.Song, error) {
	tops, err := s.repo.queries(ctx).SongGetTop50(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get top songs %w", err)
	}

	songs := make([]*model.Song, 0, len(tops))
	for _, top := range tops {
		song := model.SongModel(top.Song)
		song.PlayCount = int(top.PlayCount)

		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Song) GetTopArtists(ctx context.Context) ([]*model.Artist, error) {
	tops, err := s.repo.queries(ctx).SongArtistGetTop50(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get top artists %w", err)
	}

	artists := make([]*model.Artist, 0, len(tops))
	for _, top := range tops {
		artist := model.ArtistModel(top.SongArtist)
		artist.PlayCount = int(top.PlayCount)

		artists = append(artists, artist)
	}

	return artists, nil
}

func (s *Song) GetTopGenres(ctx context.Context) ([]*model.Genre, error) {
	tops, err := s.repo.queries(ctx).SongGenreGetTop50(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get top genres %w", err)
	}

	genres := make([]*model.Genre, 0, len(tops))
	for _, top := range tops {
		genre := model.GenreModel(top.SongGenre)
		genre.PlayCount = int(top.PlayCount)

		genres = append(genres, genre)
	}

	return genres, nil
}

func (s *Song) GetTopSongsMonthly(ctx context.Context) ([]*model.Song, error) {
	tops, err := s.repo.queries(ctx).SongGetTop50Monthly(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get top songs monthly %w", err)
	}

	songs := make([]*model.Song, 0, len(tops))
	for _, top := range tops {
		song := model.SongModel(top.Song)
		song.PlayCount = int(top.PlayCount)

		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Song) GetTopArtistsMonthly(ctx context.Context) ([]*model.Artist, error) {
	tops, err := s.repo.queries(ctx).SongArtistGetTop50Monthly(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get top artists montlhy %w", err)
	}

	artists := make([]*model.Artist, 0, len(tops))
	for _, top := range tops {
		artist := model.ArtistModel(top.SongArtist)
		artist.PlayCount = int(top.PlayCount)

		artists = append(artists, artist)
	}

	return artists, nil
}

func (s *Song) GetTopGenresMonthly(ctx context.Context) ([]*model.Genre, error) {
	tops, err := s.repo.queries(ctx).SongGenreGetTop50Monthly(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get top genres monthly %w", err)
	}

	genres := make([]*model.Genre, 0, len(tops))
	for _, top := range tops {
		genre := model.GenreModel(top.SongGenre)
		genre.PlayCount = int(top.PlayCount)

		genres = append(genres, genre)
	}

	return genres, nil
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
