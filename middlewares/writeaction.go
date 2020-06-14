package middlewares

import (
	"errors"
	"fmt"
	"github.com/avigailbr/book_service/datastore"
	"github.com/gin-gonic/gin"
	"net/http"
)

var counter float64


func WriteActionToCache() gin.HandlerFunc {

	return func(c *gin.Context) {
		cache := c.MustGet("Cache").(datastore.IActivityCacher)
		username := c.Query("username")
		if username == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("username is required"))
			return
		}
		c.Next()

		if err := cache.AddAction(username, c.Request.Method, c.Request.RequestURI); err != nil {
			fmt.Printf("Write action to Redis failed. Method: %v, RequestURI: %v\n", c.Request.Method, c.Request.RequestURI)
		}

	}

}

