package main

import (
	"github.com/daniyakubov/book_service_n/book_service"
	"github.com/daniyakubov/book_service_n/cache"
	"github.com/daniyakubov/book_service_n/config"
	"github.com/daniyakubov/book_service_n/elastic_service"
	"github.com/daniyakubov/book_service_n/http_service"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"gopkg.in/redis.v5"
)

func main() {
	client, err := elastic.NewClient(elastic.SetURL(config.BooksUrl))
	if err != nil {
		panic(err)
	}
	eHandler := elastic_service.NewElasticHandler(config.BooksUrl, client, config.MaxQueryResults)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       config.Db,
	})
	bookService := book_Service.NewBookService(cache.NewRedisCache(config.Host, config.Db, config.Expiration, config.MaxActions, redisClient), &eHandler)
	httpHandler := http_service.NewHttpHandler(bookService)

	router := gin.Default()

	book := router.Group("/book")
	{
		book.GET("", httpHandler.GetBook)
		book.DELETE("", httpHandler.DeleteBook)
		book.PUT("", httpHandler.PutBook)
		book.POST("", httpHandler.PostBook)
	}
	router.GET("/search", httpHandler.Search)
	router.GET("/store", httpHandler.StoreInfo)
	router.GET("/activity", httpHandler.Activity)

	err = router.Run(config.HttpAddress)
	if err != nil {
		panic(err)
	}
}
