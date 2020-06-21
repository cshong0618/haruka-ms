package redis

import "github.com/go-redis/redis/v8"

func InitRedis(redisUrl string) *redis.Client {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)
	if client == nil {
		panic("cannot create redis client")
	}

	return client
}
