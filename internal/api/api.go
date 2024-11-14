// Package api provides all the API endpoints
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/api/message"
	apiSpotify "github.com/zeusWPI/scc/internal/api/spotify"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/spotify"
)

// New creates a new API instance
func New(router fiber.Router, db *db.DB, spotify *spotify.Spotify) {
	message.New(router, db)
	apiSpotify.New(router, db, spotify)
}
