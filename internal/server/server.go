// Package server starts the API server
package server

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	routers "github.com/zeusWPI/scc/internal/server/api"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/internal/server/service"
	"github.com/zeusWPI/scc/internal/server/ws"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	*fiber.App
	Addr string
}

func New(service service.Service) *Server {
	env := config.GetDefaultString("app.env", "development")
	port := config.GetDefaultInt("server.port", 4000)
	host := config.GetDefaultString("server.host", "0.0.0.0")

	// Web template engine
	engine := html.New("./ui", ".html")
	engine.AddFunc("lastMessage", func(msgs []dto.Message) dto.Message {
		return msgs[len(msgs)-1]
	})
	engine.AddFunc("maxMessageID", func(days []dto.MessageDayGroup) int {
		maxID := 0
		for _, d := range days {
			for _, c := range d.Clusters {
				for _, m := range c.Messages {
					if m.ID > maxID {
						maxID = m.ID
					}
				}
			}
		}
		return maxID
	})
	engine.Reload(env != "production")

	// Construct app
	app := fiber.New(fiber.Config{
		BodyLimit:      16 * 1024 * 1024,
		ReadBufferSize: 8096,
		Views:          engine,
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))
	if env != "production" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     fmt.Sprintf("http://localhost:%d", port),
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
		}))
	}

	// Register web socket routes
	socket := app.Group("/ws")
	socket.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ws.NewMessage(socket, service)

	// Register api routes
	api := app.Group("/")

	routers.NewMessage(api, service)
	routers.NewSong(api, service)

	// Fallback
	app.All("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	srv := &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}

	return srv
}
