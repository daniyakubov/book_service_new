package book_Service

import (
	"context"
	"encoding/json"
	"github.com/daniyakubov/book_service_n/pkg/cache"
	"github.com/daniyakubov/book_service_n/pkg/db_service"
	models2 "github.com/daniyakubov/book_service_n/pkg/models"
	errors "github.com/fiverr/go_errors"
	"strings"
)

type BookService struct {
	booksCache cache.Cache
	dbHandler  db_service.DbHandler
}

func NewBookService(booksCache cache.Cache, dbHandler db_service.DbHandler) BookService {
	return BookService{
		booksCache: booksCache,
		dbHandler:  dbHandler,
	}
}

func (b *BookService) AddBook(ctx *context.Context, req *models2.Request) (string, error) {
	postBody, err := json.Marshal(req.Data)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	id, err := b.dbHandler.Put(ctx, postBody)
	if err != nil {
		return "", err
	}

	err = b.booksCache.Push(req.Data.Username, "method:Put,"+"route:"+req.Route)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (b *BookService) UpdateBook(ctx *context.Context, req *models2.Request) error {
	err := b.dbHandler.Post(ctx, req.Data.Title, req.Data.Id)
	if err != nil {
		return err
	}

	err = b.booksCache.Push(req.Data.Username, "method:Post,"+"route:"+req.Route)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookService) GetBook(ctx *context.Context, req *models2.Request) (*models2.Book, error) {
	src, err := b.dbHandler.Get(ctx, req.Data.Id)
	src.Id = req.Data.Id

	err = b.booksCache.Push(req.Data.Username, "method:Get,"+"route:"+req.Route)
	if err != nil {
		return nil, err
	}
	return src, nil
}

func (b *BookService) DeleteBook(ctx *context.Context, req *models2.Request) error {
	err := b.dbHandler.Delete(ctx, req.Data.Id)
	if err != nil {
		return err
	}

	err = b.booksCache.Push(req.Data.Username, "method:Delete,"+"route:"+req.Route)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookService) Search(req *models2.Request) ([]models2.Book, error) {
	res, err := b.dbHandler.Search(req.Data.Title, req.Data.Author, req.Data.PriceStart, req.Data.PriceEnd)
	if err != nil {
		return nil, err
	}

	err = b.booksCache.Push(req.Data.Username, "method:Get,"+"route:"+req.Route)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BookService) StoreInfo(req *models2.Request) (*models2.StoreResponse, error) {
	count, distinctAuth, err := b.dbHandler.StoreInfo()
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	resp := models2.StoreResponse{count, distinctAuth}

	err = b.booksCache.Push(req.Data.Username, "method:Get,"+"route:"+req.Route)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (b *BookService) Activity(username string) ([]models2.Action, error) {
	actions, err := b.booksCache.Get(username)
	if err != nil {
		return nil, err
	}

	res := make([]models2.Action, int(len(actions)))
	for i := 0; i < len(actions); i++ {
		s := strings.Split(actions[i], ",")
		method := strings.Split(s[0], ":")[1]
		route := strings.Split(s[1], ":")[1]
		res[i].Method = method
		res[i].Route = route
	}
	return res, nil
}
