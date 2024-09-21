package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"oauth2/app/global/helper"
	"oauth2/app/security"
	"os"
	"strings"
)

func TokenMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		errLog := helper.WriteLogWoP(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
		if authHeader == "" {
			logrus.Trace("auth header is empty")
			return helper.ResponseError(c, errLog)
		}

		// Check if the token is in the format "Bearer <token>"
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		response := security.ValidateToken(tokenString, os.Getenv("JWT_KEY_ACCESS_TOKEN"))
		if response.Error != nil {
			if errors.Is(response.Error, jwt.ErrTokenExpired) {
				errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, response.Error.Error())
				return helper.ResponseError(c, errLog)
			}
			logrus.Trace(response.Error)
			return helper.ResponseError(c, errLog)
		}
		marshalledData, err := json.Marshal(response.User)
		if err != nil {
			logrus.Trace(err)
			return helper.ResponseError(c, errLog)
		}
		c.Set("user_data", string(marshalledData))
		return c.Next()
	}
}

func BasicAuthMiddleware() fiber.Handler {
	return func(f fiber.Ctx) error {
		// Get the Authorization header
		auth := f.Get("Authorization")
		if auth == "" {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "No Authorization header")
			return helper.ResponseError(f, errLog)
		}

		// Check if it starts with "Basic"
		if !strings.HasPrefix(auth, "Basic ") {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid Authorization header")
			return helper.ResponseError(f, errLog)
		}

		// Decode the Base64 encoded credentials
		encodedCredentials := strings.TrimPrefix(auth, "Basic ")
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid credentials encoding")
			return helper.ResponseError(f, errLog)
		}

		// Split the decoded credentials into username and password
		credentials := strings.SplitN(string(decodedCredentials), ":", 2)
		if len(credentials) != 2 || credentials[0] != os.Getenv("BASICAUTH_USERNAME") || credentials[1] != os.Getenv("BASICAUTH_PASSWORD") {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid username or password")
			return helper.ResponseError(f, errLog)
		}

		// Continue to the next handler if authentication is successful
		return f.Next()
	}
}
