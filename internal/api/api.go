// Package api provides all the API endpoints
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/api/message"
	apiSong "github.com/zeusWPI/scc/internal/api/song"
	"github.com/zeusWPI/scc/internal/pkg/buzzer"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/song"
)

// New creates a new API instance
func New(router fiber.Router, db *db.DB, buzz *buzzer.Buzzer, song *song.Song) {
	message.New(router, db, buzz)
	apiSong.New(router, db, song)
}
