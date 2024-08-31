package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"net/http"
	"os"
)

type Router struct {
}

func NewRouter() *fiber.App {
	// Initialize a new Fiber app
	f := fiber.New()
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "info" || logLevel == "debug" || logLevel == "trace" {
		f.Use(logger.New())
	}
	// Define a route for the GET method on the root path '/'
	f.Get("/ping", func(c fiber.Ctx) error {
		//time.Sleep(time.Second * 10)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	return f
}
