// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: song.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createSong = `-- name: CreateSong :one

INSERT INTO song (title, album, spotify_id, duration_ms, lyrics_type, lyrics)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, title, spotify_id, duration_ms, album, lyrics_type, lyrics
`

type CreateSongParams struct {
	Title      string
	Album      string
	SpotifyID  string
	DurationMs int64
	LyricsType sql.NullString
	Lyrics     sql.NullString
}

// CRUD
func (q *Queries) CreateSong(ctx context.Context, arg CreateSongParams) (Song, error) {
	row := q.db.QueryRowContext(ctx, createSong,
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
VALUES (?, ?, ?, ?)
RETURNING id, name, spotify_id, followers, popularity
`

type CreateSongArtistParams struct {
	Name       string
	SpotifyID  string
	Followers  int64
	Popularity int64
}

func (q *Queries) CreateSongArtist(ctx context.Context, arg CreateSongArtistParams) (SongArtist, error) {
	row := q.db.QueryRowContext(ctx, createSongArtist,
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
VALUES (?, ?)
RETURNING id, artist_id, genre_id
`

type CreateSongArtistGenreParams struct {
	ArtistID int64
	GenreID  int64
}

func (q *Queries) CreateSongArtistGenre(ctx context.Context, arg CreateSongArtistGenreParams) (SongArtistGenre, error) {
	row := q.db.QueryRowContext(ctx, createSongArtistGenre, arg.ArtistID, arg.GenreID)
	var i SongArtistGenre
	err := row.Scan(&i.ID, &i.ArtistID, &i.GenreID)
	return i, err
}

const createSongArtistSong = `-- name: CreateSongArtistSong :one
INSERT INTO song_artist_song (artist_id, song_id)
VALUES (?, ?)
RETURNING id, artist_id, song_id
`

type CreateSongArtistSongParams struct {
	ArtistID int64
	SongID   int64
}

func (q *Queries) CreateSongArtistSong(ctx context.Context, arg CreateSongArtistSongParams) (SongArtistSong, error) {
	row := q.db.QueryRowContext(ctx, createSongArtistSong, arg.ArtistID, arg.SongID)
	var i SongArtistSong
	err := row.Scan(&i.ID, &i.ArtistID, &i.SongID)
	return i, err
}

const createSongGenre = `-- name: CreateSongGenre :one
INSERT INTO song_genre (genre)
VALUES (?)
RETURNING id, genre
`

func (q *Queries) CreateSongGenre(ctx context.Context, genre string) (SongGenre, error) {
	row := q.db.QueryRowContext(ctx, createSongGenre, genre)
	var i SongGenre
	err := row.Scan(&i.ID, &i.Genre)
	return i, err
}

const createSongHistory = `-- name: CreateSongHistory :one
INSERT INTO song_history (song_id)
VALUES (?)
RETURNING id, song_id, created_at
`

func (q *Queries) CreateSongHistory(ctx context.Context, songID int64) (SongHistory, error) {
	row := q.db.QueryRowContext(ctx, createSongHistory, songID)
	var i SongHistory
	err := row.Scan(&i.ID, &i.SongID, &i.CreatedAt)
	return i, err
}

const getLastSongFull = `-- name: GetLastSongFull :many
SELECT s.title AS song_title, s.spotify_id, s.album, s.duration_ms, s.lyrics_type, s.lyrics, a.name AS artist_name, g.genre AS genre
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
	SongTitle  string
	SpotifyID  string
	Album      string
	DurationMs int64
	LyricsType sql.NullString
	Lyrics     sql.NullString
	ArtistName sql.NullString
	Genre      sql.NullString
}

func (q *Queries) GetLastSongFull(ctx context.Context) ([]GetLastSongFullRow, error) {
	rows, err := q.db.QueryContext(ctx, getLastSongFull)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLastSongFullRow
	for rows.Next() {
		var i GetLastSongFullRow
		if err := rows.Scan(
			&i.SongTitle,
			&i.SpotifyID,
			&i.Album,
			&i.DurationMs,
			&i.LyricsType,
			&i.Lyrics,
			&i.ArtistName,
			&i.Genre,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
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
	row := q.db.QueryRowContext(ctx, getLastSongHistory)
	var i SongHistory
	err := row.Scan(&i.ID, &i.SongID, &i.CreatedAt)
	return i, err
}

const getSongArtistByName = `-- name: GetSongArtistByName :one
SELECT id, name, spotify_id, followers, popularity
FROM song_artist
WHERE name = ?
`

func (q *Queries) GetSongArtistByName(ctx context.Context, name string) (SongArtist, error) {
	row := q.db.QueryRowContext(ctx, getSongArtistByName, name)
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
WHERE spotify_id = ?
`

func (q *Queries) GetSongArtistBySpotifyID(ctx context.Context, spotifyID string) (SongArtist, error) {
	row := q.db.QueryRowContext(ctx, getSongArtistBySpotifyID, spotifyID)
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
WHERE spotify_id = ?
`

// Other
func (q *Queries) GetSongBySpotifyID(ctx context.Context, spotifyID string) (Song, error) {
	row := q.db.QueryRowContext(ctx, getSongBySpotifyID, spotifyID)
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
WHERE genre = ?
`

func (q *Queries) GetSongGenreByName(ctx context.Context, genre string) (SongGenre, error) {
	row := q.db.QueryRowContext(ctx, getSongGenreByName, genre)
	var i SongGenre
	err := row.Scan(&i.ID, &i.Genre)
	return i, err
}
