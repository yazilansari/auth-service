package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {

	var req dto.LoginRequest

	// Parse Request

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid request body",
			},
		)
	}

	// Basic Validation

	if req.Mobile == "" || req.Password == "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Mobile and password are required",
			},
		)
	}

	// Device Info

	deviceID := c.Get("X-Device-ID")

	ipAddress := c.IP()

	userAgent := c.Get("User-Agent")

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

		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

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
