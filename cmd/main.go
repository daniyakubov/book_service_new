package main

import (
	"github.com/daniyakubov/book_service_n/pkg/book_service"
	"github.com/daniyakubov/book_service_n/pkg/cache"
	"github.com/daniyakubov/book_service_n/pkg/consts"
	"github.com/daniyakubov/book_service_n/pkg/elastic_service"
	"github.com/daniyakubov/book_service_n/pkg/http_service"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
	"net/http"
	"time"
)

func main() {
	client := http.Client{Timeout: time.Duration(consts.ClientTimeOut) * time.Second}

	router := gin.Default()

	eHandler := elastic_service.NewElasticHandler(consts.BooksUrl, &client, consts.MaxQueryResults)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     consts.Host,
		Password: consts.Password,
		DB:       consts.Db,
	})
	bookService := book_Service.NewBookService(cache.NewRedisCache(consts.Host, consts.Db, consts.Expiration, consts.MaxActions, redisClient), &eHandler)
	httpHandler := http_service.NewHttpHandler(client, bookService)

	router.GET("/book", httpHandler.GetBook)
	router.PUT("/book", httpHandler.PutBook)
	router.POST("/book", httpHandler.PostBook)
	router.DELETE("/book", httpHandler.DeleteBook)
	router.GET("/search", httpHandler.Search)
	router.GET("/store", httpHandler.Store)
	router.GET("/activity", httpHandler.Activity)

	err := router.Run(consts.HttpAddress)
	if err != nil {
		panic(err)
	}
}
