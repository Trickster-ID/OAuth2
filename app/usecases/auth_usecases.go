package usecases

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"oauth2/app/dto"
	"oauth2/app/global/helper"
	"oauth2/app/global/model"
	"oauth2/app/models"
	"oauth2/app/repositories/mongo_repo"
	"oauth2/app/repositories/sql_repo"
	"oauth2/app/security"
	"os"
)

type IAuthUseCase interface {
	Login(req *dto.LoginRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog)
	ValidateUser(users *dto.LoginRequest, ctx context.Context) (*models.Users, *model.ErrorLog)
	RefreshToken(req *dto.RefreshTokenRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog)
}

type authUseCase struct {
	authRepository    sql_repo.IAuthRepository
	accessRepository  mongo_repo.IAccessTokenSessionsRepository
	refreshRepository mongo_repo.IRefreshTokenSessionsRepository
}

func NewAuthUseCase(authRepository sql_repo.IAuthRepository, accessRepository mongo_repo.IAccessTokenSessionsRepository, refreshRepository mongo_repo.IRefreshTokenSessionsRepository) IAuthUseCase {
	return &authUseCase{
		authRepository:    authRepository,
		accessRepository:  accessRepository,
		refreshRepository: refreshRepository,
	}
}

func (u *authUseCase) Login(req *dto.LoginRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog) {
	response := &dto.LoginResponse{}
	user, errLog := u.ValidateUser(req, ctx)
	if errLog != nil {
		return nil, errLog
	}
	if user == nil {
		errLog = helper.WriteLog(errors.New("username or Email or Password is not valid"), 400, "")
		return nil, errLog
	}
	userRequest := &models.UserDataOnJWT{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
	tokenResult, err := security.GenerateToken(userRequest)
	if err != nil {
		errLog = helper.WriteLog(err, http.StatusInternalServerError, "")
		return nil, errLog
	}
	requestAccessToken := &models.AccessTokenSession{
		AccessToken: tokenResult.AccessToken,
		Expired:     tokenResult.AccessTokenExpired,
		UserData:    userRequest,
	}
	go u.accessRepository.Insert(requestAccessToken, ctx)
	requestRefreshToken := &models.RefreshTokenSession{
		RefreshToken: tokenResult.RefreshToken,
		Expired:      tokenResult.RefreshTokenExpired,
		UserData:     userRequest,
	}
	go u.refreshRepository.Insert(requestRefreshToken, ctx)
	response.AccessToken = tokenResult.AccessToken
	response.RefreshToken = tokenResult.RefreshToken
	return response, nil
}

func (u *authUseCase) ValidateUser(request *dto.LoginRequest, ctx context.Context) (*models.Users, *model.ErrorLog) {
	user, errLog := u.authRepository.GetUserByUsernameOrEmail(request.Username, request.Email, ctx)
	if errLog != nil {
		return nil, errLog
	}

	// validate password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		errLog = helper.WriteLog(err, 401, "username or Email or Password is not valid")
		return nil, errLog
	}
	return user, nil
}

func (u *authUseCase) RefreshToken(req *dto.RefreshTokenRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog) {
	response := &dto.LoginResponse{
		RefreshToken: req.RefreshToken,
	}
	responseChan := make(chan *dto.GetByRefreshTokenChan)
	go func(responseChan chan *dto.GetByRefreshTokenChan) {
		res := &dto.GetByRefreshTokenChan{}
		res.Data, res.ErrLog = u.refreshRepository.GetByRefreshToken(req.RefreshToken, ctx)
		responseChan <- res
	}(responseChan)

	resultValidate := security.ValidateToken(req.RefreshToken, os.Getenv("JWT_KEY_REFRESH_TOKEN"))
	if resultValidate.Error != nil {
		if errors.Is(resultValidate.Error, jwt.ErrTokenExpired) {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, resultValidate.Error.Error())
			return nil, errLog
		}
		errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "invalid refresh token")
		return nil, errLog
	}
	responseGetSession := <-responseChan
	if responseGetSession.ErrLog != nil {
		return nil, responseGetSession.ErrLog
	}
	token, err := security.GenerateAccessToken(responseGetSession.Data.UserData)
	if err != nil {
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error generating access token")
		return nil, errLog
	}
	response.AccessToken = token.AccessToken

	requestRefreshToken := &models.AccessTokenSession{
		AccessToken: token.AccessToken,
		Expired:     token.AccessTokenExpired,
		UserData:    responseGetSession.Data.UserData,
	}
	go u.accessRepository.Insert(requestRefreshToken, ctx)
	return response, nil
}
