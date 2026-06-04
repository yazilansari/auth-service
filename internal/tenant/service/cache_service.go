package service

import (
	"auth-service/internal/logger"
	"auth-service/internal/redis"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

func SetTenantCache(
	key string,
	value interface{},
) error {

	logger.Log.Info(
		"setting tenant cache",

		zap.String(
			"cache_key",
			key,
		),
	)

	start := time.Now()

	data, err := json.Marshal(value)

	if err != nil {

		logger.Log.Error(
			"failed to marshal tenant cache data",

			zap.String(
				"cache_key",
				key,
			),

			zap.Error(err),
		)

		return err
	}

	err = redis.Client.Set(
		redis.Ctx,
		key,
		data,
		time.Hour,
	).Err()

	duration := time.Since(start)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"Redis.Get",
			),

			zap.String(
				"redis_key",
				key,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to store tenant cache",

			zap.String(
				"cache_key",
				key,
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
		"tenant cache stored successfully",

		zap.String(
			"cache_key",
			key,
		),

		zap.Duration(
			"ttl",
			time.Hour*24,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return nil
}

func GetTenantCache(
	key string,
	dest interface{},
) error {

	logger.Log.Info(
		"fetching tenant cache",

		zap.String(
			"cache_key",
			key,
		),
	)

	start := time.Now()

	data, err := redis.Client.Get(
		redis.Ctx,
		key,
	).Result()

	duration := time.Since(start)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"Redis.Get",
			),

			zap.String(
				"redis_key",
				key,
			),
		)
	}

	if err != nil {

		logger.Log.Warn(
			"tenant cache miss",

			zap.String(
				"cache_key",
				key,
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
		"tenant cache hit",

		zap.String(
			"cache_key",
			key,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	err = json.Unmarshal([]byte(data), dest)

	if err != nil {

		logger.Log.Error(
			"failed to unmarshal tenant cache",

			zap.String(
				"cache_key",
				key,
			),

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"tenant cache unmarshal success",

		zap.String(
			"cache_key",
			key,
		),
	)

	return nil
}
