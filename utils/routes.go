package utils

import (
	"github.com/daniyakubov/book_service_n/app_service"
	"github.com/gin-gonic/gin"
)

type RoutesService struct {
	appHandler app_service.AppHandler
}

func NewRoutesService(appHandler app_service.AppHandler) *RoutesService {
	return &RoutesService{
		appHandler: appHandler,
	}
}

func (r *RoutesService) Routes(router *gin.Engine) *gin.Engine {
	book := router.Group("/book")
	{
		book.GET("", r.appHandler.GetBook)
		book.DELETE("", r.appHandler.DeleteBook)
		book.PUT("", r.appHandler.AddBook)
		book.POST("", r.appHandler.UpdateBook)
	}
	router.GET("/search", r.appHandler.Search)
	router.GET("/store", r.appHandler.StoreInfo)
	router.GET("/activity", r.appHandler.Activity)
	return router
}
