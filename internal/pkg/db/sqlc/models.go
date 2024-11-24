// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"time"
)

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

type Spotify struct {
	ID         int64
	Title      string
	Artists    string
	SpotifyID  string
	DurationMs int64
	CreatedAt  time.Time
}

type Tap struct {
	ID             int64
	OrderID        int64
	OrderCreatedAt time.Time
	Name           string
	Category       string
	CreatedAt      time.Time
}
