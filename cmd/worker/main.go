package main

import (
	"auth-service/internal/logger"
	"auth-service/internal/queue"
	smsService "auth-service/internal/sms"

	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func main() {

	logger.Log.Info(
		"starting auth worker",
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
		"connecting to queue redis",

		zap.String(
			"redis_address",
			redisAddr,
		),

		zap.Int(
			"redis_db",
			db,
		),
	)

	// =========================
	// ASYNQ SERVER
	// =========================

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,

			Password: os.Getenv("REDIS_PASSWORD"),

			DB: db,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{

				// queue name : priority weight

				"critical": 10,

				"default": 5,

				"low": 1,
			},

			Logger: NewAsynqLogger(),
		},
	)

	logger.Log.Info(
		"asynq worker initialized",

		zap.Int(
			"concurrency",
			10,
		),
	)

	// =========================
	// TASK ROUTES
	// =========================

	mux := asynq.NewServeMux()

	mux.HandleFunc(
		queue.TaskSendOTP,
		handleSendOTP,
	)

	logger.Log.Info(
		"queue handlers registered",
	)

	// =========================
	// START WORKER
	// =========================

	logger.Log.Info(
		"starting queue worker server",
	)

	if err := srv.Run(mux); err != nil {

		logger.Log.Fatal(
			"failed to run queue worker",

			zap.Error(err),
		)
	}
}

// =========================
// OTP WORKER
// =========================

func handleSendOTP(
	ctx context.Context,
	task *asynq.Task,
) error {

	start := time.Now()

	logger.Log.Info(
		"otp worker triggered",
	)

	var payload queue.OTPPayload

	err := json.Unmarshal(
		task.Payload(),
		&payload,
	)

	if err != nil {
		logger.Log.Error(
			"failed to unmarshal otp payload",

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"sending otp",

		zap.String(
			"phone",
			payload.Phone,
		),
	)

	err = smsService.SendOTP(
		payload.Phone,
		payload.OTP,
	)

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to send otp",

			zap.String(
				"phone",
				payload.Phone,
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
		"otp sent successfully",

		zap.String(
			"phone",
			payload.Phone,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return nil
}

// =========================
// ASYNQ LOGGER
// =========================

type AsynqLogger struct{}

func NewAsynqLogger() *AsynqLogger {
	return &AsynqLogger{}
}

func (l *AsynqLogger) Debug(
	args ...interface{},
) {
	logger.Log.Debug(
		"asynq debug",
		zap.Any("data", args),
	)
}

func (l *AsynqLogger) Info(
	args ...interface{},
) {
	logger.Log.Info(
		"asynq info",
		zap.Any("data", args),
	)
}

func (l *AsynqLogger) Warn(
	args ...interface{},
) {
	logger.Log.Warn(
		"asynq warn",
		zap.Any("data", args),
	)
}

func (l *AsynqLogger) Error(
	args ...interface{},
) {
	logger.Log.Error(
		"asynq error",
		zap.Any("data", args),
	)
}

func (l *AsynqLogger) Fatal(
	args ...interface{},
) {
	logger.Log.Fatal(
		"asynq fatal",
		zap.Any("data", args),
	)
}
