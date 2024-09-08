package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"oauth2/app/controllers"
	"os"
)

type Router struct {
}

func NewRouter(authController controllers.IAuthController) *fiber.App {
	// Initialize a new Fiber app
	f := fiber.New()
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "info" || logLevel == "debug" || logLevel == "trace" {
		f.Use(logger.New())
	}
	// Define a route for the GET method on the root path '/'
	var hash []byte
	f.Get("/ping", func(c fiber.Ctx) error {
		//time.Sleep(time.Second * 10)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "success"})
	})
	f.Get("/encrypt", func(c fiber.Ctx) error {
		res, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
		if err != nil {
			logrus.Errorln(err.Error())
			return err
		}
		hash = res
		fmt.Println(res)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "success", "password": res})
	})
	f.Get("/compare", func(c fiber.Ctx) error {
		err := bcrypt.CompareHashAndPassword(hash, []byte("admin"))
		if err != nil {
			logrus.Errorln(err.Error())
			return err
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	f.Get("/login", authController.Login)
	//v1WithMiddleware := f.Group("v1", middleware.Authorization())
	return f
}
