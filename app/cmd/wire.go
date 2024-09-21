//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"oauth2/app/controllers"
	"oauth2/app/global/db"
	"oauth2/app/repositories/mongo_repo"
	"oauth2/app/repositories/sql_repo"
	"oauth2/app/routes"
	"oauth2/app/usecases"
)

var connectionSet = wire.NewSet(
	db.NewPostgresClient,
	db.NewMongoClient,
	db.NewRedisClient,
)

var controllerSet = wire.NewSet(
	controllers.NewAuthController,
)

var useCaseSet = wire.NewSet(
	usecases.NewAuthUseCase,
)

var repositorySet = wire.NewSet(
	sql_repo.NewAuthRepository,
	mongo_repo.NewAccessTokenSessionRepository,
	mongo_repo.NewRefreshTokenSessionRepository,
)

func InitializeFiberServer(postgresParam db.PostgresParam, mongoParam db.MongoParam, redisParam db.RedisParam) *fiber.App {
	wire.Build(
		connectionSet,
		controllerSet,
		useCaseSet,
		repositorySet,
		routes.NewRouter,
	)
	return nil
}
