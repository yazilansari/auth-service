package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"

	"github.com/gofiber/fiber/v2"
)

func Refresh(c *fiber.Ctx) error {

	var req dto.RefreshRequest

	if err := c.BodyParser(&req); err != nil {

		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	deviceID := c.Get("X-Device-ID")

	ipAddress := c.IP()

	userAgent := c.Get("User-Agent")

	accessToken,
		refreshToken,
		err := authService.RefreshAccessToken(
		req.RefreshToken,
		deviceID,
		ipAddress,
		userAgent,
	)

	if err != nil {

		return c.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
