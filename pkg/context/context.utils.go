package context

import (
	"appointment-service/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func WithBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			return response.BadRequest(c, "Cannot parse JSON")
		}
		validate := validator.New()
		if err := validate.Struct(body); err != nil {
			return response.BadRequest(c, err.Error())
		}
		c.Locals("body", body)
		return c.Next()
	}
}
