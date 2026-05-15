package main

import (
	"auth-service/internal/queue"
	smsService "auth-service/internal/sms"

	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

func main() {

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: os.Getenv("REDIS_HOST") +
				":" +
				os.Getenv("REDIS_PORT"),

			Password: os.Getenv("REDIS_PASSWORD"),

			DB: 1,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{

				// queue name : priority weight

				"critical": 10,

				"default": 5,

				"low": 1,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(
		queue.TaskSendOTP,
		handleSendOTP,
	)

	if err := srv.Run(mux); err != nil {

		log.Fatal(err)
	}
}

func handleSendOTP(
	ctx context.Context,
	task *asynq.Task,
) error {

	log.Println("OTP Worker Triggered")

	var payload queue.OTPPayload

	err := json.Unmarshal(
		task.Payload(),
		&payload,
	)

	if err != nil {
		return err
	}

	log.Println(
		"Sending OTP To:",
		payload.Phone,
	)

	return smsService.SendOTP(
		payload.Phone,
		payload.OTP,
	)
}
