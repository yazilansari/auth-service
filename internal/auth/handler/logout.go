package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {

	var req dto.LogoutRequest

	// Parse Request

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid request body",
			},
		)
	}

	// Validate Refresh Token

	if req.RefreshToken == "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Refresh token is required",
			},
		)
	}

	// Logout Service

	err := authService.Logout(
		req.RefreshToken,
	)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Logout failed",
			},
		)
	}

	// Success Response

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "Logged out successfully",
		},
	)
}
