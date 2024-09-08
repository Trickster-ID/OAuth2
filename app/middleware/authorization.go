package middleware

import (
	"github.com/gofiber/fiber/v3"
	"oauth2/app/security"
	"strings"
)

func Authorization() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		// Check if the token is in the format "Bearer <token>"
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		valid, err, _ := security.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
		}
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		return c.Next()
	}
}
