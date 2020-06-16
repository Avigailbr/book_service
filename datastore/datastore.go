package datastore

import "github.com/avigailbr/book_service/models"

type BookStorer interface {
	Get(id string) (*models.Book, error)
	Insert(book *models.Book) (string, error)
	Update(id string, fields map[string]interface{}) error
	Delete(id string) error
	Search(fields map[string]string) ([]models.Book, error)
	Info() (map[string]interface{}, error)
}

type ActivityCacher interface {
	AddAction(userId, method, route string) error
	GetLastActions(username string) ([]string, error)
}
