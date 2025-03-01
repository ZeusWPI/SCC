-- CRUD

-- name: CreateSong :one
INSERT INTO song (title, album, spotify_id, duration_ms, lyrics_type, lyrics)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateSongHistory :one
INSERT INTO song_history (song_id)
VALUES ($1)
RETURNING *;

-- name: CreateSongGenre :one
INSERT INTO song_genre (genre)
VALUES ($1)
RETURNING *;

-- name: CreateSongArtist :one
INSERT INTO song_artist (name, spotify_id, followers, popularity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateSongArtistSong :one
INSERT INTO song_artist_song (artist_id, song_id)
VALUES ($1, $2)
RETURNING *;

-- name: CreateSongArtistGenre :one
INSERT INTO song_artist_genre (artist_id, genre_id)
VALUES ($1, $2)
RETURNING *;


-- Other

-- name: GetSongBySpotifyID :one
SELECT *
FROM song
WHERE spotify_id = $1;

-- name: GetSongArtistBySpotifyID :one
SELECT *
FROM song_artist
WHERE spotify_id = $1;

-- name: GetLastSongHistory :one
SELECT *
FROM song_history
ORDER BY created_at DESC
LIMIT 1;

-- name: GetSongGenreByName :one
SELECT *
FROM song_genre
WHERE genre = $1;

-- name: GetSongArtistByName :one
SELECT *
FROM song_artist
WHERE name = $1;

-- name: GetLastSongFull :many
SELECT s.id, s.title AS song_title, s.spotify_id, s.album, s.duration_ms, s.lyrics_type, s.lyrics, sh.created_at, a.id AS artist_id, a.name AS artist_name, a.spotify_id AS artist_spotify_id, a.followers AS artist_followers, a.popularity AS artist_popularity, g.id AS genre_id, g.genre AS genre, sh.created_at
FROM song_history sh
JOIN song s ON sh.song_id = s.id
LEFT JOIN song_artist_song sa ON s.id = sa.song_id
LEFT JOIN song_artist a ON sa.artist_id = a.id
LEFT JOIN song_artist_genre ag ON ag.artist_id = a.id
LEFT JOIN song_genre g ON ag.genre_id = g.id
WHERE sh.created_at = (SELECT MAX(created_at) FROM song_history)
ORDER BY a.name, g.genre;

-- name: GetSongHistory :many
SELECT s.title, play_count, aggregated.created_at
FROM (
    SELECT sh.song_id, MAX(sh.created_at) AS created_at, COUNT(sh.song_id) AS play_count
    FROM song_history sh
    GROUP BY sh.song_id
) aggregated
JOIN song s ON aggregated.song_id = s.id
ORDER BY aggregated.created_at DESC
LIMIT 50;

-- name: GetTopSongs :many
SELECT s.id AS song_id, s.title, COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
GROUP BY s.id, s.title
ORDER BY play_count DESC
LIMIT 10;

-- name: GetTopArtists :many
SELECT sa.id AS artist_id, sa.name AS artist_name, COUNT(sh.id) AS total_plays
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
GROUP BY sa.id, sa.name
ORDER BY total_plays DESC
LIMIT 10;

-- name: GetTopGenres :many
SELECT g.genre AS genre_name, COUNT(sh.id) AS total_plays
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
JOIN song_artist_genre sag ON sa.id = sag.artist_id
JOIN song_genre g ON sag.genre_id = g.id
GROUP BY g.genre
ORDER BY total_plays DESC
LIMIT 10;

-- name: GetTopMonthlySongs :many
SELECT s.id AS song_id, s.title, COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY s.id, s.title
ORDER BY play_count DESC
LIMIT 10;

-- name: GetTopMonthlyArtists :many
SELECT sa.id AS artist_id, sa.name AS artist_name, COUNT(sh.id) AS total_plays
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY sa.id, sa.name
ORDER BY total_plays DESC
LIMIT 10;

-- name: GetTopMonthlyGenres :many
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
LIMIT 10;
