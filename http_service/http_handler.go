package http_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/book_service"
	"github.com/daniyakubov/book_service_n/consts/elastic_fields"
	"github.com/daniyakubov/book_service_n/consts/query_fields"
	"github.com/daniyakubov/book_service_n/models"
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

	c.JSON(http.StatusBadRequest, "book updated")
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

	c.JSON(http.StatusBadRequest, "book deleted")
}

func (h *HttpHandler) Search(c *gin.Context) {
	searchParams := make(map[string]string)
	searchParams["title"] = c.Query(elastic_fields.Title)
	searchParams["author"] = c.Query(elastic_fields.Author)
	searchParams["price_range"] = c.Query(query_fields.PriceRange)
	s, err := h.bookService.Search(context.Background(), searchParams, c.Query(query_fields.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) StoreInfo(c *gin.Context) {
	info, err := h.bookService.StoreInfo(context.Background(), c.Query(query_fields.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"store_info": info,
	})
}

func (h *HttpHandler) Activity(c *gin.Context) {
	s, err := h.bookService.Activity(c.Query(query_fields.UserName))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}
