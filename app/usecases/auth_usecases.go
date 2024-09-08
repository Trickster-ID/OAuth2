package usecases

import (
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"oauth2/app/global/helper"
	"oauth2/app/global/model"
	"oauth2/app/models"
	"oauth2/app/repositories"
	"oauth2/app/requests"
	"oauth2/app/responses"
	"oauth2/app/security"
)

type IAuthUseCase interface {
	Login(req *requests.Login) (response *responses.Login, errLog *model.ErrorLog)
	ValidateUser(users *requests.Login) (*models.Users, *model.ErrorLog)
}

type authUseCase struct {
	authRepository repositories.IAuthRepository
}

func NewAuthUseCase(authRepository repositories.IAuthRepository) IAuthUseCase {
	return &authUseCase{
		authRepository: authRepository,
	}
}

func (u *authUseCase) Login(req *requests.Login) (response *responses.Login, errLog *model.ErrorLog) {
	user, errLog := u.ValidateUser(req)
	if errLog != nil {
		return
	}
	if user == nil {
		errLog = helper.WriteLog(errors.New("username or Email or Password is not valid"), 400, "")
		return
	}
	accessToken, refreshToken, err := security.GenerateToken(user.Username)
	if err != nil {
		errLog = helper.WriteLog(err, http.StatusInternalServerError, "")
		return
	}
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken
	return
}

func (u *authUseCase) ValidateUser(users *requests.Login) (*models.Users, *model.ErrorLog) {
	user, err := u.authRepository.GetUserByUsernameOrEmail(users.Username, users.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Trace("user is not found", err)
			return nil, nil
		}
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error get data")
		return nil, errLog
	}

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(users.Password))
	if err != nil {
		logrus.Trace("password is not valid", err)
		return nil, nil
	}
	return user, nil
}

func (u *authUseCase) ValidateToken(token string) (bool, error) {
	return true, nil
}
