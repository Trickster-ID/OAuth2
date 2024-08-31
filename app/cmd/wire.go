//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"oauth2/app/routes"
)

func InitializeFiberServer() *fiber.App {
	wire.Build(
		routes.NewRouter,
	)
	return nil
}
