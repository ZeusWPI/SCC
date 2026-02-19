package ws

import (
	"context"
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/internal/server/service"
)

type Message struct {
	router  fiber.Router
	message service.Message
}

func NewMessage(router fiber.Router, service service.Service) *Message {
	ws := &Message{
		router:  router.Group("/messages"),
		message: *service.NewMessage(),
	}

	ws.createRoutes()

	return ws
}

func (m *Message) createRoutes() {
	m.router.Get("", func(c *fiber.Ctx) error {
		ip := c.IP()
		c.Locals("ip", ip)

		return websocket.New(m.create)(c)
	})
}

func (m *Message) create(c *websocket.Conn) {
	defer m.message.ListenerRemove(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ip, _ := c.Locals("ip").(string)

	for {
		_, body, err := c.ReadMessage()
		if err != nil {
			break
		}

		var frame dto.WSFrame
		if err := json.Unmarshal(body, &frame); err != nil {
			_ = c.WriteJSON(dto.WSFrame{
				Event: "error",
				Data:  map[string]any{"message": "invalid json"},
			})
			continue
		}

		// nolint:gocritic // Let me use switch grrrrr
		switch frame.Event {
		case "message":
			name, _ := frame.Data["username"].(string)
			message, _ := frame.Data["message"].(string)

			messageSave := dto.MessageSave{
				Name:    name,
				Message: message,
				IP:      ip,
			}

			if err := dto.Validate.Struct(messageSave); err != nil {
				_ = c.WriteJSON(dto.WSFrame{
					Event: "error",
					Data:  map[string]any{"message": err.Error()},
				})
			}

			created, err := m.message.Create(ctx, c, messageSave)
			if err != nil {
				_ = c.WriteJSON(dto.WSFrame{
					Event: "error",
					Data:  map[string]any{"message": "failed to save"},
				})
				continue
			}

			_ = c.WriteJSON(dto.WSFrame{
				Event: "ack",
				Data:  map[string]any{"id": created.ID},
			})
		}
	}
}
