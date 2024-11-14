// Package api provides all the API endpoints
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/api/message"
	"github.com/zeusWPI/scc/internal/pkg/db"
)

// New creates a new API instance
func New(router fiber.Router, db *db.DB) {
	message.New(router, db)
}
