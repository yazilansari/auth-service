package redis

import (
	"auth-service/internal/logger"
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Ctx = context.Background()
var Client *redis.Client

func ConnectRedis() {
	logger.Log.Info(
		"connecting to otp redis database",
	)

	// =========================
	// REDIS DB
	// =========================

	db, err := strconv.Atoi(
		os.Getenv("REDIS_OTP_DB"),
	)

	if err != nil {

		logger.Log.Warn(
			"invalid redis db value, defaulting to 0",

			zap.String(
				"redis_db",
				os.Getenv("REDIS_OTP_DB"),
			),

			zap.Error(err),
		)

		db = 0
	}

	redisAddr :=
		os.Getenv("REDIS_HOST") +
			":" +
			os.Getenv("REDIS_PORT")

	logger.Log.Info(
		"creating redis client",

		zap.String(
			"redis_address",
			redisAddr,
		),

		zap.Int(
			"redis_db",
			db,
		),
	)

	start := time.Now()

	Client = redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           db,
		PoolSize:     20,
		MinIdleConns: 5,

		MaxRetries: 3,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		PoolTimeout: 4 * time.Second,
	})

	// =========================
	// TEST CONNECTION
	// =========================

	err = Client.Ping(Ctx).Err()

	duration := time.Since(start)

	if err != nil {

		logger.Log.Fatal(
			"failed to connect redis",

			zap.String(
				"redis_address",
				redisAddr,
			),

			zap.Int(
				"redis_db",
				db,
			),

			zap.Duration(
				"connection_duration",
				duration,
			),

			zap.Error(err),
		)

		return
	}

	logger.Log.Info(
		"redis connected successfully",

		zap.String(
			"redis_address",
			redisAddr,
		),

		zap.Int(
			"redis_db",
			db,
		),

		zap.Int(
			"pool_size",
			20,
		),

		zap.Int(
			"min_idle_connections",
			5,
		),

		zap.Duration(
			"connection_duration",
			duration,
		),
	)
}
