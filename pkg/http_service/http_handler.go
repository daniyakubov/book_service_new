package http_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/pkg/book_service"
	"github.com/daniyakubov/book_service_n/pkg/consts"
	models2 "github.com/daniyakubov/book_service_n/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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
	var hit models2.Hit
	err := c.Bind(&hit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := models2.NewRequest(&hit, c.FullPath())
	id, err := h.bookService.AddBook(&ctx, &r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *HttpHandler) PostBook(c *gin.Context) {
	ctx := context.Background()
	var hit models2.Hit
	err := c.Bind(&hit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := models2.NewRequest(&hit, c.FullPath())
	err = h.bookService.UpdateBook(&ctx, &r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) GetBook(c *gin.Context) {
	ctx := context.Background()
	var hit models2.Hit
	hit.Id = c.Query(consts.Id)
	hit.Username = c.Query(consts.UserName)

	r := models2.NewRequest(&hit, c.FullPath())
	s, err := h.bookService.GetBook(&ctx, &r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) DeleteBook(c *gin.Context) {
	ctx := context.Background()
	var hit models2.Hit
	hit.Id = c.Query(consts.Id)
	hit.Username = c.Query(consts.UserName)

	r := models2.NewRequest(&hit, c.FullPath())
	err := h.bookService.DeleteBook(&ctx, &r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *HttpHandler) Search(c *gin.Context) {
	var hit models2.Hit
	hit.Title = c.Query(consts.Title)
	hit.Author = c.Query(consts.Author)
	sRange := c.Query(consts.PriceRange)
	if sRange == "" {
		hit.PriceStart = 0
		hit.PriceEnd = 0
	} else {
		priceRange := strings.Split(sRange, "-")
		hit.PriceStart, _ = strconv.ParseFloat(priceRange[0], 32)
		hit.PriceEnd, _ = strconv.ParseFloat(priceRange[1], 32)
	}
	hit.Username = c.Query(consts.UserName)

	r := models2.NewRequest(&hit, c.FullPath())
	s, err := h.bookService.Search(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) StoreInfo(c *gin.Context) {
	var hit models2.Hit
	hit.Username = c.Query(consts.UserName)
	r := models2.NewRequest(&hit, c.FullPath())
	s, err := h.bookService.StoreInfo(&r)
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
