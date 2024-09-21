package controllers

import (
	"github.com/gofiber/fiber/v3"
	"oauth2/app/dto"
	"oauth2/app/global/helper"
	"oauth2/app/usecases"
)

type IAuthController interface {
	Login(f fiber.Ctx) error
	RefreshToken(f fiber.Ctx) error
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
	loginModel := new(dto.LoginRequest)
	if err := f.Bind().JSON(loginModel); err != nil {
		errLog := helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.ResponseError(f, errLog)
	}
	res, errLog := c.authUseCase.Login(loginModel, f.Context())
	if errLog != nil {
		return helper.ResponseError(f, errLog)
	}
	return helper.Response(f, res)
}

func (c *authController) RefreshToken(f fiber.Ctx) error {
	request := new(dto.RefreshTokenRequest)
	if err := f.Bind().JSON(request); err != nil {
		errLog := helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.ResponseError(f, errLog)
	}
	resultRefreshToken, errLog := c.authUseCase.RefreshToken(request, f.Context())
	if errLog != nil {
		return helper.ResponseError(f, errLog)
	}
	return helper.Response(f, resultRefreshToken)
}
