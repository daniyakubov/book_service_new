package http_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/app_service"
	"github.com/daniyakubov/book_service_n/book_service"
	"github.com/daniyakubov/book_service_n/consts"
	"github.com/daniyakubov/book_service_n/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ app_service.AppHandler = &HttpHandler{}

type HttpHandler struct {
	bookService book_Service.BookService
}

func NewHttpHandler(bookService book_Service.BookService) HttpHandler {
	return HttpHandler{
		bookService: bookService,
	}
}

func (h *HttpHandler) AddBook(c *gin.Context) {
	ctx := context.Background()
	var book models.Book
	if err := c.Bind(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.bookService.AddBook(&ctx, &book, c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *HttpHandler) UpdateBook(c *gin.Context) {
	ctx := context.Background()
	var book models.Book
	err := c.Bind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.bookService.UpdateBook(&ctx, book.Title, book.Id, c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) GetBook(c *gin.Context) {
	ctx := context.Background()

	s, err := h.bookService.GetBook(&ctx, c.Query(consts.Id), c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) DeleteBook(c *gin.Context) {
	ctx := context.Background()
	err := h.bookService.DeleteBook(&ctx, c.Query(consts.Id), c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	return
}

func (h *HttpHandler) Search(c *gin.Context) {
	s, err := h.bookService.Search(c.Query(consts.Title), c.Query(consts.Author), c.Query(consts.UserName), c.FullPath(), c.Query(consts.PriceRange))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) StoreInfo(c *gin.Context) {
	info, err := h.bookService.StoreInfo(c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books_num": info["books_num"], "distinct_authors_num": info["distinct_authors_num"]})
}

func (h *HttpHandler) Activity(c *gin.Context) {
	s, err := h.bookService.Activity(c.Query(consts.UserName))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}
