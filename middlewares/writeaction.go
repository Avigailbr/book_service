package middlewares

import (
	"fmt"
	"github.com/avigailbr/book_service/connectors"
	"github.com/avigailbr/book_service/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
	"net/http"
)

var counter float64

func WriteToRedis(userId, method, route string) error {
	//score := float64(time.Now().Nanosecond()) - don't work
	score:= counter
	counter+=1
	value := "method: " + method + ", route: " + route

	rdb := connectors.BookConnectors.RedisClients
	_, err := rdb.ZAdd(userId, redis.Z{score, value}).Result()
	if err != nil {
		return err
	}
	return nil
}

func WriteActionToRedis() gin.HandlerFunc {

	// Flush redis - only happens once
	rdb := connectors.BookConnectors.RedisClients
	rdb.FlushAll()

	return func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.NewStringError("username is required").Error())
			return
		}
		c.Next()

		if err := WriteToRedis(username, c.Request.Method, c.Request.RequestURI); err != nil {
			fmt.Printf("Write action to Redis failed. Method: %v, RequestURI: %v\n", c.Request.Method, c.Request.RequestURI)
		}

	}

}

