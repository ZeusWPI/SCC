// Package server starts the API server
package server

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	routers "github.com/zeusWPI/scc/internal/server/api"
	"github.com/zeusWPI/scc/internal/server/service"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	*fiber.App
	Addr string
}

func New(service service.Service, pool *pgxpool.Pool) *Server {
	env := config.GetDefaultString("app.env", "development")

	// Construct app
	app := fiber.New(fiber.Config{
		BodyLimit:      16 * 1024 * 1024,
		ReadBufferSize: 8096,
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))
	if env != "production" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:3000",
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
		}))
	}

	// Register routes
	api := app.Group("/api")

	routers.NewMessage(api, service)
	routers.NewSong(api, service)

	// Fallback
	app.All("/api*", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	port := config.GetDefaultInt("server.port", 4000)
	host := config.GetDefaultString("server.host", "0.0.0.0")

	srv := &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}

	return srv
}
