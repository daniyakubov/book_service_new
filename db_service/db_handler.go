package db_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/models"
)

type DbHandler interface {
	Update(ctx *context.Context, title string, id string) error
	Add(ctx *context.Context, body []byte) (string, error)
	Get(ctx *context.Context, id string) (*models.Book, error)
	Delete(ctx *context.Context, id string) error
	Search(title string, author string, priceRange string) ([]models.Book, error)
	StoreInfo() (int64, int, error)
}
