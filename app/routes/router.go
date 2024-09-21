package routes

import (
	"github.com/gofiber/fiber/v3"
	"oauth2/app/configuration"
	"oauth2/app/controllers"
	"oauth2/app/global/helper"
	"oauth2/app/middleware"
)

func NewRouter(authController controllers.IAuthController) *fiber.App {
	f := fiber.New()
	configuration.FiberInitLogger(f)

	f.Get("/ping", func(c fiber.Ctx) error { return helper.Response(c, nil) })

	authGroup := f.Group("/auth")
	{
		authGroup.Post("/login", authController.Login, middleware.BasicAuthMiddleware())
		authGroup.Get("/verify-token", func(c fiber.Ctx) error { return helper.Response(c, nil) }, middleware.TokenMiddleware())
		authGroup.Post("/refresh-token", authController.RefreshToken, middleware.BasicAuthMiddleware())
	}
	return f
}
