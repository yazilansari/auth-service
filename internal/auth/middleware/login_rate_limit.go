package middleware

import (
	"auth-service/internal/logger"
	redisClient "auth-service/internal/redis"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

func LoginRateLimiter() fiber.Handler {

	logger.Log.Info(
		"login rate limiter initialized",

		zap.Int(
			"max_requests",
			10,
		),

		zap.Duration(
			"expiration",
			15*time.Minute,
		),
	)

	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 15 * time.Minute,

		KeyGenerator: func(c *fiber.Ctx) string {

			start := time.Now()

			ip := c.IP()
			key := "login_limit:" + ip

			duration := time.Since(start)

			// =========================
			// SLOW KEY GENERATION DETECTION
			// =========================

			if duration > time.Second {

				logger.Log.Warn(
					"slow login rate limiter key generation detected",

					zap.Duration(
						"duration",
						duration,
					),

					zap.String(
						"operation",
						"LoginRateLimiter.KeyGenerator",
					),

					zap.String(
						"ip_address",
						ip,
					),
				)
			}

			logger.Log.Info(
				"login rate limiter key generated",

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
				"login rate limit exceeded",

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
					"message": "Too many login attempts",
				},
			)
		},
	})
}
