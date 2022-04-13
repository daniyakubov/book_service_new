package book_Service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/daniyakubov/book_service_n/datastore"
	"github.com/daniyakubov/book_service_n/models"
	errors "github.com/fiverr/go_errors"
)

type BookService struct {
	dbHandler datastore.BookStorer
}

func NewBookService(dbHandler datastore.BookStorer) BookService {
	return BookService{
		dbHandler: dbHandler,
	}
}

func (b *BookService) InsertBook(ctx context.Context, book *models.Book) (string, error) {
	body, err := json.Marshal(book)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("couldn't unmarshal result of book with id: %s, in AddBook function ", book.Id))
	}

	id, err := b.dbHandler.InsertBook(ctx, body)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *BookService) UpdateBook(ctx context.Context, title string, id string) error {
	err := b.dbHandler.UpdateBook(ctx, title, id)
	if err != nil {
		return err
	}

	return nil
}

func (b *BookService) GetBook(ctx context.Context, id string) (*models.Book, error) {
	src, err := b.dbHandler.GetBook(ctx, id)
	if err != nil {
		return nil, err
	}

	src.Id = id
	return src, nil
}

func (b *BookService) DeleteBook(ctx context.Context, id string) error {
	err := b.dbHandler.DeleteBook(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (b *BookService) Search(ctx context.Context, fields map[string]string) ([]models.Book, error) {

	res, err := b.dbHandler.Search(ctx, fields)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BookService) StoreInfo(ctx context.Context) (info map[string]interface{}, err error) {
	info, err = b.dbHandler.StoreInfo(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}
