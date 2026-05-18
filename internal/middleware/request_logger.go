package middleware

import (
	"time"

	"catalog-service/internal/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {

		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		requestID, _ := c.Locals("request_id").(string)

		tenantCode, _ := c.Locals("tenant_code").(string)

		countryCode, _ := c.Locals("country_code").(string)

		logger.Log.Info(
			"incoming request",

			zap.String("request_id", requestID),

			zap.String("tenant_code", tenantCode),

			zap.String("country_code", countryCode),

			zap.String("method", c.Method()),

			zap.String("path", c.OriginalURL()),

			zap.Int("status", c.Response().StatusCode()),

			zap.String("ip", c.IP()),

			zap.String("user_agent", c.Get("User-Agent")),

			zap.Duration("duration", duration),
		)

		return err
	}
}
