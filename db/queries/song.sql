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
