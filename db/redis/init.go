package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	// rConfig := conf.Config.Redis
	// pathRead := strings.Join([]string{rConfig.RedisHost, ":", rConfig.RedisPort}, "")
	cache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  //use default DB
	})

	_, err := cache.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return cache
}
