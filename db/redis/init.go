package redis

import (
	"github.com/go-redis/redis/v8"
	conf "github.com/mp3ismp3/gocrypto/config"
)

var cache *redis.Client

func NewRedisClient() {
	rConfig := conf.Config.Redis
	cache = redis.NewClient(&redis.Options{
		Addr:     rConfig.RedisHost,
		Password: "", // no password set
		DB:       0,  //use default DB
	})
}

func GetCache() *redis.Client {
	return cache
}
