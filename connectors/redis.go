package connectors

import (
	"gopkg.in/redis.v5"
)

func (c *Connectors) NewRedisClient(addr string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	c.RedisClients = client
	return client, nil
}
