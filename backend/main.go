package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	pool := mustDB()
	defer pool.Close()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	registerAuthRoutes(app, pool)

	registerUX(app, pool)
	registerQA(app, pool)
	registerBA(app, pool)

	registerScoreRoutes(app, pool)
	registerManagerRoutes(app, pool)
	registerProjectRoutes(app, pool)

	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
