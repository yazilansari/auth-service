package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"
	"auth-service/internal/logger"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Logout(c *fiber.Ctx) error {

	start := time.Now()

	logger.Log.Info(
		"logout request started",

		zap.String(
			"ip_address",
			c.IP(),
		),

		zap.String(
			"user_agent",
			c.Get("User-Agent"),
		),
	)

	var req dto.LogoutRequest

	// =========================
	// Parse Request
	// =========================

	parseStart := time.Now()

	if err := c.BodyParser(&req); err != nil {

		parseDuration := time.Since(parseStart)

		if parseDuration > time.Second {

			logger.Log.Warn(
				"slow logout request parsing detected",

				zap.Duration(
					"duration",
					parseDuration,
				),

				zap.String(
					"operation",
					"BodyParser",
				),
			)
		}

		logger.Log.Warn(
			"invalid logout request body",

			zap.String(
				"ip_address",
				c.IP(),
			),

			zap.Error(err),
		)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid request body",
			},
		)
	}

	// =========================
	// Validation
	// =========================

	if req.RefreshToken == "" {

		logger.Log.Warn(
			"missing refresh token in logout request",

			zap.String(
				"ip_address",
				c.IP(),
			),
		)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Refresh token is required",
			},
		)
	}

	logger.Log.Info(
		"logout attempt received",
	)

	// =========================
	// Logout Service
	// =========================

	logoutStart := time.Now()

	err := authService.Logout(
		req.RefreshToken,
	)

	logoutDuration := time.Since(logoutStart)

	if logoutDuration > time.Second {

		logger.Log.Warn(
			"slow logout operation detected",

			zap.Duration(
				"duration",
				logoutDuration,
			),

			zap.String(
				"operation",
				"Logout",
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"logout failed",

			zap.String(
				"ip_address",
				c.IP(),
			),

			zap.Error(err),
		)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Logout failed",
			},
		)
	}

	totalDuration := time.Since(start)

	logger.Log.Info(
		"logout successful",

		zap.String(
			"ip_address",
			c.IP(),
		),

		zap.Duration(
			"duration",
			totalDuration,
		),
	)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "Logged out successfully",
		},
	)
}
