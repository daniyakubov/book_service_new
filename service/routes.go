package service

import (
	"github.com/daniyakubov/book_service_n/http_service"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine, h *http_service.HttpHandler) *gin.Engine {
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
	return router
}
