package datastore

import "github.com/avigailbr/book_service/models"

type IBookStorer interface {
	GetBook(id string) (*models.Book, error)
	InsertBook(book *models.Book) (string, error)
	UpdateBook(id string, fields map[string]interface{}) error
	DeleteBook(id string) error
	Search(fields map[string]string) ([]models.Book, error)
	Info() (map[string]interface{}, error)
}

type IActivityCacher interface {
	AddAction(userId, method, route string) error
	GetLastActions(username string) ([]string, error)
}
