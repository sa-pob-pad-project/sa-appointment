package middleware

import (
	"appointment-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware(jwtService *jwt.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get cookie named "access_token"
		token := c.Cookies("access_token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		claims, err := jwtService.Parse(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
