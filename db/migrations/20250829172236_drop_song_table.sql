-- +goose Up
-- +goose StatementBegin
DROP TABLE song_history;
DROP TABLE song_artist;
DROP TABLE song_genre;
DROP TABLE song_artist_song;
DROP TABLE song_artist_genre;
DROP TABLE song;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE song (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    duration_ms INTEGER NOT NULL,
    album TEXT NOT NULL
    lyrics_type TEXT,
    lyrics TEXT
);

CREATE TABLE song_history (
    id SERIAL PRIMARY KEY,
    song_id INTEGER NOT NULL REFERENCES song (id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() - INTERVAL '1 second')
);

CREATE TABLE song_genre (
    id SERIAL PRIMARY KEY,
    genre TEXT NOT NULL
);

CREATE TABLE song_artist (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    followers INTEGER NOT NULL,
    popularity INTEGER NOT NULL
);

CREATE TABLE song_artist_song (
    id SERIAL PRIMARY KEY,
    artist_id INTEGER NOT NULL,
    song_id INTEGER NOT NULL,
    FOREIGN KEY(artist_id) REFERENCES song_artist(id)  ON DELETE CASCADE,
    FOREIGN KEY(song_id) REFERENCES song(id)  ON DELETE CASCADE
);

CREATE TABLE song_artist_genre (
    id SERIAL PRIMARY KEY,
    artist_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    FOREIGN KEY(artist_id) REFERENCES song_artist(id)  ON DELETE CASCADE,
    FOREIGN KEY(genre_id) REFERENCES song_genre(id)  ON DELETE CASCADE
);
-- +goose StatementEnd
