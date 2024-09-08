package controllers

import (
	"github.com/gofiber/fiber/v3"
	"oauth2/app/global/helper"
	"oauth2/app/requests"
	"oauth2/app/usecases"
)

type IAuthController interface {
	Login(f fiber.Ctx) error
}

type authController struct {
	authUseCase usecases.IAuthUseCase
}

func NewAuthController(authUseCase usecases.IAuthUseCase) IAuthController {
	return &authController{
		authUseCase: authUseCase,
	}
}

func (c *authController) Login(f fiber.Ctx) error {
	loginModel := new(requests.Login)
	if err := f.Bind().JSON(loginModel); err != nil {
		errLog := helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.ResponseError(f, errLog)
	}
	res, errLog := c.authUseCase.Login(loginModel)
	if errLog != nil {
		return helper.ResponseError(f, errLog)
	}
	return helper.Response(f, res)
}
