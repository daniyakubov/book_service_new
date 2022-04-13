package book_Service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/daniyakubov/book_service_n/cache"
	"github.com/daniyakubov/book_service_n/db_service"
	"github.com/daniyakubov/book_service_n/models"
	errors "github.com/fiverr/go_errors"
)

type BookService struct {
	booksCache cache.Cache
	dbHandler  db_service.DBHandler
}

func NewBookService(booksCache cache.Cache, dbHandler db_service.DBHandler) BookService {
	return BookService{
		booksCache: booksCache,
		dbHandler:  dbHandler,
	}
}

func (b *BookService) AddBook(ctx context.Context, book *models.Book, username string, routeName string) (string, error) {
	body, err := json.Marshal(book)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("couldn't unmarshal result of book with id: %s, in AddBook function ", book.Id))
	}

	id, err := b.dbHandler.AddBook(ctx, body)
	if err != nil {
		return "", err
	}

	err = b.booksCache.AddAction(username, "Put", routeName)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (b *BookService) UpdateBook(ctx context.Context, title string, id string, username string, routeName string) error {
	err := b.dbHandler.UpdateBook(ctx, title, id)
	if err != nil {
		return err
	}

	err = b.booksCache.AddAction(username, "Post", routeName)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookService) GetBook(ctx context.Context, id string, username string, routeName string) (*models.Book, error) {
	src, err := b.dbHandler.GetBook(ctx, id)
	src.Id = id

	err = b.booksCache.AddAction(username, "Get", routeName)
	if err != nil {
		return nil, err
	}
	return src, nil
}

func (b *BookService) DeleteBook(ctx context.Context, id string, username string, routeName string) error {
	err := b.dbHandler.DeleteBook(ctx, id)
	if err != nil {
		return err
	}

	err = b.booksCache.AddAction(username, "Delete", routeName)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookService) Search(title string, author string, username string, routeName string, priceRange string) ([]models.Book, error) {
	res, err := b.dbHandler.Search(title, author, priceRange)
	if err != nil {
		return nil, err
	}

	err = b.booksCache.AddAction(username, "Get", routeName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BookService) StoreInfo(username string, routeName string) (info map[string]interface{}, err error) {
	info, err = b.dbHandler.StoreInfo()
	if err != nil {
		return nil, err
	}

	err = b.booksCache.AddAction(username, "Get", routeName)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (b *BookService) Activity(username string) ([]models.Action, error) {
	actions, err := b.booksCache.GetLastActions(username)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
