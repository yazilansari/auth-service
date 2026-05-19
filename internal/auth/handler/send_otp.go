package handler

import (
	"auth-service/internal/auth/dto"
	"auth-service/internal/logger"
	otpService "auth-service/internal/otp/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func SendOTP(c *fiber.Ctx) error {

	start := time.Now()

	logger.Log.Info(
		"send otp request started",

		zap.String(
			"ip_address",
			c.IP(),
		),

		zap.String(
			"user_agent",
			c.Get("User-Agent"),
		),
	)

	var req dto.SendOTPRequest

	// =========================
	// Parse Request
	// =========================

	parseStart := time.Now()

	if err := c.BodyParser(&req); err != nil {
		parseDuration := time.Since(parseStart)

		// =========================
		// SLOW REQUEST PARSING
		// =========================

		if parseDuration > time.Second {

			logger.Log.Warn(
				"slow request body parsing detected",

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
			"invalid otp request body",

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

	// =========================
	// Validation
	// =========================

	if req.Mobile == "" {

		logger.Log.Warn(
			"mobile number missing in otp request",

			zap.String(
				"ip_address",
				c.IP(),
			),
		)

		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Mobile is required",
			},
		)
	}

	// =========================
	// Generate OTP
	// =========================

	generateOTPStart := time.Now()

	otp := otpService.GenerateOTP()

	generateOTPDuration := time.Since(generateOTPStart)

	if generateOTPDuration > time.Second {

		logger.Log.Warn(
			"slow otp generation detected",

			zap.Duration(
				"duration",
				generateOTPDuration,
			),

			zap.String(
				"operation",
				"GenerateOTP",
			),
		)
	}

	// =========================
	// Save OTP
	// =========================

	saveOTPStart := time.Now()

	err := otpService.SaveOTP(req.Mobile, otp)

	saveOTPDuration := time.Since(saveOTPStart)

	// =========================
	// SLOW OTP SAVE DETECTION
	// =========================

	if saveOTPDuration > time.Second {

		logger.Log.Warn(
			"slow otp save detected",

			zap.Duration(
				"duration",
				saveOTPDuration,
			),

			zap.String(
				"operation",
				"SaveOTP",
			),

			zap.String(
				"mobile",
				req.Mobile,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to save otp",

			zap.String(
				"mobile",
				req.Mobile,
			),

			zap.Duration(
				"duration",
				saveOTPDuration,
			),

			zap.Error(err),
		)
		return c.Status(500).JSON(fiber.Map{
			"message": "OTP save failed",
		})
	}

	// =========================
	// SMS Provider Placeholder
	// =========================

	logger.Log.Info(
		"sms provider integration pending",

		zap.String(
			"mobile",
			req.Mobile,
		),
	)

	totalDuration := time.Since(start)

	// =========================
	// SLOW REQUEST DETECTION
	// =========================

	if totalDuration > time.Second {

		logger.Log.Warn(
			"slow send otp request detected",

			zap.Duration(
				"duration",
				totalDuration,
			),

			zap.String(
				"mobile",
				req.Mobile,
			),

			zap.String(
				"ip_address",
				c.IP(),
			),
		)
	}

	logger.Log.Info(
		"otp sent successfully",

		zap.String(
			"mobile",
			req.Mobile,
		),

		zap.Duration(
			"duration",
			totalDuration,
		),

		zap.String(
			"ip_address",
			c.IP(),
		),
	)

	return c.JSON(fiber.Map{
		"message": "OTP Sent Successfully",
		"otp":     otp,
	})
}
