package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/internal/server/service"
)

type Song struct {
	router fiber.Router
	song   service.Song
}

func NewSong(router fiber.Router, service service.Service) *Song {
	api := &Song{
		router: router.Group("/song"),
		song:   *service.NewSong(),
	}

	api.createRoutes()

	return api
}

func (s *Song) createRoutes() {
	s.router.Post("/", s.new)
}

func (s *Song) new(c *fiber.Ctx) error {
	var song dto.Song

	if err := c.BodyParser(&song); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(song); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.song.New(c.Context(), song); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
