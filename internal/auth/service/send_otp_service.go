package service

import (
	"auth-service/internal/queue"
	redisClient "auth-service/internal/redis"

	"fmt"
	"math/rand"
	"time"
)

func SendOTP(
	phone string,
) error {

	// Generate OTP

	rand.Seed(time.Now().UnixNano())

	otp := fmt.Sprintf(
		"%04d",
		rand.Intn(10000),
	)

	// Store In Redis

	key := "otp:" + phone

	err := redisClient.Client.Set(
		redisClient.Ctx,
		key,
		otp,
		time.Minute*2,
	).Err()

	if err != nil {
		return err
	}

	// Push Queue Job

	err = queue.EnqueueOTP(
		phone,
		otp,
	)

	if err != nil {
		return err
	}

	return nil
}
