package service

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/auth/repository"
	jwtService "auth-service/internal/jwt"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	// Parse JWT

	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {

			return []byte(jwtService.GetJWTSecret()), nil
		},
	)

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	customerID := uint64(claims["customer_id"].(float64))

	tenantCode := claims["tenant_code"].(string)

	countryCode := claims["country_code"].(string)

	// Hash Token

	hash := sha256.Sum256([]byte(refreshToken))

	tokenHash := hex.EncodeToString(hash[:])

	// Check DB

	storedToken, err := repository.FindRefreshToken(
		tokenHash,
	)

	if err != nil {
		return "", "", errors.New("refresh token not found")
	}

	// Expiry Check

	if time.Now().After(storedToken.ExpiresAt) {

		return "", "", errors.New("refresh token expired")
	}

	// =========================
	// REPLAY ATTACK DETECTION
	// =========================

	if storedToken.Revoked {

		// Kill Entire Token Family

		_ = repository.RevokeTokenFamily(
			storedToken.TokenFamilyID,
		)

		return "", "", errors.New(
			"refresh token replay detected",
		)
	}

	// Expiry Check

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

	if err != nil {
		return "", "", err
	}

	// Generate New Tokens

	newAccessToken, err := jwtService.GenerateToken(
		customerID,
		tenantCode,
		countryCode,
	)

	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := jwtService.GenerateRefreshToken(
		customerID,
		tenantCode,
		countryCode,
	)

	if err != nil {
		return "", "", err
	}

	// Save New Refresh Token

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

	if err != nil {
		return "", "", err
	}

	return newAccessToken,
		newRefreshToken,
		nil
}
