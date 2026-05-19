package middleware

import (
	"auth-service/internal/logger"
	redisClient "auth-service/internal/redis"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

func OTPRateLimiter() fiber.Handler {

	logger.Log.Info(
		"otp rate limiter initialized",

		zap.Int(
			"max_requests",
			5,
		),

		zap.Duration(
			"expiration",
			10*time.Minute,
		),
	)

	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 10 * time.Minute,

		KeyGenerator: func(c *fiber.Ctx) string {

			start := time.Now()

			phone := c.FormValue("phone")

			if phone == "" {
				phone = c.IP()
			}

			key := "otp_limit:" + phone

			duration := time.Since(start)

			// =========================
			// SLOW OPERATION DETECTION
			// =========================

			if duration > time.Second {

				logger.Log.Warn(
					"slow limiter key generation detected",

					zap.Duration(
						"duration",
						duration,
					),

					zap.String(
						"operation",
						"OTPRateLimiter.KeyGenerator",
					),

					zap.String(
						"key",
						key,
					),
				)
			}

			logger.Log.Info(
				"otp limiter key generated",

				zap.String(
					"key",
					key,
				),

				zap.String(
					"ip_address",
					c.IP(),
				),
			)

			return key
		},

		Storage: redisClient.FiberStorage,

		LimitReached: func(c *fiber.Ctx) error {

			logger.Log.Warn(
				"otp rate limit exceeded",

				zap.String(
					"ip_address",
					c.IP(),
				),

				zap.String(
					"user_agent",
					c.Get("User-Agent"),
				),

				zap.String(
					"phone",
					c.FormValue("phone"),
				),
			)

			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{
					"success": false,
					"message": "Too many OTP requests",
				},
			)
		},
	})
}
