package cache

import "github.com/daniyakubov/book_service_n/models"

type Cache interface {
	Push(key string, method string, route string) error
	Get(key string) ([]models.Action, error)
}
