package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/server/service"
)

type Message struct {
	router  fiber.Router
	message service.Message
}

func NewMessage(router fiber.Router, service service.Service) *Message {
	web := &Message{
		router:  router.Group("/messages"),
		message: *service.NewMessage(),
	}

	web.createRoutes()

	return web
}

func (m *Message) createRoutes() {
	m.router.Get("/", m.index)
}

func (m *Message) index(c *fiber.Ctx) error {
	groups, err := m.message.Get(c.Context(), -1, 100)
	if err != nil {
		return err
	}

	lastID := 0
	for _, g := range groups {
		for _, cl := range g.Clusters {
			for _, msg := range cl.Messages {
				if msg.ID > lastID {
					lastID = msg.ID
				}
			}
		}
	}

	return c.Render("pages/index", fiber.Map{
		"Days":   groups,
		"LastID": lastID,
	}, "layout/main")
}
