package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"oauth2/app/constants"
	"oauth2/app/dto"
	"oauth2/app/models"
	"os"
	"time"
)

// getToken will generate token with JWT lib
func getToken(userData *models.UserDataOnJWT, expiredTime time.Time, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"id":       userData.Id,
		"username": userData.Username,
		"email":    userData.Email,
		"exp":      expiredTime.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "oauth2-pikri",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(secretKey))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return token, nil
}

func GenerateToken(userData *models.UserDataOnJWT) (*dto.GenerateTokenResponse, error) {
	response := &dto.GenerateTokenResponse{
		AccessTokenExpired:  time.Now().Add(constants.ACCESS_TOKEN_EXPIRED),
		RefreshTokenExpired: time.Now().Add(constants.REFRESH_TOKEN_EXPIRED),
	}
	accessToken, err := getToken(userData, response.AccessTokenExpired, os.Getenv("JWT_KEY_ACCESS_TOKEN"))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	response.AccessToken = accessToken
	refreshToken, err := getToken(userData, response.RefreshTokenExpired, os.Getenv("JWT_KEY_REFRESH_TOKEN"))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	response.RefreshToken = refreshToken
	return response, nil
}

func GenerateAccessToken(userData *models.UserDataOnJWT) (*dto.GenerateAccessTokenResponse, error) {
	response := &dto.GenerateAccessTokenResponse{
		AccessTokenExpired: time.Now().Add(constants.ACCESS_TOKEN_EXPIRED),
	}
	accessToken, err := getToken(userData, response.AccessTokenExpired, os.Getenv("JWT_KEY_ACCESS_TOKEN"))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	response.AccessToken = accessToken
	return response, nil
}

// ValidateToken will validate jwt token for middleware function
func ValidateToken(tokenString, secretKey string) *dto.ValidateTokenResponse {
	response := &dto.ValidateTokenResponse{}
	// Parse the token and validate its signature
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			// Make sure the signing method is ECDSA (ES256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil || !token.Valid {
		response.Error = errors.New("invalid token")
		return response
	}
	expTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		response.Error = err
		return response
	}
	if expTime.Before(time.Now()) {
		response.Error = jwt.ErrTokenExpired
		return response
	}
	userData := &models.UserDataOnJWT{
		Id:       int64(token.Claims.(jwt.MapClaims)["id"].(float64)),
		Username: token.Claims.(jwt.MapClaims)["username"].(string),
		Email:    token.Claims.(jwt.MapClaims)["email"].(string),
	}
	response.User = userData
	return response
}

func RefreshToken(refreshToken string) (string, error) {
	response := ValidateToken(refreshToken, os.Getenv("JWT_KEY_REFRESH_TOKEN"))
	if response.Error != nil {
		logrus.Error(response.Error)
		return "", response.Error
	}
	newToken, err := getToken(response.User, time.Now().Add(time.Hour*24), os.Getenv("JWT_KEY_ACCESS_TOKEN"))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newToken, nil
}
