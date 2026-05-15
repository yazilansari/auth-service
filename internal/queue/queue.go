package queue

import (
	"os"

	"github.com/hibiken/asynq"
)

var Client *asynq.Client

func InitQueue() {

	Client = asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: os.Getenv("REDIS_HOST") +
				":" +
				os.Getenv("REDIS_PORT"),

			Password: os.Getenv("REDIS_PASSWORD"),

			DB: 1,
		},
	)
}
