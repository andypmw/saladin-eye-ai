package cache

import (
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/redis/go-redis/v9"
)

// Singleton
var (
	redisClient redis.Cmdable
	once        sync.Once
)

func New() redis.Cmdable {
	once.Do(func() {
		redisAddr := os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			log.Fatal().Msg("REDIS_ADDR environment variable not set")
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr: redisAddr,
			DB:   0, // use default DB
		})
	})
	return redisClient
}
