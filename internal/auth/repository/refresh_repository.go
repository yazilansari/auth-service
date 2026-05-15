package repository

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/database"
)

func CreateRefreshToken(
	token *model.RefreshToken,
) error {

	return database.DB.Create(token).Error
}

func FindRefreshToken(
	tokenHash string,
) (*model.RefreshToken, error) {

	var token model.RefreshToken

	err := database.DB.
		Where("token_hash = ?", tokenHash).
		First(&token).Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func RevokeRefreshToken(
	id uint64,
) error {

	return database.DB.
		Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func RevokeRefreshTokenByHash(
	tokenHash string,
) error {

	return database.DB.
		Model(&model.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error
}

func RevokeTokenFamily(
	familyID string,
) error {

	return database.DB.
		Model(&model.RefreshToken{}).
		Where("token_family_id = ?", familyID).
		Update("revoked", true).Error
}
