package connectors

import (
	"gopkg.in/redis.v5"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(addr string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisClient := &RedisClient{
		client,
	}

	if _, err := redisClient.Ping().Result(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
