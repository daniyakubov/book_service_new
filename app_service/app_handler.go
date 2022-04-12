package app_service

import "github.com/gin-gonic/gin"

type AppHandler interface {
	AddBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	GetBook(c *gin.Context)
	DeleteBook(c *gin.Context)
	Search(c *gin.Context)
	StoreInfo(c *gin.Context)
	Activity(c *gin.Context)
}
