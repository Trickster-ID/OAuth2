package usecases

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"oauth2/app/dto"
	"oauth2/app/global/helper"
	mocksMongo "oauth2/app/mock/repositories/mongo_repo"
	mocksSql "oauth2/app/mock/repositories/sql_repo"
	"oauth2/app/models"

	"testing"
)

func TestLogin(t *testing.T) {
	goMockController := gomock.NewController(t)
	mockAuthRepo := mocksSql.NewMockIAuthRepository(goMockController)
	mockAccessRepo := mocksMongo.NewMockIAccessTokenSessionsRepository(goMockController)
	mockRefreshRepo := mocksMongo.NewMockIRefreshTokenSessionsRepository(goMockController)
	useCase := NewAuthUseCase(mockAuthRepo, mockAccessRepo, mockRefreshRepo)

	t.Run("+: admin", func(t *testing.T) {
		returnValue := &models.Users{
			Id:           1,
			Username:     "admin",
			Email:        "admin@live.com",
			PasswordHash: "$2a$04$ZP1.DVAdR677eHTBUpDzE.0hHnp31JcyRK/eMF9Z7Y.iOWJkE/JNi",
			RoleID:       1,
			CreatedAt:    nil,
			UpdatedAt:    nil,
		}
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("admin", "", context.Background()).Return(returnValue, nil)
		request := &dto.LoginRequest{
			Username: "admin",
			Email:    "",
			Password: "admin",
		}
		response, errLog := useCase.ValidateUser(request, context.Background())
		assert.Nil(t, errLog)
		assert.NotNil(t, response)
		assert.Equal(t, returnValue, response)
	})
}

func TestValidateUser(t *testing.T) {
	goMockController := gomock.NewController(t)
	mockAuthRepo := mocksSql.NewMockIAuthRepository(goMockController)
	mockAccessRepo := mocksMongo.NewMockIAccessTokenSessionsRepository(goMockController)
	mockRefreshRepo := mocksMongo.NewMockIRefreshTokenSessionsRepository(goMockController)
	useCase := NewAuthUseCase(mockAuthRepo, mockAccessRepo, mockRefreshRepo)

	t.Run("+: admin", func(t *testing.T) {
		returnValue := &models.Users{
			Id:           1,
			Username:     "admin",
			Email:        "admin@live.com",
			PasswordHash: "$2a$04$ZP1.DVAdR677eHTBUpDzE.0hHnp31JcyRK/eMF9Z7Y.iOWJkE/JNi",
			RoleID:       1,
			CreatedAt:    nil,
			UpdatedAt:    nil,
		}
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("admin", "", context.Background()).Return(returnValue, nil)
		request := &dto.LoginRequest{
			Username: "admin",
			Email:    "",
			Password: "admin",
		}
		response, errLog := useCase.ValidateUser(request, context.Background())
		assert.Nil(t, errLog)
		assert.NotNil(t, response)
		assert.Equal(t, returnValue, response)
	})

	t.Run("-: does not exist user", func(t *testing.T) {
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("userNotExist", "", context.Background()).Return(nil, helper.WriteLog(pgx.ErrNoRows, 404, "please enter valid username"))
		request := &dto.LoginRequest{
			Username: "userNotExist",
			Email:    "",
			Password: "userNotExist",
		}
		response, errLog := useCase.ValidateUser(request, context.Background())

		expectedError := helper.WriteLog(pgx.ErrNoRows, 404, "please enter valid username")
		assert.Nil(t, response)
		assert.NotNil(t, errLog)
		assert.Equal(t, expectedError, errLog)
	})
}
