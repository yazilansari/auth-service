package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"
	"auth-service/internal/logger"
	otpService "auth-service/internal/otp/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func VerifyOTP(c *fiber.Ctx) error {

	start := time.Now()

	logger.Log.Info(
		"verify otp request started",

		zap.String(
			"ip_address",
			c.IP(),
		),

		zap.String(
			"user_agent",
			c.Get("User-Agent"),
		),
	)

	var req dto.VerifyOTPRequest

	// =========================
	// Parse Request
	// =========================

	parseStart := time.Now()

	if err := c.BodyParser(&req); err != nil {

		parseDuration := time.Since(parseStart)

		if parseDuration > time.Second {

			logger.Log.Warn(
				"slow request parsing detected",

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
			"invalid verify otp request body",

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

	logger.Log.Info(
		"otp verification payload received",

		zap.String(
			"mobile",
			req.Mobile,
		),
	)

	// =========================
	// OTP Verification
	// =========================

	verifyStart := time.Now()

	valid := otpService.VerifyOTP(
		req.Mobile,
		req.OTP,
	)

	verifyDuration := time.Since(verifyStart)

	if verifyDuration > time.Second {

		logger.Log.Warn(
			"slow otp verification detected",

			zap.Duration(
				"duration",
				verifyDuration,
			),

			zap.String(
				"mobile",
				req.Mobile,
			),
		)
	}

	if !valid {
		logger.Log.Warn(
			"invalid otp provided",

			zap.String(
				"mobile",
				req.Mobile,
			),

			zap.String(
				"ip_address",
				c.IP(),
			),
		)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid OTP",
		})
	}

	logger.Log.Info(
		"otp verified successfully",

		zap.String(
			"mobile",
			req.Mobile,
		),
	)

	tenantCode := c.Locals("tenant_code").(string)

	countryCode := c.Locals("country_code").(string)

	// =========================
	// Delete OTP
	// =========================

	deleteStart := time.Now()

	otpService.DeleteOTP(req.Mobile)

	deleteDuration := time.Since(deleteStart)

	if deleteDuration > time.Second {

		logger.Log.Warn(
			"slow otp delete detected",

			zap.Duration(
				"duration",
				deleteDuration,
			),

			zap.String(
				"mobile",
				req.Mobile,
			),
		)
	}

	// =========================
	// Signup
	// =========================

	signupStart := time.Now()

	token, err := authService.Signup(
		tenantCode,
		countryCode,
		req.Name,
		req.Email,
		req.Mobile,
		req.Password,
	)

	signupDuration := time.Since(signupStart)

	if signupDuration > time.Second {

		logger.Log.Warn(
			"slow signup detected",

			zap.Duration(
				"duration",
				signupDuration,
			),

			zap.String(
				"mobile",
				req.Mobile,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"signup failed after otp verification",

			zap.String(
				"mobile",
				req.Mobile,
			),

			zap.Error(err),
		)
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	totalDuration := time.Since(start)

	logger.Log.Info(
		"otp verification + signup completed",

		zap.String(
			"mobile",
			req.Mobile,
		),

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.Duration(
			"total_duration",
			totalDuration,
		),
	)

	return c.JSON(fiber.Map{
		"message":      "Signup Successful",
		"access_token": token,
	})
}
