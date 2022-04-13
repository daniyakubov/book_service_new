package http_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/book_service"
	"github.com/daniyakubov/book_service_n/elastic_fields"
	"github.com/daniyakubov/book_service_n/models"
	"github.com/daniyakubov/book_service_n/query_fields"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpHandler struct {
	bookService book_Service.BookService
}

func NewHttpHandler(bookService book_Service.BookService) HttpHandler {
	return HttpHandler{
		bookService: bookService,
	}
}

func (h *HttpHandler) AddBook(c *gin.Context) {
	var book models.Book
	if err := c.Bind(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.bookService.InsertBook(context.Background(), &book, c.Query(query_fields.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *HttpHandler) UpdateBook(c *gin.Context) {
	var book models.Book
	err := c.Bind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.bookService.UpdateBook(context.Background(), book.Title, book.Id, c.Query(query_fields.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) GetBook(c *gin.Context) {
	s, err := h.bookService.GetBook(context.Background(), c.Query(elastic_fields.Id), c.Query(query_fields.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) DeleteBook(c *gin.Context) {
	err := h.bookService.DeleteBook(context.Background(), c.Query(elastic_fields.Id), c.Query(query_fields.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	return
}

func (h *HttpHandler) Search(c *gin.Context) {
	s, err := h.bookService.Search(c.Query(elastic_fields.Title), c.Query(elastic_fields.Author), c.Query(query_fields.UserName), c.FullPath(), c.Query(query_fields.PriceRange))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) StoreInfo(c *gin.Context) {
	info, err := h.bookService.StoreInfo(c.Query(query_fields.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books_num": info["books_num"], "distinct_authors_num": info["distinct_authors_num"]})
}

func (h *HttpHandler) Activity(c *gin.Context) {
	s, err := h.bookService.Activity(c.Query(query_fields.UserName))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}
