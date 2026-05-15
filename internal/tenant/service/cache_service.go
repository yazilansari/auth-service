package service

import (
	"auth-service/internal/redis"
	"encoding/json"
	"time"
)

func SetTenantCache(
	key string,
	value interface{},
) error {

	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return redis.Client.Set(
		redis.Ctx,
		key,
		data,
		time.Hour,
	).Err()
}

func GetTenantCache(
	key string,
	dest interface{},
) error {

	data, err := redis.Client.Get(
		redis.Ctx,
		key,
	).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}
