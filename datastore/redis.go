package datastore

import (
	"fmt"
	"github.com/avigailbr/book_service/connectors"
)

type RedisCache struct {
	Client *connectors.RedisClient
}

func NewRedisCache(addr string) (IActivityCacher, error) {
	redisClient, err := connectors.NewRedisClient(addr)
	if err != nil {
		return nil, err
	}
	redisCache := &RedisCache{
		Client: redisClient,
	}
	return redisCache, nil
}

func (r *RedisCache) AddAction(userId, method, route string) error {
	value := "method: " + method + ", route: " + route

	rdb := r.Client
	_, err := rdb.LPush(userId, value).Result()
	if err != nil {
		return err
	}
	return nil

}
func (r *RedisCache) GetLastActions(username string) ([]string, error){
	rdb := r.Client

	items, err := rdb.LRange(username, 0, 2).Result()
	if err !=nil{
		return nil, err
	}

	var actions []string
	for _, val := range items {
		actions = append(actions, fmt.Sprint(val))
	}

	return actions, nil

}
