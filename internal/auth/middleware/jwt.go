package middleware

import (
	"auth-service/internal/logger"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Protected() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")

	logger.Log.Info(
		"jwt protected middleware initialized",
	)

	if secret == "" {

		logger.Log.Warn(
			"jwt secret is empty - middleware insecure",
		)
	} else {

		logger.Log.Info(
			"jwt secret loaded successfully",
		)
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(secret),
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			logger.Log.Warn(
				"unauthorized access attempt",

				zap.String(
					"ip_address",
					c.IP(),
				),

				zap.String(
					"path",
					c.Path(),
				),

				zap.String(
					"user_agent",
					c.Get("User-Agent"),
				),

				zap.Error(err),
			)

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "unauthorized",
			})
		},
	})
}
