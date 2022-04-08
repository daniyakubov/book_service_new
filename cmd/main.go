package main

import (
	"context"
	"github.com/daniyakubov/book_service_n/pkg/book_service"
	"github.com/daniyakubov/book_service_n/pkg/cache"
	"github.com/daniyakubov/book_service_n/pkg/consts"
	"github.com/daniyakubov/book_service_n/pkg/elastic_service"
	"github.com/daniyakubov/book_service_n/pkg/http_service"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"gopkg.in/redis.v5"
)

func main() {
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL(consts.BooksUrl))
	if err != nil {
		panic(err)
	}
	//	client := http.Client{Timeout: time.Duration(consts.ClientTimeOut) * time.Second}

	router := gin.Default()

	eHandler := elastic_service.NewElasticHandler(&ctx, consts.BooksUrl, client, consts.MaxQueryResults)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     consts.Host,
		Password: consts.Password,
		DB:       consts.Db,
	})
	bookService := book_Service.NewBookService(cache.NewRedisCache(consts.Host, consts.Db, consts.Expiration, consts.MaxActions, redisClient), &eHandler)
	httpHandler := http_service.NewHttpHandler(bookService)

	router.GET("/book", httpHandler.GetBook)
	router.DELETE("/book", httpHandler.DeleteBook)
	router.PUT("/book", httpHandler.PutBook)
	router.POST("/book", httpHandler.PostBook)
	router.GET("/search", httpHandler.Search)
	router.GET("/store", httpHandler.Store)
	router.GET("/activity", httpHandler.Activity)

	err = router.Run(consts.HttpAddress)
	if err != nil {
		panic(err)
	}
}
