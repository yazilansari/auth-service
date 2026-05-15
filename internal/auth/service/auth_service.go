package service

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/auth/repository"
	jwtService "auth-service/internal/jwt"

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
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

	customer, err := repository.FindCustomerByPhone(phone)

	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(customer.Password),
		[]byte(password),
	)

	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := jwtService.GenerateToken(
		customer.ID,
		customer.TenantCode,
		customer.CountryCode,
	)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwtService.GenerateRefreshToken(
		customer.ID,
		customer.TenantCode,
		customer.CountryCode,
	)

	if err != nil {
		return "", "", err
	}

	// Create Token Family

	familyID := uuid.NewString()

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

	if err != nil {
		return "", "", err
	}

	return accessToken,
		refreshToken,
		nil
}
