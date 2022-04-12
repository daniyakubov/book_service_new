package cache

import "github.com/daniyakubov/book_service_n/models"

type Cache interface {
	AddAction(key string, method string, routeName string) error
	Get(key string) ([]models.Action, error)
}
