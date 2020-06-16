package main

import (
	"github.com/avigailbr/book_service/middlewares"
	"github.com/avigailbr/book_service/service"
	"github.com/gin-gonic/gin"
)
func main() {
	router := gin.New()

	// middelwares
	router.Use(gin.Logger(), gin.Recovery(),middlewares.Database(), middlewares.WriteActionToCache())

	service.Routes(router)

	router.Run(":8888")
}

