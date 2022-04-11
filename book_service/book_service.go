package book_Service

import (
	"context"
	"encoding/json"
	"github.com/daniyakubov/book_service_n/cache"
	"github.com/daniyakubov/book_service_n/db_service"
	"github.com/daniyakubov/book_service_n/models"
	errors "github.com/fiverr/go_errors"
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

func (b *BookService) AddBook(ctx *context.Context, book *models.Book, username string, route string) (string, error) {
	body, err := json.Marshal(book)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	id, err := b.dbHandler.Add(ctx, body)
	if err != nil {
		return "", err
	}

	err = b.booksCache.Push(username, "Put", route)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (b *BookService) UpdateBook(ctx *context.Context, book *models.Book, username string, route string) error {
	err := b.dbHandler.Update(ctx, book.Title, book.Id)
	if err != nil {
		return err
	}

	err = b.booksCache.Push(username, "Post", route)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookService) GetBook(ctx *context.Context, book *models.Book, username string, route string) (*models.Book, error) {
	src, err := b.dbHandler.Get(ctx, book.Id)
	src.Id = book.Id

	err = b.booksCache.Push(username, "Get", route)
	if err != nil {
		return nil, err
	}
	return src, nil
}

func (b *BookService) DeleteBook(ctx *context.Context, book *models.Book, username string, route string) error {
	err := b.dbHandler.Delete(ctx, book.Id)
	if err != nil {
		return err
	}

	err = b.booksCache.Push(username, "Delete", route)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookService) Search(book *models.Book, username string, route string, priceRange string) ([]models.Book, error) {
	res, err := b.dbHandler.Search(book.Title, book.Author, priceRange)
	if err != nil {
		return nil, err
	}

	err = b.booksCache.Push(username, "Get", route)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BookService) StoreInfo(username string, route string) (int64, int, error) {
	count, distinctAuth, err := b.dbHandler.StoreInfo()
	if err != nil {
		return 0, 0, errors.Wrap(err, err.Error())
	}

	err = b.booksCache.Push(username, "Get", route)
	if err != nil {
		return 0, 0, err
	}
	return count, distinctAuth, nil
}

func (b *BookService) Activity(username string) ([]models.Action, error) {
	actions, err := b.booksCache.Get(username)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
