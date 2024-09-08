package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// getToken will generate token with JWT lib
func getToken(username string, expiredTime time.Time, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"username":   username,
		"expires_at": expiredTime.Unix(),
		"issuer":     "oauth2-pikri",
		"issued_at":  time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token, err := t.SignedString([]byte(secretKey))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return token, nil
}

func GenerateToken(username string) (string, string, error) {
	accessToken, err := getToken(username, time.Now().Add(time.Hour*24), os.Getenv("JWT_KEY_ACCESS_TOKEN"))
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}
	refreshToken, err := getToken(username, time.Now().Add(time.Hour*24*7), os.Getenv("JWT_KEY_REFRESH_TOKEN"))
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// ValidateToken will validate jwt token for middleware function
func ValidateToken(tokenString string) (valid bool, err error, username string) {
	// Parse the token and validate its signature
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			// Make sure the signing method is ECDSA (ES256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return os.Getenv("JWT_KEY_ACCESS_TOKEN"), nil
		},
	)
	if err != nil || !token.Valid {
		return
	}
	expired := token.Claims.(jwt.MapClaims)["expires_at"].(int64)
	if time.Now().Unix() > expired {
		err = errors.New("token is expired")
		return
	}
	valid = true
	username = token.Claims.(jwt.MapClaims)["username"].(string)
	return
}

func RefreshToken(refreshToken string) (string, error) {
	valid, err, username := ValidateToken(refreshToken)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if !valid {
		err := errors.New("refresh token is invalid")
		logrus.Error(err)
		return "", err
	}
	newToken, err := getToken(username, time.Now().Add(time.Hour*24), os.Getenv("JWT_KEY_ACCESS_TOKEN"))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newToken, nil
}
