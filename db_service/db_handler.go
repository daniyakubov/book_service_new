package db_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/models"
)

type DBHandler interface {
	UpdateBook(ctx context.Context, title string, id string) error
	AddBook(ctx context.Context, body []byte) (string, error)
	GetBook(ctx context.Context, id string) (*models.Book, error)
	DeleteBook(ctx context.Context, id string) error
	Search(title string, author string, priceRange string) ([]models.Book, error)
	StoreInfo() (map[string]interface{}, error)
}
