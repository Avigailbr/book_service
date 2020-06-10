package connectors

import (
	"github.com/avigailbr/book_service/config"
	"gopkg.in/redis.v5"
)

var BookConnectors *Connectors

type Connectors struct {
	ElasticClient *ElasticClient
	RedisClients  *redis.Client
}

func NewConnectors() error {
	BookConnectors = &Connectors{
	}
	_, err := BookConnectors.NewElasticClient(config.ElasticUrl, config.IndexName, config.TypeName)
	if err != nil {
		return err
	}
	_, err = BookConnectors.NewRedisClient(config.RedisAdd, config.RedisDB)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connectors) Shutdown() {
	if c.RedisClients != nil {
		c.RedisClients.Close()
	}
}
