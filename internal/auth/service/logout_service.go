package service

import (
	"auth-service/internal/auth/repository"
	"crypto/sha256"
	"encoding/hex"
)

func Logout(
	refreshToken string,
) error {

	// Hash Refresh Token

	hash := sha256.Sum256(
		[]byte(refreshToken),
	)

	tokenHash := hex.EncodeToString(
		hash[:],
	)

	// Revoke Token

	return repository.RevokeRefreshTokenByHash(
		tokenHash,
	)
}
