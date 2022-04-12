package elastic_service

import "github.com/daniyakubov/book_service_n/models"

type StoreInformation interface {
	Info(key string) ([]models.Action, error)
}
