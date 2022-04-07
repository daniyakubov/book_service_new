package http_service

import (
	"github.com/daniyakubov/book_service_n/pkg/book_service"
	"github.com/daniyakubov/book_service_n/pkg/book_service/models"
	"github.com/daniyakubov/book_service_n/pkg/consts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type HttpHandler struct {
	client      http.Client
	bookService book_Service.BookService
}

func NewHttpHandler(client http.Client, bookService book_Service.BookService) HttpHandler {
	return HttpHandler{
		client:      client,
		bookService: bookService,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (h *HttpHandler) PutBook(c *gin.Context) {
	var hit models.Hit

	err := c.Bind(&hit)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	r := models.NewRequest(&hit, c.FullPath())
	s, err := h.bookService.PutBook(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) PostBook(c *gin.Context) {

	var hit models.Hit

	err := c.Bind(&hit)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	r := models.NewRequest(&hit, c.FullPath())
	err = h.bookService.PostBook(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (h *HttpHandler) GetBook(c *gin.Context) {
	var hit models.Hit
	hit.Id = c.Query(consts.Id)
	hit.Username = c.Query(consts.UserName)

	r := models.NewRequest(&hit, c.FullPath())

	s, err := h.bookService.GetBook(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, s.Source)

}

func (h *HttpHandler) DeleteBook(c *gin.Context) {
	var hit models.Hit
	hit.Id = c.Query(consts.Id)
	hit.Username = c.Query(consts.UserName)

	r := models.NewRequest(&hit, c.FullPath())
	err := h.bookService.DeleteBook(&r)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (h *HttpHandler) Search(c *gin.Context) {
	var hit models.Hit
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
	r := models.NewRequest(&hit, c.FullPath())

	s, err := h.bookService.Search(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, s)

}

func (h *HttpHandler) Store(c *gin.Context) {

	var hit models.Hit
	hit.Username = c.Query(consts.UserName)

	r := models.NewRequest(&hit, c.FullPath())
	s, err := h.bookService.Store(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) Activity(c *gin.Context) {
	user := c.Query(consts.UserName)

	s, err := h.bookService.Activity(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, s)

}
