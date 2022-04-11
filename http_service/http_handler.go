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

func (h *HttpHandler) AddBook(c *gin.Context) {
	ctx := context.Background()
	var book models.Book
	err := c.Bind(&book)
	if err != nil {
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

	err = h.bookService.UpdateBook(&ctx, &book, c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) GetBook(c *gin.Context) {
	ctx := context.Background()
	var book models.Book
	book.Id = c.Query(consts.Id)

	s, err := h.bookService.GetBook(&ctx, &book, c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) DeleteBook(c *gin.Context) {
	ctx := context.Background()
	var book models.Book
	book.Id = c.Query(consts.Id)
	err := h.bookService.DeleteBook(&ctx, &book, c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) Search(c *gin.Context) {
	var book models.Book
	book.Title = c.Query(consts.Title)
	book.Author = c.Query(consts.Author)
	priceRange := c.Query(consts.PriceRange)

	s, err := h.bookService.Search(&book, c.Query(consts.UserName), c.FullPath(), priceRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) StoreInfo(c *gin.Context) {
	count, distinctAuthors, err := h.bookService.StoreInfo(c.Query(consts.UserName), c.FullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books_num": count, "distinct_authors_num": distinctAuthors})
}

func (h *HttpHandler) Activity(c *gin.Context) {
	s, err := h.bookService.Activity(c.Query(consts.UserName))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) ApplyRoutes(router *gin.Engine) {
	book := router.Group("/book")
	{
		book.GET("", h.GetBook)
		book.DELETE("", h.DeleteBook)
		book.PUT("", h.AddBook)
		book.POST("", h.UpdateBook)
	}
	router.GET("/search", h.Search)
	router.GET("/store", h.StoreInfo)
	router.GET("/activity", h.Activity)
}
