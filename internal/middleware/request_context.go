package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestContextMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {

		requestID := c.Get("X-Request-ID")

		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("X-Request-ID", requestID)

		c.Locals("request_id", requestID)

		c.Locals("ip_address", c.IP())

		c.Locals("user_agent", c.Get("User-Agent"))

		return c.Next()
	}
}
