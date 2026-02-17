package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/internal/server/service"
)

type Message struct {
	router fiber.Router

	message service.Message
	reply   service.Reply
}

func NewMessage(router fiber.Router, service service.Service) *Message {
	api := &Message{
		router:  router.Group("/messages"),
		message: *service.NewMessage(),
		reply:   *service.NewReply(),
	}

	api.createRoutes()

	return api
}

func (m *Message) createRoutes() {
	m.router.Get("/last", m.getLast)
	m.router.Post("/", m.create)
	m.router.Post("/:id/reply", m.replyMsg)
}

func (m *Message) getLast(c *fiber.Ctx) error {
	msg, err := m.message.GetLast(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(msg)
}

func (m *Message) create(c *fiber.Ctx) error {
	var message dto.MessageSave

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

	_, err = m.reply.Create(c.Context(), reply)
	if err != nil {
		return err
	}

	ref := c.Get("Referer")
	if ref == "" {
		ref = "/" // fallback
	}

	return c.Redirect(ref, fiber.StatusSeeOther)
}
