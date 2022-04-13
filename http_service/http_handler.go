package http_service

import (
	"context"
	"github.com/daniyakubov/book_service_n/book_service"
	"github.com/daniyakubov/book_service_n/cache"
	"github.com/daniyakubov/book_service_n/consts/elastic_fields"
	"github.com/daniyakubov/book_service_n/consts/query_fields"
	"github.com/daniyakubov/book_service_n/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type HttpHandler struct {
	bookService     book_Service.BookService
	activityHandler cache.ActivityCacher
}

func NewHttpHandler(bookService book_Service.BookService, activityHandler cache.ActivityCacher) HttpHandler {
	return HttpHandler{
		bookService:     bookService,
		activityHandler: activityHandler,
	}
}

type username struct {
	Username string `json:"username"`
}

func (h *HttpHandler) AddBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindBodyWith(&book, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, "error: couldn't bind body in AddBook request")
		return
	}

	id, err := h.bookService.InsertBook(context.Background(), &book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *HttpHandler) UpdateBook(c *gin.Context) {
	var book models.Book
	err := c.ShouldBindBodyWith(&book, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error: couldn't bind body in UpdateBook request")
		return
	}

	err = h.bookService.UpdateBook(context.Background(), book.Title, book.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, "book updated")
}

func (h *HttpHandler) GetBook(c *gin.Context) {
	s, err := h.bookService.GetBook(context.Background(), c.Query(elastic_fields.Id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) DeleteBook(c *gin.Context) {
	err := h.bookService.DeleteBook(context.Background(), c.Query(elastic_fields.Id))
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
	s, err := h.bookService.Search(context.Background(), searchParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *HttpHandler) StoreInfo(c *gin.Context) {
	info, err := h.bookService.StoreInfo(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"store_info": info,
	})
}

func (h *HttpHandler) Activity(c *gin.Context) {

	actions, err := h.activityHandler.GetLastActions(c.Query(query_fields.UserName))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, actions)
}

func (h *HttpHandler) AddAction(c *gin.Context) {
	un := username{c.Query(query_fields.UserName)}

	if c.Request.Method == "PUT" || c.Request.Method == "POST" {
		if err := c.ShouldBindBodyWith(&un, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, "error: couldn't bind body in Activity request")
			return
		}
	}

	if c.FullPath() != "/activity" {
		err := h.activityHandler.AddAction(un.Username, c.Request.Method, c.FullPath())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}
