package config

import (
	"go-servie/utils"

	"github.com/redis/go-redis/v9"
)

func RedisConfig() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     utils.REDIS_HOST,
		Password: "",
		DB:       0,
	})
}
