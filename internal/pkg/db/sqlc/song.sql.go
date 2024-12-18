// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: song.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSong = `-- name: CreateSong :one

INSERT INTO song (title, album, spotify_id, duration_ms, lyrics_type, lyrics)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, title, spotify_id, duration_ms, album, lyrics_type, lyrics
`

type CreateSongParams struct {
	Title      string
	Album      string
	SpotifyID  string
	DurationMs int32
	LyricsType pgtype.Text
	Lyrics     pgtype.Text
}

// CRUD
func (q *Queries) CreateSong(ctx context.Context, arg CreateSongParams) (Song, error) {
	row := q.db.QueryRow(ctx, createSong,
		arg.Title,
		arg.Album,
		arg.SpotifyID,
		arg.DurationMs,
		arg.LyricsType,
		arg.Lyrics,
	)
	var i Song
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.SpotifyID,
		&i.DurationMs,
		&i.Album,
		&i.LyricsType,
		&i.Lyrics,
	)
	return i, err
}

const createSongArtist = `-- name: CreateSongArtist :one
INSERT INTO song_artist (name, spotify_id, followers, popularity)
VALUES ($1, $2, $3, $4)
RETURNING id, name, spotify_id, followers, popularity
`

type CreateSongArtistParams struct {
	Name       string
	SpotifyID  string
	Followers  int32
	Popularity int32
}

func (q *Queries) CreateSongArtist(ctx context.Context, arg CreateSongArtistParams) (SongArtist, error) {
	row := q.db.QueryRow(ctx, createSongArtist,
		arg.Name,
		arg.SpotifyID,
		arg.Followers,
		arg.Popularity,
	)
	var i SongArtist
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SpotifyID,
		&i.Followers,
		&i.Popularity,
	)
	return i, err
}

const createSongArtistGenre = `-- name: CreateSongArtistGenre :one
INSERT INTO song_artist_genre (artist_id, genre_id)
VALUES ($1, $2)
RETURNING id, artist_id, genre_id
`

type CreateSongArtistGenreParams struct {
	ArtistID int32
	GenreID  int32
}

func (q *Queries) CreateSongArtistGenre(ctx context.Context, arg CreateSongArtistGenreParams) (SongArtistGenre, error) {
	row := q.db.QueryRow(ctx, createSongArtistGenre, arg.ArtistID, arg.GenreID)
	var i SongArtistGenre
	err := row.Scan(&i.ID, &i.ArtistID, &i.GenreID)
	return i, err
}

const createSongArtistSong = `-- name: CreateSongArtistSong :one
INSERT INTO song_artist_song (artist_id, song_id)
VALUES ($1, $2)
RETURNING id, artist_id, song_id
`

type CreateSongArtistSongParams struct {
	ArtistID int32
	SongID   int32
}

func (q *Queries) CreateSongArtistSong(ctx context.Context, arg CreateSongArtistSongParams) (SongArtistSong, error) {
	row := q.db.QueryRow(ctx, createSongArtistSong, arg.ArtistID, arg.SongID)
	var i SongArtistSong
	err := row.Scan(&i.ID, &i.ArtistID, &i.SongID)
	return i, err
}

const createSongGenre = `-- name: CreateSongGenre :one
INSERT INTO song_genre (genre)
VALUES ($1)
RETURNING id, genre
`

func (q *Queries) CreateSongGenre(ctx context.Context, genre string) (SongGenre, error) {
	row := q.db.QueryRow(ctx, createSongGenre, genre)
	var i SongGenre
	err := row.Scan(&i.ID, &i.Genre)
	return i, err
}

const createSongHistory = `-- name: CreateSongHistory :one
INSERT INTO song_history (song_id)
VALUES ($1)
RETURNING id, song_id, created_at
`

func (q *Queries) CreateSongHistory(ctx context.Context, songID int32) (SongHistory, error) {
	row := q.db.QueryRow(ctx, createSongHistory, songID)
	var i SongHistory
	err := row.Scan(&i.ID, &i.SongID, &i.CreatedAt)
	return i, err
}

const getLastSongFull = `-- name: GetLastSongFull :many
SELECT s.id, s.title AS song_title, s.spotify_id, s.album, s.duration_ms, s.lyrics_type, s.lyrics, sh.created_at, a.id AS artist_id, a.name AS artist_name, a.spotify_id AS artist_spotify_id, a.followers AS artist_followers, a.popularity AS artist_popularity, g.id AS genre_id, g.genre AS genre, sh.created_at
FROM song_history sh
JOIN song s ON sh.song_id = s.id
LEFT JOIN song_artist_song sa ON s.id = sa.song_id
LEFT JOIN song_artist a ON sa.artist_id = a.id
LEFT JOIN song_artist_genre ag ON ag.artist_id = a.id
LEFT JOIN song_genre g ON ag.genre_id = g.id
WHERE sh.created_at = (SELECT MAX(created_at) FROM song_history)
ORDER BY a.name, g.genre
`

type GetLastSongFullRow struct {
	ID               int32
	SongTitle        string
	SpotifyID        string
	Album            string
	DurationMs       int32
	LyricsType       pgtype.Text
	Lyrics           pgtype.Text
	CreatedAt        pgtype.Timestamptz
	ArtistID         pgtype.Int4
	ArtistName       pgtype.Text
	ArtistSpotifyID  pgtype.Text
	ArtistFollowers  pgtype.Int4
	ArtistPopularity pgtype.Int4
	GenreID          pgtype.Int4
	Genre            pgtype.Text
	CreatedAt_2      pgtype.Timestamptz
}

func (q *Queries) GetLastSongFull(ctx context.Context) ([]GetLastSongFullRow, error) {
	rows, err := q.db.Query(ctx, getLastSongFull)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLastSongFullRow
	for rows.Next() {
		var i GetLastSongFullRow
		if err := rows.Scan(
			&i.ID,
			&i.SongTitle,
			&i.SpotifyID,
			&i.Album,
			&i.DurationMs,
			&i.LyricsType,
			&i.Lyrics,
			&i.CreatedAt,
			&i.ArtistID,
			&i.ArtistName,
			&i.ArtistSpotifyID,
			&i.ArtistFollowers,
			&i.ArtistPopularity,
			&i.GenreID,
			&i.Genre,
			&i.CreatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastSongHistory = `-- name: GetLastSongHistory :one
SELECT id, song_id, created_at
FROM song_history
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetLastSongHistory(ctx context.Context) (SongHistory, error) {
	row := q.db.QueryRow(ctx, getLastSongHistory)
	var i SongHistory
	err := row.Scan(&i.ID, &i.SongID, &i.CreatedAt)
	return i, err
}

const getSongArtistByName = `-- name: GetSongArtistByName :one
SELECT id, name, spotify_id, followers, popularity
FROM song_artist
WHERE name = $1
`

func (q *Queries) GetSongArtistByName(ctx context.Context, name string) (SongArtist, error) {
	row := q.db.QueryRow(ctx, getSongArtistByName, name)
	var i SongArtist
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SpotifyID,
		&i.Followers,
		&i.Popularity,
	)
	return i, err
}

const getSongArtistBySpotifyID = `-- name: GetSongArtistBySpotifyID :one
SELECT id, name, spotify_id, followers, popularity
FROM song_artist
WHERE spotify_id = $1
`

func (q *Queries) GetSongArtistBySpotifyID(ctx context.Context, spotifyID string) (SongArtist, error) {
	row := q.db.QueryRow(ctx, getSongArtistBySpotifyID, spotifyID)
	var i SongArtist
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SpotifyID,
		&i.Followers,
		&i.Popularity,
	)
	return i, err
}

const getSongBySpotifyID = `-- name: GetSongBySpotifyID :one

SELECT id, title, spotify_id, duration_ms, album, lyrics_type, lyrics
FROM song
WHERE spotify_id = $1
`

// Other
func (q *Queries) GetSongBySpotifyID(ctx context.Context, spotifyID string) (Song, error) {
	row := q.db.QueryRow(ctx, getSongBySpotifyID, spotifyID)
	var i Song
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.SpotifyID,
		&i.DurationMs,
		&i.Album,
		&i.LyricsType,
		&i.Lyrics,
	)
	return i, err
}

const getSongGenreByName = `-- name: GetSongGenreByName :one
SELECT id, genre
FROM song_genre
WHERE genre = $1
`

func (q *Queries) GetSongGenreByName(ctx context.Context, genre string) (SongGenre, error) {
	row := q.db.QueryRow(ctx, getSongGenreByName, genre)
	var i SongGenre
	err := row.Scan(&i.ID, &i.Genre)
	return i, err
}

const getSongHistory = `-- name: GetSongHistory :many
SELECT s.title, play_count, aggregated.created_at
FROM (
    SELECT sh.song_id, MAX(sh.created_at) AS created_at, COUNT(sh.song_id) AS play_count
    FROM song_history sh
    GROUP BY sh.song_id
) aggregated
JOIN song s ON aggregated.song_id = s.id
ORDER BY aggregated.created_at DESC
LIMIT 10
`

type GetSongHistoryRow struct {
	Title     string
	PlayCount int64
	CreatedAt interface{}
}

func (q *Queries) GetSongHistory(ctx context.Context) ([]GetSongHistoryRow, error) {
	rows, err := q.db.Query(ctx, getSongHistory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSongHistoryRow
	for rows.Next() {
		var i GetSongHistoryRow
		if err := rows.Scan(&i.Title, &i.PlayCount, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopArtists = `-- name: GetTopArtists :many
SELECT sa.id AS artist_id, sa.name AS artist_name, COUNT(sh.id) AS total_plays
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
GROUP BY sa.id, sa.name
ORDER BY total_plays DESC
LIMIT 10
`

type GetTopArtistsRow struct {
	ArtistID   int32
	ArtistName string
	TotalPlays int64
}

func (q *Queries) GetTopArtists(ctx context.Context) ([]GetTopArtistsRow, error) {
	rows, err := q.db.Query(ctx, getTopArtists)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopArtistsRow
	for rows.Next() {
		var i GetTopArtistsRow
		if err := rows.Scan(&i.ArtistID, &i.ArtistName, &i.TotalPlays); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopGenres = `-- name: GetTopGenres :many
SELECT g.genre AS genre_name, COUNT(sh.id) AS total_plays
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
JOIN song_artist_genre sag ON sa.id = sag.artist_id
JOIN song_genre g ON sag.genre_id = g.id
GROUP BY g.genre
ORDER BY total_plays DESC
LIMIT 10
`

type GetTopGenresRow struct {
	GenreName  string
	TotalPlays int64
}

func (q *Queries) GetTopGenres(ctx context.Context) ([]GetTopGenresRow, error) {
	rows, err := q.db.Query(ctx, getTopGenres)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopGenresRow
	for rows.Next() {
		var i GetTopGenresRow
		if err := rows.Scan(&i.GenreName, &i.TotalPlays); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopMonthlyArtists = `-- name: GetTopMonthlyArtists :many
SELECT sa.id AS artist_id, sa.name AS artist_name, COUNT(sh.id) AS total_plays
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY sa.id, sa.name
ORDER BY total_plays DESC
LIMIT 10
`

type GetTopMonthlyArtistsRow struct {
	ArtistID   int32
	ArtistName string
	TotalPlays int64
}

func (q *Queries) GetTopMonthlyArtists(ctx context.Context) ([]GetTopMonthlyArtistsRow, error) {
	rows, err := q.db.Query(ctx, getTopMonthlyArtists)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopMonthlyArtistsRow
	for rows.Next() {
		var i GetTopMonthlyArtistsRow
		if err := rows.Scan(&i.ArtistID, &i.ArtistName, &i.TotalPlays); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopMonthlyGenres = `-- name: GetTopMonthlyGenres :many
SELECT g.genre AS genre_name, COUNT(sh.id) AS total_plays
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
JOIN song_artist_genre sag ON sa.id = sag.artist_id
JOIN song_genre g ON sag.genre_id = g.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY g.genre
ORDER BY total_plays DESC
LIMIT 10
`

type GetTopMonthlyGenresRow struct {
	GenreName  string
	TotalPlays int64
}

func (q *Queries) GetTopMonthlyGenres(ctx context.Context) ([]GetTopMonthlyGenresRow, error) {
	rows, err := q.db.Query(ctx, getTopMonthlyGenres)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopMonthlyGenresRow
	for rows.Next() {
		var i GetTopMonthlyGenresRow
		if err := rows.Scan(&i.GenreName, &i.TotalPlays); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopMonthlySongs = `-- name: GetTopMonthlySongs :many
SELECT s.id AS song_id, s.title, COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY s.id, s.title
ORDER BY play_count DESC
LIMIT 10
`

type GetTopMonthlySongsRow struct {
	SongID    int32
	Title     string
	PlayCount int64
}

func (q *Queries) GetTopMonthlySongs(ctx context.Context) ([]GetTopMonthlySongsRow, error) {
	rows, err := q.db.Query(ctx, getTopMonthlySongs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopMonthlySongsRow
	for rows.Next() {
		var i GetTopMonthlySongsRow
		if err := rows.Scan(&i.SongID, &i.Title, &i.PlayCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopSongs = `-- name: GetTopSongs :many
SELECT s.id AS song_id, s.title, COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
GROUP BY s.id, s.title
ORDER BY play_count DESC
LIMIT 10
`

type GetTopSongsRow struct {
	SongID    int32
	Title     string
	PlayCount int64
}

func (q *Queries) GetTopSongs(ctx context.Context) ([]GetTopSongsRow, error) {
	rows, err := q.db.Query(ctx, getTopSongs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopSongsRow
	for rows.Next() {
		var i GetTopSongsRow
		if err := rows.Scan(&i.SongID, &i.Title, &i.PlayCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
