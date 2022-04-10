package http_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/book_service"
	"github.com/daniyakubov/book_service_n/consts"
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

func (h *HttpHandler) PutBook(c *gin.Context) {
	ctx := context.Background()
	var book models.UserBook
	err := c.Bind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.bookService.AddBook(&ctx, &book, c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *HttpHandler) PostBook(c *gin.Context) {
	ctx := context.Background()
	var book models.UserBook
	err := c.Bind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.bookService.UpdateBook(&ctx, &book, c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) GetBook(c *gin.Context) {
	ctx := context.Background()
	var book models.UserBook
	book.Id = c.Query(consts.Id)
	book.Username = c.Query(consts.UserName)

	s, err := h.bookService.GetBook(&ctx, &book, c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) DeleteBook(c *gin.Context) {
	ctx := context.Background()
	var book models.UserBook
	book.Id = c.Query(consts.Id)
	book.Username = c.Query(consts.UserName)

	err := h.bookService.DeleteBook(&ctx, &book, c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) Search(c *gin.Context) {
	var book models.UserBook
	book.Title = c.Query(consts.Title)
	book.Author = c.Query(consts.Author)
	priceRange := c.Query(consts.PriceRange)
	book.Username = c.Query(consts.UserName)

	s, err := h.bookService.Search(&book, c.FullPath(), priceRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) StoreInfo(c *gin.Context) {
	var book models.UserBook
	book.Username = c.Query(consts.UserName)
	s, err := h.bookService.StoreInfo(&book, c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) Activity(c *gin.Context) {
	user := c.Query(consts.UserName)
	s, err := h.bookService.Activity(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}
