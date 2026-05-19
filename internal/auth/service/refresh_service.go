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

	"github.com/golang-jwt/jwt/v5"

	"go.uber.org/zap"
)

func RefreshAccessToken(
	refreshToken string,
	deviceID string,
	ipAddress string,
	userAgent string,
) (
	string,
	string,
	error,
) {

	start := time.Now()

	logger.Log.Info(
		"refresh access token started",

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
	// Parse JWT
	// =========================

	parseStart := time.Now()

	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {

			return []byte(jwtService.GetJWTSecret()), nil
		},
	)

	parseDuration := time.Since(parseStart)

	if parseDuration > time.Second {

		logger.Log.Warn(
			"slow jwt parse detected",

			zap.Duration(
				"duration",
				parseDuration,
			),

			zap.String(
				"operation",
				"jwt.Parse",
			),
		)
	}

	if err != nil || !token.Valid {
		logger.Log.Warn(
			"invalid refresh token",

			zap.String(
				"device_id",
				deviceID,
			),

			zap.String(
				"ip_address",
				ipAddress,
			),

			zap.Error(err),
		)
		return "", "", errors.New("invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		logger.Log.Warn(
			"invalid token type",

			zap.String(
				"device_id",
				deviceID,
			),

			zap.String(
				"ip_address",
				ipAddress,
			),
		)
		return "", "", errors.New("invalid token type")
	}

	customerID := uint64(claims["customer_id"].(float64))

	tenantCode := claims["tenant_code"].(string)

	countryCode := claims["country_code"].(string)

	logger.Log.Info(
		"refresh token claims parsed",

		zap.Uint64(
			"customer_id",
			customerID,
		),

		zap.String(
			"tenant_code",
			tenantCode,
		),
	)

	// =========================
	// Hash Token
	// =========================

	hashStart := time.Now()

	hash := sha256.Sum256([]byte(refreshToken))

	tokenHash := hex.EncodeToString(hash[:])

	hashDuration := time.Since(hashStart)

	if hashDuration > time.Second {

		logger.Log.Warn(
			"slow hash computation detected",

			zap.Duration(
				"duration",
				hashDuration,
			),

			zap.String(
				"operation",
				"sha256",
			),
		)
	}

	// =========================
	// DB Lookup
	// =========================

	dbStart := time.Now()

	storedToken, err := repository.FindRefreshToken(
		tokenHash,
	)

	dbDuration := time.Since(dbStart)

	if dbDuration > time.Second {

		logger.Log.Warn(
			"slow database query detected",

			zap.Duration(
				"duration",
				dbDuration,
			),

			zap.String(
				"operation",
				"FindRefreshToken",
			),

			zap.Uint64(
				"customer_id",
				customerID,
			),
		)
	}

	if err != nil {
		logger.Log.Warn(
			"refresh token not found",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.Error(err),
		)
		return "", "", errors.New("refresh token not found")
	}

	// =========================
	// Expiry Check
	// =========================

	if time.Now().After(storedToken.ExpiresAt) {
		logger.Log.Warn(
			"refresh token expired",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.String(
				"token_id",
				storedToken.TokenFamilyID,
			),
		)
		return "", "", errors.New("refresh token expired")
	}

	// =========================
	// REPLAY ATTACK DETECTION
	// =========================

	if storedToken.Revoked {

		// Kill Entire Token Family

		logger.Log.Error(
			"refresh token replay detected",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.String(
				"family_id",
				storedToken.TokenFamilyID,
			),

			zap.String(
				"ip_address",
				ipAddress,
			),
		)

		_ = repository.RevokeTokenFamily(
			storedToken.TokenFamilyID,
		)

		return "", "", errors.New(
			"refresh token replay detected",
		)
	}

	// =========================
	// Revoke Old Token
	// =========================

	revokeStart := time.Now()

	if time.Now().After(
		storedToken.ExpiresAt,
	) {

		return "", "", errors.New(
			"refresh token expired",
		)
	}

	// Revoke Old Token (Rotation)

	err = repository.RevokeRefreshToken(
		storedToken.ID,
	)

	revokeDuration := time.Since(revokeStart)

	if revokeDuration > time.Second {

		logger.Log.Warn(
			"slow token revoke detected",

			zap.Duration(
				"duration",
				revokeDuration,
			),

			zap.Uint64(
				"customer_id",
				customerID,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to revoke refresh token",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.Error(err),
		)
		return "", "", err
	}

	// =========================
	// Generate New Tokens
	// =========================

	tokenGenStart := time.Now()

	newAccessToken, err := jwtService.GenerateToken(
		customerID,
		tenantCode,
		countryCode,
	)

	if err != nil {
		logger.Log.Error(
			"failed to generate access token",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.Error(err),
		)
		return "", "", err
	}

	newRefreshToken, err := jwtService.GenerateRefreshToken(
		customerID,
		tenantCode,
		countryCode,
	)

	tokenGenDuration := time.Since(tokenGenStart)

	if tokenGenDuration > time.Second {

		logger.Log.Warn(
			"slow token generation detected",

			zap.Duration(
				"duration",
				tokenGenDuration,
			),

			zap.Uint64(
				"customer_id",
				customerID,
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

			zap.Error(err),
		)
		return "", "", err
	}

	// =========================
	// Save New Refresh Token
	// =========================

	saveStart := time.Now()

	newHash := sha256.Sum256(
		[]byte(newRefreshToken),
	)

	newTokenHash := hex.EncodeToString(
		newHash[:],
	)

	refreshModel := model.RefreshToken{
		CustomerID:    customerID,
		TokenHash:     newTokenHash,
		TokenFamilyID: storedToken.TokenFamilyID,
		DeviceID:      deviceID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExpiresAt:     time.Now().Add(time.Hour * 24 * 30),
		Revoked:       false,
	}

	err = repository.CreateRefreshToken(
		&refreshModel,
	)

	saveDuration := time.Since(saveStart)

	if saveDuration > time.Second {

		logger.Log.Warn(
			"slow refresh token save detected",

			zap.Duration(
				"duration",
				saveDuration,
			),

			zap.Uint64(
				"customer_id",
				customerID,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to save refresh token",

			zap.Uint64(
				"customer_id",
				customerID,
			),

			zap.Error(err),
		)
		return "", "", err
	}

	totalDuration := time.Since(start)

	logger.Log.Info(
		"refresh token completed",

		zap.Uint64(
			"customer_id",
			customerID,
		),

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"ip_address",
			ipAddress,
		),

		zap.Duration(
			"total_duration",
			totalDuration,
		),
	)

	return newAccessToken,
		newRefreshToken,
		nil
}
