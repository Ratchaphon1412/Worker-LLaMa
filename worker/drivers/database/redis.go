package database

import (
	"runtime"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/go-redis/redis/v8"
)

type RedisInstance struct {
	// Redis client
	Rd *redis.Client
}

var Redis RedisInstance

// ConnectRedis connects to Redis
func ConnectRedis(cfg *configs.Config) {
	// Initialize custom config
	store := redis.NewClient(&redis.Options{
		Addr:      cfg.REDIS_ADDR,
		Username:  cfg.REDIS_USERNAME,
		Password:  cfg.REDIS_PASSWORD,
		DB:        cfg.REDIS_DATABASE,
		PoolFIFO:  cfg.REDIS_POOLFIFO,
		TLSConfig: nil,
		PoolSize:  cfg.REDIS_POOL_SIZE * runtime.GOMAXPROCS(0),
	})

	// Set the Redis Client
	Redis.Rd = store

}
