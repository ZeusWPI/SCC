-- +goose Up
-- +goose StatementBegin
ALTER TABLE song
DROP COLUMN artists;

CREATE TABLE IF NOT EXISTS song_genre (
    id SERIAL PRIMARY KEY,
    genre TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS song_artist (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    followers INTEGER NOT NULL,
    popularity INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS song_artist_song (
    id SERIAL PRIMARY KEY,
    artist_id INTEGER NOT NULL,
    song_id INTEGER NOT NULL,
    FOREIGN KEY(artist_id) REFERENCES song_artist(id)  ON DELETE CASCADE,
    FOREIGN KEY(song_id) REFERENCES song(id)  ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS song_artist_genre (
    id SERIAL PRIMARY KEY,
    artist_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    FOREIGN KEY(artist_id) REFERENCES song_artist(id)  ON DELETE CASCADE,
    FOREIGN KEY(genre_id) REFERENCES song_genre(id)  ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS song_artist_genre;
DROP TABLE IF EXISTS song_artist_song;
DROP TABLE IF EXISTS song_artist;
DROP TABLE IF EXISTS song_genre;

ALTER TABLE song
ADD COLUMN artists TEXT NOT NULL DEFAULT 'Unknown';
-- +goose StatementEnd
