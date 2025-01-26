package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler centralizes error handling
func ErrorHandler(c *fiber.Ctx, err error) error {
	// If there's an error, return a 500 internal server error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
			"details": err.Error(),
		})
	}
	// If no error, continue the request
	return c.Next()
}
