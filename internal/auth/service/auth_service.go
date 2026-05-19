package service

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/auth/repository"
	jwtService "auth-service/internal/jwt"
	"auth-service/internal/logger"

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Login(
	phone string,
	password string,
	deviceID string,
	ipAddress string,
	userAgent string,
) (
	string,
	string,
	error,
) {

	logger.Log.Info(
		"login started",

		zap.String(
			"phone",
			phone,
		),

		zap.String(
			"device_id",
			deviceID,
		),

		zap.String(
			"ip_address",
			ipAddress,
		),
	)

	// =========================
	// Find Customer
	// =========================

	findCustomerStart := time.Now()

	customer, err := repository.FindCustomerByPhone(phone)

	findCustomerDuration := time.Since(findCustomerStart)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if findCustomerDuration > time.Second {

		logger.Log.Warn(
			"slow database query detected",

			zap.Duration(
				"duration",
				findCustomerDuration,
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

			zap.Error(err),
		)
		return "", "", err
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
	)

	// =========================
	// Password Verification
	// =========================

	passwordCheckStart := time.Now()

	err = bcrypt.CompareHashAndPassword(
		[]byte(customer.Password),
		[]byte(password),
	)

	passwordCheckDuration := time.Since(passwordCheckStart)

	// =========================
	// SLOW PASSWORD CHECK
	// =========================

	if passwordCheckDuration > time.Second {

		logger.Log.Warn(
			"slow password verification detected",

			zap.Duration(
				"duration",
				passwordCheckDuration,
			),

			zap.Uint64(
				"customer_id",
				customer.ID,
			),
		)
	}

	if err != nil {
		logger.Log.Warn(
			"invalid credentials",

			zap.Uint64(
				"customer_id",
				customer.ID,
			),

			zap.String(
				"phone",
				phone,
			),

			zap.String(
				"ip_address",
				ipAddress,
			),
		)
		return "", "", errors.New("invalid credentials")
	}

	// =========================
	// Generate Access Token
	// =========================

	accessTokenStart := time.Now()

	accessToken, err := jwtService.GenerateToken(
		customer.ID,
		customer.TenantCode,
		customer.CountryCode,
	)

	accessTokenDuration := time.Since(accessTokenStart)

	if accessTokenDuration > time.Second {

		logger.Log.Warn(
			"slow access token generation detected",

			zap.Duration(
				"duration",
				accessTokenDuration,
			),

			zap.Uint64(
				"customer_id",
				customer.ID,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to generate access token",

			zap.Uint64(
				"customer_id",
				customer.ID,
			),

			zap.Error(err),
		)
		return "", "", err
	}

	// =========================
	// Generate Refresh Token
	// =========================

	refreshTokenStart := time.Now()

	refreshToken, err := jwtService.GenerateRefreshToken(
		customer.ID,
		customer.TenantCode,
		customer.CountryCode,
	)

	refreshTokenDuration := time.Since(refreshTokenStart)

	if refreshTokenDuration > time.Second {

		logger.Log.Warn(
			"slow refresh token generation detected",

			zap.Duration(
				"duration",
				refreshTokenDuration,
			),

			zap.Uint64(
				"customer_id",
				customer.ID,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to generate refresh token",

			zap.Uint64(
				"customer_id",
				customer.ID,
			),

			zap.Error(err),
		)
		return "", "", err
	}

	// Create Token Family

	familyID := uuid.NewString()

	logger.Log.Info(
		"token family created",

		zap.Uint64(
			"customer_id",
			customer.ID,
		),

		zap.String(
			"family_id",
			familyID,
		),
	)

	// =========================
	// Store Refresh Token
	// =========================

	storeTokenStart := time.Now()

	// Store Refresh Token

	hash := sha256.Sum256([]byte(refreshToken))

	tokenHash := hex.EncodeToString(hash[:])

	refreshModel := model.RefreshToken{
		CustomerID:    customer.ID,
		TokenHash:     tokenHash,
		TokenFamilyID: familyID,
		DeviceID:      deviceID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExpiresAt:     time.Now().Add(time.Hour * 24 * 30),
		Revoked:       false,
	}

	err = repository.CreateRefreshToken(
		&refreshModel,
	)

	storeTokenDuration := time.Since(storeTokenStart)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if storeTokenDuration > time.Second {

		logger.Log.Warn(
			"slow database query detected",

			zap.Duration(
				"duration",
				storeTokenDuration,
			),

			zap.String(
				"operation",
				"CreateRefreshToken",
			),

			zap.Uint64(
				"customer_id",
				customer.ID,
			),

			zap.String(
				"tenant_code",
				customer.TenantCode,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to store refresh token",

			zap.Uint64(
				"customer_id",
				customer.ID,
			),

			zap.Error(err),
		)
		return "", "", err
	}

	logger.Log.Info(
		"login successful",

		zap.Uint64(
			"customer_id",
			customer.ID,
		),

		zap.String(
			"tenant_code",
			customer.TenantCode,
		),

		zap.String(
			"device_id",
			deviceID,
		),

		zap.String(
			"ip_address",
			ipAddress,
		),
	)

	return accessToken,
		refreshToken,
		nil
}
