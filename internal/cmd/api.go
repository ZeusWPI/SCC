package cmd

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zeusWPI/scc/internal/api"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/song"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// API starts the API server
func API(db *db.DB, song *song.Song) {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024,
	})
	app.Use(
		fiberzap.New(fiberzap.Config{
			Logger: zap.L(),
		}),
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
		}),
	)

	apiGroup := app.Group("/api")
	api.New(apiGroup, db, song)

	host := config.GetDefaultString("server.host", "localhost")
	port := config.GetDefaultInt("server.port", 3000)

	zap.S().Fatal("API: Fatal server error", app.Listen(fmt.Sprintf("%s:%d", host, port)))
}
