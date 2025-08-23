package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/internal/server/service"
)

type Message struct {
	router  fiber.Router
	message service.Message
}

func NewMessage(router fiber.Router, service service.Service) *Message {
	api := &Message{
		router:  router.Group("/messages"),
		message: *service.NewMessage(),
	}

	api.createRoutes()

	return api
}

func (m *Message) createRoutes() {
	m.router.Post("/", m.create)
}

func (m *Message) create(c *fiber.Ctx) error {
	var message dto.Message

	if err := c.BodyParser(&message); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(message); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newMessage, err := m.message.Create(c.Context(), message)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(newMessage)
}
