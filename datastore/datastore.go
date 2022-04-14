package datastore

import (
	"context"
	"github.com/daniyakubov/book_service_n/models"
)

type BookStorer interface {
	UpdateBook(ctx context.Context, title string, id string) error
	InsertBook(ctx context.Context, body []byte) (string, error)
	GetBook(ctx context.Context, id string) (*models.Book, error)
	DeleteBook(ctx context.Context, id string) error
	Search(ctx context.Context, fields map[string]string) ([]models.Book, error)
	StoreInfo(ctx context.Context) (map[string]interface{}, error)
}
