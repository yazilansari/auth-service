package service

import (
	"auth-service/internal/queue"
	"auth-service/internal/redis"

	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())

	return fmt.Sprintf("%04d", rand.Intn(9000)+1000)
}

func SaveOTP(phone string, otp string) error {

	key := "otp:" + phone

	err := redis.Client.Set(
		redis.Ctx,
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

func VerifyOTP(phone string, otp string) bool {

	fmt.Println("phone:"+phone, "otp:"+otp)

	key := "otp:" + phone

	fmt.Println("Redis Key:" + key)

	storedOTP, err := redis.Client.Get(
		redis.Ctx,
		key,
	).Result()

	fmt.Println("stored OTP:" + storedOTP)

	if err != nil {
		return false
	}

	return storedOTP == otp
}

func DeleteOTP(phone string) {
	key := "otp:" + phone

	redis.Client.Del(redis.Ctx, key)
}
