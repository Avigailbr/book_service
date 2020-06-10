package main

import (
	"fmt"
	"github.com/avigailbr/book_service/connectors"
	"github.com/avigailbr/book_service/controllers"
	"github.com/avigailbr/book_service/middlewares"
	"github.com/avigailbr/book_service/service"
	"github.com/gin-gonic/gin"
)
func main() {
	err := connectors.NewConnectors()
	if err != nil {
		fmt.Printf("Initialization failed, %v" ,err.Error())
		return
	}
	controllers.BookController = controllers.NewElasticBookController()

	router := gin.New()

	// middelwares
	router.Use(gin.Logger(), gin.Recovery(), middlewares.WriteActionToRedis())

	service.Routes(router)

	router.Run(":8888")
	//panic(connectors.BookConnectors.Shutdown) - TODO (ASK)- do i need to run `router.Run` inside a goroutine

}

