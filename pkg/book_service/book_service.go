package book_Service

import (
	"encoding/json"
	action "github.com/daniyakubov/book_service_n/pkg/action"
	"github.com/daniyakubov/book_service_n/pkg/book_service/models"
	"github.com/daniyakubov/book_service_n/pkg/cache"
	"github.com/daniyakubov/book_service_n/pkg/elastic_service"
	errors "github.com/fiverr/go_errors"
	"strings"
)

type BookService struct {
	booksCache     cache.Cache
	elasticHandler *elastic_service.ElasticHandler
}

func NewBookService(booksCache cache.Cache, elasticHandler *elastic_service.ElasticHandler) BookService {
	return BookService{
		booksCache:     booksCache,
		elasticHandler: elasticHandler,
	}
}

func (b *BookService) PutBook(req *models.Request) (*models.PutResponse, error) {

	postBody, err := json.Marshal(req.Data)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	id, err := b.elasticHandler.Put(postBody)

	if err != nil {
		return nil, err
	}

	err = b.booksCache.Push(req.Data.Username, "method:Put,"+"route:"+req.Route)
	if err != nil {
		return nil, err
	}
	res := models.PutResponse{id}
	return &res, nil

}

func (b *BookService) PostBook(req *models.Request) error {

	err := b.elasticHandler.Post(req.Data.Title, req.Data.Id)
	if err != nil {
		return err
	}
	err = b.booksCache.Push(req.Data.Username, "method:Post,"+"route:"+req.Route)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookService) GetBook(req *models.Request) (*models.Source, error) {

	src, err := b.elasticHandler.Get(req.Data.Id)

	src.Id = req.Data.Id

	err = b.booksCache.Push(req.Data.Username, "method:Get,"+"route:"+req.Route)
	if err != nil {
		return nil, err
	}
	return src, nil

}

func (b *BookService) DeleteBook(req *models.Request) error {

	err := b.elasticHandler.Delete(req.Data.Id)
	if err != nil {
		return err
	}

	err = b.booksCache.Push(req.Data.Username, "method:Delete,"+"route:"+req.Route)
	if err != nil {
		return err
	}
	return nil

}

func (b *BookService) Search(req *models.Request) ([]models.Source, error) {

	res, err := b.elasticHandler.Search(req.Data.Title, req.Data.Author, req.Data.PriceStart, req.Data.PriceEnd)
	if err != nil {
		return nil, err
	}

	err = b.booksCache.Push(req.Data.Username, "method:Get,"+"route:"+req.Route)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (b *BookService) Store(req *models.Request) (*models.StoreResponse, error) {

	count, distinctAuth, err := b.elasticHandler.Store()
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	resp := models.StoreResponse{count, distinctAuth}
	return &resp, nil
}

func (b *BookService) Activity(username string) ([]action.Action, error) {

	actions, err := b.booksCache.Get(username)
	if err != nil {
		return nil, err
	}
	var res []action.Action = make([]action.Action, int(len(actions)))

	for i := 0; i < len(actions); i++ {
		s := strings.Split(actions[i], ",")
		method := strings.Split(s[0], ":")[1]
		route := strings.Split(s[1], ":")[1]
		res[i].Method = method
		res[i].Route = route
	}

	return res, nil

}
