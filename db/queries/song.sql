-- CRUD

-- name: CreateSong :one
INSERT INTO song (title, album, spotify_id, duration_ms, lyrics_type, lyrics)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: CreateSongHistory :one
INSERT INTO song_history (song_id)
VALUES (?)
RETURNING *;

-- name: CreateSongGenre :one
INSERT INTO song_genre (genre)
VALUES (?)
RETURNING *;

-- name: CreateSongArtist :one
INSERT INTO song_artist (name, spotify_id, followers, popularity)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: CreateSongArtistSong :one
INSERT INTO song_artist_song (artist_id, song_id)
VALUES (?, ?)
RETURNING *;

-- name: CreateSongArtistGenre :one
INSERT INTO song_artist_genre (artist_id, genre_id)
VALUES (?, ?)
RETURNING *;


-- Other

-- name: GetSongBySpotifyID :one
SELECT *
FROM song
WHERE spotify_id = ?;

-- name: GetSongArtistBySpotifyID :one
SELECT *
FROM song_artist
WHERE spotify_id = ?;

-- name: GetLastSongHistory :one
SELECT *
FROM song_history
ORDER BY created_at DESC
LIMIT 1;

-- name: GetSongGenreByName :one
SELECT *
FROM song_genre
WHERE genre = ?;

-- name: GetSongArtistByName :one
SELECT *
FROM song_artist
WHERE name = ?;

-- name: GetLastSongFull :many
SELECT s.title AS song_title, s.spotify_id, s.album, s.duration_ms, s.lyrics_type, s.lyrics, a.name AS artist_name, g.genre AS genre
FROM song_history sh
JOIN song s ON sh.song_id = s.id
LEFT JOIN song_artist_song sa ON s.id = sa.song_id
LEFT JOIN song_artist a ON sa.artist_id = a.id
LEFT JOIN song_artist_genre ag ON ag.artist_id = a.id
LEFT JOIN song_genre g ON ag.genre_id = g.id
WHERE sh.created_at = (SELECT MAX(created_at) FROM song_history)
ORDER BY a.name, g.genre;
