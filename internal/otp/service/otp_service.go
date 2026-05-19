package service

import (
	"auth-service/internal/logger"
	"auth-service/internal/queue"
	"auth-service/internal/redis"

	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

func GenerateOTP() string {
	start := time.Now()

	rand.Seed(time.Now().UnixNano())

	otp := fmt.Sprintf("%04d", rand.Intn(9000)+1000)

	duration := time.Since(start)

	// =========================
	// SLOW OTP GENERATION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow otp generation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"GenerateOTP",
			),
		)
	}

	logger.Log.Info(
		"otp generated",

		zap.Duration(
			"duration",
			duration,
		),
	)

	return otp
}

func SaveOTP(phone string, otp string) error {
	start := time.Now()

	key := "otp:" + phone

	logger.Log.Info(
		"saving otp to redis",

		zap.String(
			"phone",
			phone,
		),

		zap.String(
			"redis_key",
			key,
		),
	)

	// =========================
	// Save OTP To Redis
	// =========================

	redisStart := time.Now()

	err := redis.Client.Set(
		redis.Ctx,
		key,
		otp,
		time.Minute*2,
	).Err()

	redisDuration := time.Since(redisStart)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if redisDuration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				redisDuration,
			),

			zap.String(
				"operation",
				"Redis.Set",
			),

			zap.String(
				"redis_key",
				key,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to save otp in redis",

			zap.String(
				"phone",
				phone,
			),

			zap.Duration(
				"duration",
				redisDuration,
			),

			zap.Error(err),
		)
		return err
	}

	logger.Log.Info(
		"otp saved to redis",

		zap.String(
			"phone",
			phone,
		),

		zap.Duration(
			"duration",
			redisDuration,
		),
	)

	// =========================
	// Push Queue Job
	// =========================

	queueStart := time.Now()

	err = queue.EnqueueOTP(
		phone,
		otp,
	)

	queueDuration := time.Since(queueStart)

	// =========================
	// SLOW QUEUE DETECTION
	// =========================

	if queueDuration > time.Second {

		logger.Log.Warn(
			"slow queue operation detected",

			zap.Duration(
				"duration",
				queueDuration,
			),

			zap.String(
				"operation",
				"EnqueueOTP",
			),

			zap.String(
				"phone",
				phone,
			),
		)
	}

	if err != nil {
		logger.Log.Error(
			"failed to enqueue otp job",

			zap.String(
				"phone",
				phone,
			),

			zap.Duration(
				"duration",
				queueDuration,
			),

			zap.Error(err),
		)
		return err
	}

	totalDuration := time.Since(start)

	// =========================
	// SLOW SAVE OTP DETECTION
	// =========================

	if totalDuration > time.Second {

		logger.Log.Warn(
			"slow save otp operation detected",

			zap.Duration(
				"duration",
				totalDuration,
			),

			zap.String(
				"phone",
				phone,
			),
		)
	}

	logger.Log.Info(
		"otp save completed",

		zap.String(
			"phone",
			phone,
		),

		zap.Duration(
			"duration",
			totalDuration,
		),
	)

	return nil
}

func VerifyOTP(phone string, otp string) bool {
	start := time.Now()

	key := "otp:" + phone

	logger.Log.Info(
		"verifying otp",

		zap.String(
			"phone",
			phone,
		),

		zap.String(
			"redis_key",
			key,
		),
	)

	// =========================
	// Redis Get
	// =========================

	redisStart := time.Now()

	storedOTP, err := redis.Client.Get(
		redis.Ctx,
		key,
	).Result()

	redisDuration := time.Since(redisStart)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if redisDuration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				redisDuration,
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
			"otp verification failed",

			zap.String(
				"phone",
				phone,
			),

			zap.Duration(
				"duration",
				redisDuration,
			),

			zap.Error(err),
		)
		return false
	}

	isValid := storedOTP == otp

	totalDuration := time.Since(start)

	// =========================
	// SLOW OTP VERIFICATION
	// =========================

	if totalDuration > time.Second {

		logger.Log.Warn(
			"slow otp verification detected",

			zap.Duration(
				"duration",
				totalDuration,
			),

			zap.String(
				"phone",
				phone,
			),
		)
	}

	if isValid {

		logger.Log.Info(
			"otp verified successfully",

			zap.String(
				"phone",
				phone,
			),

			zap.Duration(
				"duration",
				totalDuration,
			),
		)

	} else {

		logger.Log.Warn(
			"invalid otp entered",

			zap.String(
				"phone",
				phone,
			),

			zap.Duration(
				"duration",
				totalDuration,
			),
		)
	}

	return isValid
}

func DeleteOTP(phone string) {
	start := time.Now()

	logger.Log.Info(
		"deleting otp",

		zap.String(
			"phone",
			phone,
		),
	)

	key := "otp:" + phone

	err := redis.Client.Del(redis.Ctx, key).Err()

	duration := time.Since(start)

	// =========================
	// SLOW REDIS DELETE
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow redis delete detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"Redis.Del",
			),

			zap.String(
				"key",
				key,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to delete otp",

			zap.String(
				"phone",
				phone,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return
	}

	logger.Log.Info(
		"otp deleted successfully",

		zap.String(
			"phone",
			phone,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)
}
