package handler

import (
	"auth-service/internal/auth/dto"
	authService "auth-service/internal/auth/service"
	otpService "auth-service/internal/otp/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func VerifyOTP(c *fiber.Ctx) error {

	var req dto.VerifyOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	fmt.Printf("%+v\n", req)

	valid := otpService.VerifyOTP(
		req.Mobile,
		req.OTP,
	)

	if !valid {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid OTP",
		})
	}

	tenantCode := c.Locals("tenant_code").(string)

	countryCode := c.Locals("country_code").(string)

	otpService.DeleteOTP(req.Mobile)

	token, err := authService.Signup(
		tenantCode,
		countryCode,
		req.Name,
		req.Email,
		req.Mobile,
		req.Password,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Signup Successful",
		"access_token": token,
	})
}
