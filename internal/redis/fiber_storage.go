package redis

import (
	"auth-service/internal/logger"
	"os"
	"strconv"
	"time"

	redisStorage "github.com/gofiber/storage/redis/v3"
	"go.uber.org/zap"
)

var FiberStorage *redisStorage.Storage

func InitFiberStorage() {

	logger.Log.Info(
		"initializing rate limit redis storage",
	)

	// =========================
	// REDIS RATE LIMIT DB
	// =========================

	db, err := strconv.Atoi(
		os.Getenv("REDIS_RATE_LIMIT_DB"),
	)

	if err != nil {

		logger.Log.Warn(
			"invalid redis db value, defaulting to 2",

			zap.String(
				"redis_db",
				os.Getenv("REDIS_RATE_LIMIT_DB"),
			),

			zap.Error(err),
		)

		db = 2
	}

	port, err := strconv.Atoi(
		os.Getenv("REDIS_PORT"),
	)

	if err != nil {

		logger.Log.Fatal(
			"invalid redis port",

			zap.String(
				"redis_port",
				os.Getenv("REDIS_PORT"),
			),

			zap.Error(err),
		)

		return
	}

	redisHost := os.Getenv("REDIS_HOST")

	logger.Log.Info(
		"creating redis storage for rate limiting",

		zap.String(
			"redis_host",
			redisHost,
		),

		zap.Int(
			"redis_port",
			port,
		),

		zap.Int(
			"redis_db",
			db,
		),
	)

	start := time.Now()

	// =========================
	// CREATE STORAGE
	// =========================

	FiberStorage = redisStorage.New(
		redisStorage.Config{
			Host: os.Getenv("REDIS_HOST"),

			Port: port,

			Password: os.Getenv("REDIS_PASSWORD"),

			Database: db,

			Reset: false,
		},
	)

	duration := time.Since(start)

	logger.Log.Info(
		"rate limit redis storage initialized successfully",

		zap.String(
			"redis_host",
			redisHost,
		),

		zap.Int(
			"redis_port",
			port,
		),

		zap.Int(
			"redis_db",
			db,
		),

		zap.Bool(
			"reset",
			false,
		),

		zap.Duration(
			"initialization_duration",
			duration,
		),
	)
}
