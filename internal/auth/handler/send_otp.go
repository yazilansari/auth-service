package handler

import (
	"auth-service/internal/auth/dto"
	otpService "auth-service/internal/otp/service"

	"github.com/gofiber/fiber/v2"
)

func SendOTP(c *fiber.Ctx) error {

	var req dto.SendOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	if req.Mobile == "" {

		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Mobile is required",
			},
		)
	}

	otp := otpService.GenerateOTP()

	err := otpService.SaveOTP(req.Mobile, otp)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "OTP save failed",
		})
	}

	// SMS Provider Integration Here

	return c.JSON(fiber.Map{
		"message": "OTP Sent Successfully",
		"otp":     otp,
	})
}
