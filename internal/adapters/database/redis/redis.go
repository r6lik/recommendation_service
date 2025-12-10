package redis

import (
	"time"

	"github.com/r6lik/recommendation_service/internal/adapters/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,

		PoolSize:     10,
		MinIdleConns: 5,

		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	})

	return client
}
