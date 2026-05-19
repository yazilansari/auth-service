package repository

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/database"
	"auth-service/internal/logger"

	"time"

	"go.uber.org/zap"
)

func CreateRefreshToken(
	token *model.RefreshToken,
) error {

	// =========================
	// Query Start
	// =========================

	start := time.Now()

	logger.Log.Info(
		"creating refresh token",

		zap.Uint64(
			"customer_id",
			token.CustomerID,
		),

		zap.String(
			"family_id",
			token.TokenFamilyID,
		),
	)

	err := database.DB.Create(token).Error

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
				"CreateRefreshToken",
			),

			zap.Uint64(
				"customer_id",
				token.CustomerID,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to create refresh token",

			zap.Uint64(
				"customer_id",
				token.CustomerID,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"refresh token created",

		zap.Uint64(
			"customer_id",
			token.CustomerID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return nil
}

func FindRefreshToken(
	tokenHash string,
) (*model.RefreshToken, error) {

	// =========================
	// Query Start
	// =========================

	start := time.Now()

	logger.Log.Info(
		"finding refresh token",
	)

	var token model.RefreshToken

	err := database.DB.
		Where("token_hash = ?", tokenHash).
		First(&token).Error

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
				"FindRefreshToken",
			),
		)
	}

	if err != nil {
		logger.Log.Warn(
			"refresh token lookup failed",

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)
		return nil, err
	}

	logger.Log.Info(
		"refresh token found",

		zap.Uint64(
			"customer_id",
			token.CustomerID,
		),

		zap.String(
			"family_id",
			token.TokenFamilyID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return &token, nil
}

func RevokeRefreshToken(
	id uint64,
) error {

	// =========================
	// Query Start
	// =========================

	start := time.Now()

	logger.Log.Info(
		"revoking refresh token",

		zap.Uint64(
			"token_id",
			id,
		),
	)

	err := database.DB.
		Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error

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
				"RevokeRefreshToken",
			),

			zap.Uint64(
				"token_id",
				id,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to revoke refresh token",

			zap.Uint64(
				"token_id",
				id,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"refresh token revoked",

		zap.Uint64(
			"token_id",
			id,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return nil
}

func RevokeRefreshTokenByHash(
	tokenHash string,
) error {

	// =========================
	// Query Start
	// =========================

	start := time.Now()

	logger.Log.Info(
		"revoking refresh token by hash",
	)

	err := database.DB.
		Model(&model.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error

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
				"RevokeRefreshTokenByHash",
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to revoke refresh token by hash",

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"refresh token revoked by hash",

		zap.Duration(
			"duration",
			duration,
		),
	)

	return nil
}

func RevokeTokenFamily(
	familyID string,
) error {

	// =========================
	// Query Start
	// =========================

	start := time.Now()

	logger.Log.Info(
		"revoking token family",

		zap.String(
			"family_id",
			familyID,
		),
	)

	err := database.DB.
		Model(&model.RefreshToken{}).
		Where("token_family_id = ?", familyID).
		Update("revoked", true).Error

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
				"RevokeTokenFamily",
			),

			zap.String(
				"family_id",
				familyID,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to revoke token family",

			zap.String(
				"family_id",
				familyID,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"token family revoked",

		zap.String(
			"family_id",
			familyID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return nil
}
