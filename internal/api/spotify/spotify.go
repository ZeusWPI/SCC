// Package spotify provides the API regarding spotify integration
package spotify

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/internal/pkg/spotify"
	"go.uber.org/zap"
)

// Router is the spotify API router
type Router struct {
	router  fiber.Router
	db      *db.DB
	spotify *spotify.Spotify
}

// New creates a new spotify API instance
func New(router fiber.Router, db *db.DB, spotify *spotify.Spotify) *Router {
	api := &Router{
		router:  router.Group("/spotify"),
		db:      db,
		spotify: spotify,
	}
	api.createRoutes()

	return api
}

func (r *Router) createRoutes() {
	r.router.Post("/", r.new)
}

func (r *Router) new(c *fiber.Ctx) error {
	spotify := new(dto.Spotify)

	if err := c.BodyParser(spotify); err != nil {
		zap.S().Error("API: Spotify body parser\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := dto.Validate.Struct(spotify); err != nil {
		zap.S().Error("API: Spotify validation\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	go func() {
		err := r.spotify.Track(spotify)
		if err != nil {
			zap.S().Error("Spotify: Get Track\n", err)
		}
	}()

	return c.SendStatus(fiber.StatusOK)
}
