package service

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) *gin.Engine {
	router.GET("/ping", Ping)
	router.GET("/book/:id", BookInfo)
	router.PUT("/book", AddBook)
	router.POST("/book", UpdateBookTitle)
	router.DELETE("/book/:id", DeleteBook)
	router.GET("/search", searchBooks)
	router.GET("/store", storeInfo)
	router.GET("/activity", activityInfo)

	return router
}



