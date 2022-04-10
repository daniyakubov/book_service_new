package db_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/pkg/models"
)

type DbHandler interface {
	Post(ctx *context.Context, title string, id string) error
	Put(ctx *context.Context, postBody []byte) (string, error)
	Get(ctx *context.Context, id string) (*models.Book, error)
	Delete(ctx *context.Context, id string) error
	Search(title string, author string, priceStart float64, priceEnd float64) ([]models.Book, error)
	StoreInfo() (int64, int, error)
}
