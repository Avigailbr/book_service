package service

import (
	"fmt"
	"github.com/avigailbr/book_service/datastore"
	"github.com/avigailbr/book_service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateBook struct {
	Id    string `json:"id" binding:"required`
	Title string `json:"title" binding:"required"`
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func BookInfo(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("DB").(datastore.IBookStorer)
	book, err := db.GetBook(id)
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
	db := c.MustGet("DB").(datastore.IBookStorer)
	id, err := db.InsertBook(&book)
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

	var book UpdateBook
	if err := c.ShouldBind(&book); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := c.MustGet("DB").(datastore.IBookStorer)
	err := db.UpdateBook(book.Id, map[string]interface{}{"title": book.Title})
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

	db := c.MustGet("DB").(datastore.IBookStorer)
	err := db.DeleteBook(id)
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
	db := c.MustGet("DB").(datastore.IBookStorer)
	books, err := db.Search(fields)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"hits": books,
	})

}

func storeInfo(c *gin.Context) {

	db := c.MustGet("DB").(datastore.IBookStorer)
	info, err := db.Info()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"store_info": info,
	})

}

func activityInfo(c *gin.Context) {
	username := c.Query("username")

	cache := c.MustGet("Cache").(datastore.IActivityCacher)

	actions, err := cache.GetLastActions(username)

	if err !=nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":     username,
		"last_actions": actions,
	})
}
