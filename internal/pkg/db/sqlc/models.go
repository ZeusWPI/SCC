// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql"
	"time"
)

type Event struct {
	ID           int64
	Name         string
	Date         time.Time
	AcademicYear string
	Location     string
}

type Gamification struct {
	ID     int64
	Name   string
	Score  int64
	Avatar string
}

type Message struct {
	ID        int64
	Name      string
	Ip        string
	Message   string
	CreatedAt time.Time
}

type Scan struct {
	ID       int64
	ScanTime time.Time
}

type Season struct {
	ID      int64
	Name    string
	Start   time.Time
	End     time.Time
	Current bool
}

type Song struct {
	ID         int64
	Title      string
	SpotifyID  string
	DurationMs int64
	Album      string
	LyricsType sql.NullString
	Lyrics     sql.NullString
}

type SongArtist struct {
	ID         int64
	Name       string
	SpotifyID  string
	Followers  int64
	Popularity int64
}

type SongArtistGenre struct {
	ID       int64
	ArtistID int64
	GenreID  int64
}

type SongArtistSong struct {
	ID       int64
	ArtistID int64
	SongID   int64
}

type SongGenre struct {
	ID    int64
	Genre string
}

type SongHistory struct {
	ID        int64
	SongID    int64
	CreatedAt time.Time
}

type Tap struct {
	ID             int64
	OrderID        int64
	OrderCreatedAt time.Time
	Name           string
	Category       string
	CreatedAt      time.Time
}
