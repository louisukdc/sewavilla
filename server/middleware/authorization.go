package middleware

import (
	"log"
	"net/http"
	"strings"
	"uas/database"
	"uas/model"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from the "Authorization" header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Authorization header is missing",
			})
		}

		// Check if the header contains "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid Authorization header format",
			})
		}

		// The token is expected in the form of "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the JWT token
		token, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Printf("Error parsing token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or expired token",
			})
		}

		// Extract the user ID from the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid token",
			})
		}

		// Find the user based on the token's claims (user ID)
		userID := claims["user_id"].(float64) // Assume user_id is in the token claims
		var user model.User
		if err := database.DBConn.First(&user, uint(userID)).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		}

		// Add the user to the context for use in handlers
		c.Locals("user", user)

		// Proceed to the next handler
		return c.Next()
	}
}
