package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"
	"auth-service/internal/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Login(c *fiber.Ctx) error {

	var req dto.LoginRequest

	// Device Info

	deviceID := c.Get("X-Device-ID")

	ipAddress := c.IP()

	userAgent := c.Get("User-Agent")

	requestID := c.Locals("requestid")

	// Parse Request

	if err := c.BodyParser(&req); err != nil {

		logger.Log.Warn("Invalid login request body",
			zap.String("ip", ipAddress),
			zap.String("user_agent", userAgent),
			zap.Any("request_id", requestID),
			zap.Error(err),
		)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid request body",
			},
		)
	}

	// Basic Validation

	if req.Mobile == "" || req.Password == "" {

		logger.Log.Warn("Login validation failed",
			zap.String("mobile", req.Mobile),
			zap.String("ip", ipAddress),
			zap.Any("request_id", requestID),
		)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Mobile and password are required",
			},
		)
	}

	// Login Attempt

	logger.Log.Info("Login attempt",
		zap.String("mobile", req.Mobile),
		zap.String("device_id", deviceID),
		zap.String("ip", ipAddress),
		zap.String("user_agent", userAgent),
		zap.Any("request_id", requestID),
	)

	// Login Service

	accessToken,
		refreshToken,
		err := authService.Login(
		req.Mobile,
		req.Password,
		deviceID,
		ipAddress,
		userAgent,
	)

	if err != nil {

		logger.Log.Warn("Login failed",
			zap.String("mobile", req.Mobile),
			zap.String("ip", ipAddress),
			zap.String("device_id", deviceID),
			zap.Any("request_id", requestID),
			zap.Error(err),
		)

		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	// Success Log

	logger.Log.Info("Login successful",
		zap.String("mobile", req.Mobile),
		zap.String("device_id", deviceID),
		zap.String("ip", ipAddress),
		zap.Any("request_id", requestID),
	)

	// Success Response

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,

			"message": "Login successful",

			"data": fiber.Map{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"token_type":    "Bearer",
				"expires_in":    900,
			},
		},
	)
}
