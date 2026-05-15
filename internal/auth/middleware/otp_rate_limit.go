package middleware

import (
	redisClient "auth-service/internal/redis"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func OTPRateLimiter() fiber.Handler {

	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 10 * time.Minute,

		KeyGenerator: func(c *fiber.Ctx) string {

			phone := c.FormValue("phone")

			if phone == "" {
				phone = c.IP()
			}

			return "otp_limit:" + phone
		},

		Storage: redisClient.FiberStorage,

		LimitReached: func(c *fiber.Ctx) error {

			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{
					"success": false,
					"message": "Too many OTP requests",
				},
			)
		},
	})
}
