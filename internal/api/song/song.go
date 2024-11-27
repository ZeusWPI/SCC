// Package song provides the API regarding songs integration
package song

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/internal/pkg/song"
	"go.uber.org/zap"
)

// Router is the song API router
type Router struct {
	router fiber.Router
	db     *db.DB
	song   *song.Song
}

// New creates a new song API instance
func New(router fiber.Router, db *db.DB, song *song.Song) *Router {
	api := &Router{
		router: router.Group("/song"),
		db:     db,
		song:   song,
	}
	api.createRoutes()

	return api
}

func (r *Router) createRoutes() {
	r.router.Post("/", r.new)
}

func (r *Router) new(c *fiber.Ctx) error {
	song := new(dto.Song)

	if err := c.BodyParser(song); err != nil {
		zap.S().Error("API: Song body parser\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := dto.Validate.Struct(song); err != nil {
		zap.S().Error("API: Song validation\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	go func() {
		err := r.song.Track(song)
		if err != nil {
			zap.S().Error("Song: Get Track\n", err)
		}
	}()

	return c.SendStatus(fiber.StatusOK)
}
