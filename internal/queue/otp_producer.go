package queue

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type OTPPayload struct {
	Phone string `json:"phone"`

	OTP string `json:"otp"`
}

func EnqueueOTP(
	phone string,
	otp string,
) error {

	payload, err := json.Marshal(
		OTPPayload{
			Phone: phone,
			OTP:   otp,
		},
	)

	if err != nil {
		return err
	}

	task := asynq.NewTask(
		TaskSendOTP,
		payload,
	)

	_, err = Client.Enqueue(
		task,

		// Retry failed jobs 10 times

		asynq.MaxRetry(10),

		// Optional Queue Name

		asynq.Queue("critical"),

		// Optional Timeout

		asynq.Timeout(30*time.Second),
	)

	return err
}
