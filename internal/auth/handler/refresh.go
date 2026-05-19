package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"
	"auth-service/internal/logger"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Refresh(c *fiber.Ctx) error {

	start := time.Now()

	logger.Log.Info(
		"refresh token request started",

		zap.String(
			"ip_address",
			c.IP(),
		),

		zap.String(
			"user_agent",
			c.Get("User-Agent"),
		),
	)

	var req dto.RefreshRequest

	// =========================
	// Parse Request
	// =========================

	parseStart := time.Now()

	if err := c.BodyParser(&req); err != nil {

		parseDuration := time.Since(parseStart)

		if parseDuration > time.Second {

			logger.Log.Warn(
				"slow refresh request parsing detected",

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
			"invalid refresh request body",

			zap.String(
				"ip_address",
				c.IP(),
			),

			zap.Error(err),
		)

		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	deviceID := c.Get("X-Device-ID")

	ipAddress := c.IP()

	userAgent := c.Get("User-Agent")

	logger.Log.Info(
		"refresh attempt received",

		zap.String(
			"device_id",
			deviceID,
		),

		zap.String(
			"ip_address",
			ipAddress,
		),
	)

	// =========================
	// Refresh Token Service
	// =========================

	refreshStart := time.Now()

	accessToken,
		refreshToken,
		err := authService.RefreshAccessToken(
		req.RefreshToken,
		deviceID,
		ipAddress,
		userAgent,
	)

	refreshDuration := time.Since(refreshStart)

	// =========================
	// SLOW REFRESH DETECTION
	// =========================

	if refreshDuration > time.Second {

		logger.Log.Warn(
			"slow refresh token operation detected",

			zap.Duration(
				"duration",
				refreshDuration,
			),

			zap.String(
				"device_id",
				deviceID,
			),

			zap.String(
				"ip_address",
				ipAddress,
			),
		)
	}

	if err != nil {

		logger.Log.Warn(
			"refresh token failed",

			zap.String(
				"device_id",
				deviceID,
			),

			zap.String(
				"ip_address",
				ipAddress,
			),

			zap.Error(err),
		)

		return c.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	totalDuration := time.Since(start)

	logger.Log.Info(
		"refresh token successful",

		zap.String(
			"device_id",
			deviceID,
		),

		zap.String(
			"ip_address",
			ipAddress,
		),

		zap.Duration(
			"total_duration",
			totalDuration,
		),
	)

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
