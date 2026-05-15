package middleware

import (
	redisClient "auth-service/internal/redis"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func LoginRateLimiter() fiber.Handler {

	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 15 * time.Minute,

		KeyGenerator: func(c *fiber.Ctx) string {

			return "login_limit:" + c.IP()
		},

		Storage: redisClient.FiberStorage,

		LimitReached: func(c *fiber.Ctx) error {

			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{
					"success": false,
					"message": "Too many login attempts",
				},
			)
		},
	})
}
