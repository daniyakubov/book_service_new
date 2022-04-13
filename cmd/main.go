package main

import (
	"fmt"
	"github.com/daniyakubov/book_service_n/book_service"
	"github.com/daniyakubov/book_service_n/cache"
	"github.com/daniyakubov/book_service_n/config"
	"github.com/daniyakubov/book_service_n/datastore"
	"github.com/daniyakubov/book_service_n/http_service"
	"github.com/daniyakubov/book_service_n/service"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"gopkg.in/redis.v5"
)

func main() {
	client, err := elastic.NewClient(elastic.SetURL(config.BooksUrl))
	if err != nil {
		panic(err)
	}
	eHandler := datastore.NewElasticHandler(client)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.Password,
		DB:       config.DB,
	})
	bookService := book_Service.NewBookService(cache.NewRedisCache(config.RedisAddress, config.DB, config.Expiration, config.MaxActions, redisClient), &eHandler)
	httpHandler := http_service.NewHttpHandler(bookService)

	router := gin.Default()
	router = service.Routes(router, &httpHandler)

	err = router.Run(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		panic(err)
	}
}
