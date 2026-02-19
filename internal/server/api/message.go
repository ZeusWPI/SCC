package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/internal/server/service"
)

type Message struct {
	router fiber.Router

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
	m.router.Get("/last", m.getLast)
	m.router.Post("/:id/reply", m.replyMsg)
}

func (m *Message) getLast(c *fiber.Ctx) error {
	msg, err := m.message.GetLast(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(msg)
}

func (m *Message) replyMsg(c *fiber.Ctx) error {
	var reply dto.ReplySave

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	reply.MessageID = id

	if err := c.BodyParser(&reply); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(reply); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = m.message.Reply(c.Context(), reply)
	wsGone := false
	if err != nil {
		if !errors.Is(err, fiber.ErrGone) {
			return err
		}
		wsGone = true
	}

	ref := c.Get("Referer")
	if ref == "" {
		ref = "/"
	}

	if wsGone {
		c.Cookie(&fiber.Cookie{
			Name:     "flash_ws_gone",
			Value:    "1",
			Path:     "/",
			HTTPOnly: true,
			SameSite: "Lax",
			MaxAge:   60,
		})
	}

	return c.Redirect(ref, fiber.StatusSeeOther)
}
