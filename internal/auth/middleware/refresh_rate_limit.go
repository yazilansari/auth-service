package middleware

import (
	"auth-service/internal/logger"
	redisClient "auth-service/internal/redis"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

func RefreshRateLimiter() fiber.Handler {

	logger.Log.Info(
		"refresh rate limiter initialized",

		zap.Int(
			"max_requests",
			20,
		),

		zap.Duration(
			"expiration",
			15*time.Minute,
		),
	)

	return limiter.New(limiter.Config{
		Max:        20,
		Expiration: 15 * time.Minute,

		KeyGenerator: func(c *fiber.Ctx) string {

			start := time.Now()

			ip := c.IP()
			key := "refresh_limit:" + ip

			duration := time.Since(start)

			// =========================
			// SLOW KEY GENERATION DETECTION
			// =========================

			if duration > time.Second {

				logger.Log.Warn(
					"slow refresh rate limiter key generation detected",

					zap.Duration(
						"duration",
						duration,
					),

					zap.String(
						"operation",
						"RefreshRateLimiter.KeyGenerator",
					),

					zap.String(
						"ip_address",
						ip,
					),
				)
			}

			logger.Log.Info(
				"refresh rate limiter key generated",

				zap.String(
					"key",
					key,
				),

				zap.String(
					"ip_address",
					ip,
				),
			)

			return key
		},

		Storage: redisClient.FiberStorage,

		LimitReached: func(c *fiber.Ctx) error {

			logger.Log.Warn(
				"refresh rate limit exceeded",

				zap.String(
					"ip_address",
					c.IP(),
				),

				zap.String(
					"user_agent",
					c.Get("User-Agent"),
				),

				zap.String(
					"path",
					c.Path(),
				),
			)

			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{
					"success": false,
					"message": "Too many refresh attempts",
				},
			)
		},
	})
}
