package queue

import (
	"auth-service/internal/logger"
	"os"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

var Client *asynq.Client

func InitQueue() {

	logger.Log.Info(
		"initializing asynq queue client",
	)

	// =========================
	// REDIS QUEUE DB
	// =========================

	db, err := strconv.Atoi(
		os.Getenv("REDIS_QUEUE_DB"),
	)

	if err != nil {

		logger.Log.Warn(
			"invalid redis db value, defaulting to 1",

			zap.String(
				"redis_db",
				os.Getenv("REDIS_QUEUE_DB"),
			),

			zap.Error(err),
		)

		db = 1
	}

	redisAddr :=
		os.Getenv("REDIS_HOST") +
			":" +
			os.Getenv("REDIS_PORT")

	logger.Log.Info(
		"connecting queue client to redis",

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

	// =========================
	// CREATE CLIENT
	// =========================

	Client = asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: redisAddr,

			Password: os.Getenv("REDIS_PASSWORD"),

			DB: db,
		},
	)

	duration := time.Since(start)

	logger.Log.Info(
		"asynq queue client initialized successfully",

		zap.String(
			"redis_address",
			redisAddr,
		),

		zap.Int(
			"redis_db",
			db,
		),

		zap.Duration(
			"initialization_duration",
			duration,
		),
	)
}
