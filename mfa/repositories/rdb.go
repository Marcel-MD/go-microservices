package repositories

import (
	"context"
	"mfa/config"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	rdbOnce sync.Once
	rdb     *redis.Client
)

func GetRDB() *redis.Client {
	rdbOnce.Do(func() {
		log.Info().Msg("Initializing redis")
		cfg := config.GetConfig()

		opt, err := redis.ParseURL(cfg.RedisUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse redis connection string")
		}

		redisDb := redis.NewClient(opt)
		redisCtx := context.Background()

		status := redisDb.Ping(redisCtx)
		if status.Err() != nil {
			log.Fatal().Err(status.Err()).Msg("Failed to connect to redis")
		}

		rdb = redisDb
	})

	return rdb
}
