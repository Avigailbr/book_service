package service

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) *gin.Engine {
	router.GET("/ping", Ping)
	book := router.Group("/book")
	{
		book.GET("/:id", BookInfo)
		book.PUT("/", AddBook)
		book.POST("/", UpdateBookTitle)
		book.DELETE("/:id", DeleteBook)

	}
	router.GET("/search", searchBooks)
	router.GET("/store", storeInfo)
	router.GET("/activity", activityInfo)

	return router
}



