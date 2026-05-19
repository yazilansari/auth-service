package jwt

import (
	"os"
	"time"

	"auth-service/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func GenerateToken(
	customerID uint64,
	tenantCode string,
	countryCode string,
) (string, error) {

	start := time.Now()

	logger.Log.Info(
		"generating access token",

		zap.Uint64(
			"customer_id",
			customerID,
		),

		zap.String(
			"tenant_code",
			tenantCode,
		),
	)

	claims := jwt.MapClaims{
		"customer_id":  customerID,
		"tenant_code":  tenantCode,
		"country_code": countryCode,
		"type":         "access",
		"exp":          time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	duration := time.Since(start)

	if duration > time.Second {

		logger.Log.Warn(
			"slow access token generation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.String(
				"operation",
				"GenerateToken",
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to generate access token",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"access token generated successfully",

		zap.Uint64(
			"customer_id",
			customerID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return signedToken, nil
}

func GenerateRefreshToken(
	customerID uint64,
	tenantCode string,
	countryCode string,
) (string, error) {

	start := time.Now()

	logger.Log.Info(
		"generating refresh token",

		zap.Uint64(
			"customer_id",
			customerID,
		),

		zap.String(
			"tenant_code",
			tenantCode,
		),
	)

	claims := jwt.MapClaims{
		"customer_id":  customerID,
		"tenant_code":  tenantCode,
		"country_code": countryCode,
		"type":         "refresh",
		"exp": time.Now().
			Add(time.Hour * 24 * 30).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)

	duration := time.Since(start)

	if duration > time.Second {

		logger.Log.Warn(
			"slow refresh token generation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.String(
				"operation",
				"GenerateRefreshToken",
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to generate refresh token",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"refresh token generated successfully",

		zap.Uint64(
			"customer_id",
			customerID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return signedToken, nil
}
