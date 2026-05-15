package redis

import (
	"os"
	"strconv"

	redisStorage "github.com/gofiber/storage/redis/v3"
)

var FiberStorage *redisStorage.Storage

func InitFiberStorage() {

	port, _ := strconv.Atoi(
		os.Getenv("REDIS_PORT"),
	)

	FiberStorage = redisStorage.New(
		redisStorage.Config{
			Host: os.Getenv("REDIS_HOST"),

			Port: port,

			Password: os.Getenv("REDIS_PASSWORD"),

			Database: 2,

			Reset: false,
		},
	)
}
