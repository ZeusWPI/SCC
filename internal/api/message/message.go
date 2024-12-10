// Package message provides the API regarding the cammie chat messages
package message

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/buzzer"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/util"
	"go.uber.org/zap"
)

// Router is the message API router
type Router struct {
	router fiber.Router
	db     *db.DB
	buzz   *buzzer.Buzzer
}

// New creates a new message API instance
func New(router fiber.Router, db *db.DB, buzz *buzzer.Buzzer) *Router {
	api := &Router{
		router: router.Group("/messages"),
		db:     db,
		buzz:   buzz,
	}
	api.createRoutes()

	return api
}

func (r *Router) createRoutes() {
	r.router.Get("/", r.getAll)
	r.router.Post("/", r.create)
}

func (r *Router) getAll(c *fiber.Ctx) error {
	messages, err := r.db.Queries.GetAllMessages(c.Context())
	if err != nil {
		zap.S().Error("DB: Get all messages\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(util.SliceMap(messages, dto.MessageDTO))
}

func (r *Router) create(c *fiber.Ctx) error {
	message := new(dto.Message)

	if err := c.BodyParser(message); err != nil {
		zap.S().Error("API: Message body parser\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := dto.Validate.Struct(message); err != nil {
		zap.S().Error("API: Message  validation\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	messageDB, err := r.db.Queries.CreateMessage(c.Context(), message.CreateParams())
	if err != nil {
		zap.S().Error("DB: Create message\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	r.buzz.Play()

	return c.Status(fiber.StatusCreated).JSON(dto.MessageDTO(messageDB))
}
