package service

import (
	"fmt"
	"github.com/avigailbr/book_service/connectors"
	"github.com/avigailbr/book_service/controllers"
	"github.com/avigailbr/book_service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func BookInfo(c *gin.Context) {
	id := c.Param("id")
	book, err := controllers.BookController.GetBook(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"book_details": book,
	})
}

func AddBook(c *gin.Context) {

	var book models.Book
	if err := c.ShouldBind(&book); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := controllers.BookController.InsertBook(&book)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Indexed book successfully! id: %s\n", id)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Book successfully indexed",
		"id":      id,
	})

}

func UpdateBookTitle(c *gin.Context) {

	var book models.UpdateBook
	if err := c.ShouldBind(&book); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := controllers.BookController.UpdateBook(book.Id, map[string]interface{}{"title": book.Title})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("updated book successfully, id: ", book.Id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully",
	})

}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	err := controllers.BookController.DeleteBook(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})

}

func searchBooks(c *gin.Context) {

	title := c.Query("title")
	author_name := c.Query("author_name")
	price_range := c.Query("price_range")

	fields := map[string]string{

		"title":       title,
		"author_name": author_name,
		"price_range": price_range,
	}
	books, err := controllers.BookController.Search(fields)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"hits": books,
	})

}

func storeInfo(c *gin.Context) {

	info, err := controllers.BookController.Info()
	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"store_info": info,
	})

}

func activityInfo(c *gin.Context) {
	username := c.Query("username")
	rdb := connectors.BookConnectors.RedisClients

	items, err := rdb.ZRevRangeWithScores(username, 0, 2).Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var actions []string
	for _, zItem := range items {
		actions = append(actions, fmt.Sprint(zItem.Member))
	}

	c.JSON(http.StatusOK, gin.H{
		"username":     username,
		"last_actions": actions,
	})
}
