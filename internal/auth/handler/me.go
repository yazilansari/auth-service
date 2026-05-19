package handler

import (
	"auth-service/internal/database"
	"auth-service/internal/logger"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Customer struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func Me(c *fiber.Ctx) error {

	start := time.Now()

	logger.Log.Info(
		"me request started",

		zap.String(
			"ip_address",
			c.IP(),
		),

		zap.String(
			"user_agent",
			c.Get("User-Agent"),
		),
	)

	// =========================
	// Extract JWT
	// =========================

	token := c.Locals("user").(*jwt.Token)

	// Extract claims
	claims := token.Claims.(jwt.MapClaims)

	// Get customer_id from token
	customerID := uint64(claims["customer_id"].(float64))

	logger.Log.Info(
		"jwt claims extracted",

		zap.Uint64(
			"customer_id",
			customerID,
		),
	)

	// =========================
	// DB Query
	// =========================

	dbStart := time.Now()

	// Customer response object
	var customer Customer

	// Query database
	err := database.DB.Table("customers").
		Select("id, name, email, phone").
		Where("id = ?", customerID).
		First(&customer).Error

	dbDuration := time.Since(dbStart)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if dbDuration > time.Second {

		logger.Log.Warn(
			"slow database query detected",

			zap.Duration(
				"duration",
				dbDuration,
			),

			zap.String(
				"operation",
				"Me.CustomerQuery",
			),

			zap.Uint64(
				"customer_id",
				customerID,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"customer not found",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.Error(err),
		)

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Customer not found",
		})
	}

	totalDuration := time.Since(start)

	logger.Log.Info(
		"me request completed",

		zap.Uint64(
			"customer_id",
			customerID,
		),

		zap.Duration(
			"total_duration",
			totalDuration,
		),
	)

	// Return customer data
	return c.JSON(customer)
}
