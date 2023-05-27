package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	conf "github.com/mp3ismp3/gocrypto/config"
)

func NewRedisClient() *redis.Client {
	rConfig := conf.Config.Redis
	fmt.Println("hello")
	cache := redis.NewClient(&redis.Options{
		Addr:     rConfig.RedisHost,
		Password: "", // no password set
		DB:       0,  //use default DB
	})

	return cache
}
