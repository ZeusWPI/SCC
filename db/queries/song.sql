-- name: SongGetLastPopulated :many
SELECT sqlc.embed(h), sqlc.embed(s), sqlc.embed(a)
FROM song_history h
JOIN song s ON s.id = h.song_id
LEFT JOIN song_artist_song sa ON sa.song_id = s.id
LEFT JOIN song_artist a ON a.id = sa.artist_id
WHERE h.created_at = (SELECT MAX(created_at) FROM song_history)
ORDER BY a.name;

-- name: SongGetLast50 :many
SELECT sqlc.embed(s), play_count
FROM (
    SELECT sh.song_id, MAX(sh.created_at) AS created_at, COUNT(sh.song_id) AS play_count
    FROM song_history sh
    GROUP BY sh.song_id
) aggregated
JOIN song s ON s.id = aggregated.song_id
ORDER BY aggregated.created_at DESC
LIMIT 50;

-- name: SongGetTop50 :many
SELECT sqlc.embed(s), COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
GROUP BY s.id, s.title
ORDER BY play_count DESC
LIMIT 50;

-- name: SongGetTop50Monthly :many
SELECT sqlc.embed(s), COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY s.id, s.title
ORDER BY play_count DESC
LIMIT 50;

-- name: SongArtistGetTop50 :many
SELECT sqlc.embed(a), COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist a ON sas.artist_id = a.id
GROUP BY a.id, a.name
ORDER BY play_count DESC
LIMIT 50;

-- name: SongArtistGetTop50Monthly :many
SELECT sqlc.embed(a), COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist a ON sas.artist_id = a.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY a.id, a.name
ORDER BY play_count DESC
LIMIT 50;

-- name: SongGenreGetTop50 :many
SELECT sqlc.embed(g), COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
JOIN song_artist_genre sag ON sa.id = sag.artist_id
JOIN song_genre g ON sag.genre_id = g.id
GROUP BY g.genre, g.id
ORDER BY play_count DESC
LIMIT 50;

-- name: SongGenreGetTop50Monthly :many
SELECT sqlc.embed(g), COUNT(sh.id) AS play_count
FROM song_history sh
JOIN song s ON sh.song_id = s.id
JOIN song_artist_song sas ON s.id = sas.song_id
JOIN song_artist sa ON sas.artist_id = sa.id
JOIN song_artist_genre sag ON sa.id = sag.artist_id
JOIN song_genre g ON sag.genre_id = g.id
WHERE sh.created_at > CURRENT_TIMESTAMP - INTERVAL '1 month'
GROUP BY g.genre, g.id
ORDER BY play_count DESC
LIMIT 50;

-- name: SongGetBySpotify :one
SELECT *
FROM song
WHERE spotify_id = $1;

-- name: SongArtistGetBySpotify :one
SELECT *
FROM song_artist
WHERE spotify_id = $1;

-- name: SongGenreGetByGenre :one
SELECT *
FROM song_genre
WHERE genre = $1;

-- name: SongCreate :one
INSERT INTO song (title, album, spotify_id, duration_ms, lyrics_type, lyrics)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: SongArtistCreate :one
INSERT INTO song_artist (name, spotify_id)
VALUES ($1, $2)
RETURNING id;

-- name: SongArtistSongCreate :one
INSERT INTO song_artist_song (artist_id, song_id)
VALUES ($1, $2)
RETURNING id;

-- name: SongGenreCreate :one
INSERT INTO song_genre (genre)
VALUES ($1)
RETURNING id;

-- name: SongArtistGenreCreate :one
INSERT INTO song_artist_genre (artist_id, genre_id)
VALUES ($1, $2)
RETURNING id;

-- name: SongHistoryCreate :one
INSERT INTO song_history (song_id)
VALUES ($1)
RETURNING id;
