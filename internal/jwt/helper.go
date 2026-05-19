package jwt

import (
	"os"

	"auth-service/internal/logger"

	"go.uber.org/zap"
)

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")

	// =========================
	// SECRET VALIDATION LOGGING
	// =========================

	if secret == "" {

		logger.Log.Error(
			"JWT_SECRET is not set in environment variables",
		)
	} else {

		logger.Log.Info(
			"JWT_SECRET loaded successfully",
			zap.Int(
				"length",
				len(secret),
			),
		)
	}

	return secret
}
