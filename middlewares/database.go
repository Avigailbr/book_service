package middlewares

import (
	"github.com/avigailbr/book_service/config"
	"github.com/avigailbr/book_service/datastore"
	"github.com/gin-gonic/gin"
)

func Database() gin.HandlerFunc {
	db, err := datastore.NewElasticBookStore(config.ElasticUrl, config.IndexName, config.TypeName)
	if err != nil {
		panic(err)
	}
	cache, err := datastore.NewRedisCache(config.RedisAdd)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("Cache", cache)
		c.Next()
	}
}
