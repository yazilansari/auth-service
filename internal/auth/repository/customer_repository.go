package repository

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/database"
	"auth-service/internal/logger"

	"time"

	"go.uber.org/zap"
)

func FindCustomerByPhone(phone string) (*model.Customer, error) {

	// =========================
	// Query Start
	// =========================

	start := time.Now()

	logger.Log.Info(
		"finding customer by phone",

		zap.String(
			"phone",
			phone,
		),
	)

	var customer model.Customer

	err := database.DB.
		Where("phone = ?", phone).
		First(&customer).Error

	duration := time.Since(start)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow database query detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"FindCustomerByPhone",
			),

			zap.String(
				"phone",
				phone,
			),
		)
	}

	if err != nil {
		logger.Log.Warn(
			"customer lookup failed",

			zap.String(
				"phone",
				phone,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)
		return nil, err
	}

	logger.Log.Info(
		"customer found",

		zap.Uint64(
			"customer_id",
			customer.ID,
		),

		zap.String(
			"tenant_code",
			customer.TenantCode,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return &customer, nil
}
