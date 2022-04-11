package elastic_service

import (
	"context"
	"encoding/json"
	"github.com/daniyakubov/book_service_n/config"
	"github.com/daniyakubov/book_service_n/consts"
	"github.com/daniyakubov/book_service_n/db_service"
	"github.com/daniyakubov/book_service_n/models"
	errors "github.com/fiverr/go_errors"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

var _ db_service.DbHandler = &ElasticHandler{}

type ElasticHandler struct {
	Url          string
	Client       *elastic.Client
	maxSizeQuery int
}

func NewElasticHandler(url string, client *elastic.Client, maxSizeQuery int) ElasticHandler {
	return ElasticHandler{
		url,
		client,
		maxSizeQuery,
	}
}

func (e *ElasticHandler) Update(ctx *context.Context, title string, id string) (err error) {
	_, err = e.Client.Update().
		Index(consts.Index).
		Id(id).
		Doc(map[string]interface{}{consts.Title: title}).
		Do(*ctx)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	return nil
}

func (e *ElasticHandler) Add(ctx *context.Context, body []byte) (string, error) {
	res, err := e.Client.Index().
		Index(consts.Index).
		BodyString(string(body)).
		Do(*ctx)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	return res.Id, err
}
func (e *ElasticHandler) Get(ctx *context.Context, id string) (src *models.Book, err error) {
	get, err := e.Client.Get().
		Index(consts.Index).
		Id(id).
		Do(*ctx)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	var book models.Book
	if err = json.Unmarshal(get.Source, &book); err != nil {
		return nil, err
	}
	return &book, err
}

func (e *ElasticHandler) Delete(ctx *context.Context, id string) error {
	_, err := e.Client.Delete().
		Index(consts.Index).
		Id(id).
		Do(*ctx)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	return err
}

func buildQueryForSearch(client *elastic.Client, title string, author string, priceStart float64, priceEnd float64) *elastic.SearchService {
	all := elastic.NewMatchAllQuery()
	builder := client.Search().Index(consts.Index).Query(all).Pretty(true)

	if title != "" {
		builder = builder.Query(elastic.NewBoolQuery().Must(elastic.NewTermQuery(consts.Title, title)))
	}
	if author != "" {
		builder = builder.Query(elastic.NewBoolQuery().Must(elastic.NewTermQuery(consts.Author, author)))
	}
	if priceEnd != 0 {
		builder = builder.Query(elastic.NewRangeQuery(consts.Price).From(priceStart).To(priceEnd))
	}
	return builder
}

func (e *ElasticHandler) Search(title string, author string, priceRange string) (res []models.Book, err error) {
	var priceStart float64
	var priceEnd float64
	if priceRange == "" {
		priceStart = 0
		priceEnd = 0
	} else {
		priceSplit := strings.Split(priceRange, "-")
		priceStart, _ = strconv.ParseFloat(priceSplit[0], 32)
		priceEnd, _ = strconv.ParseFloat(priceSplit[1], 32)
	}

	builder := buildQueryForSearch(e.Client, title, author, priceStart, priceEnd)
	searchResult, err := builder.Pretty(true).Size(config.MaxQueryResults).Do(context.TODO())

	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	if searchResult.Hits == nil {
		return nil, errors.New("expected SearchResult.Hits != nil; got nil")
	}

	for _, s := range searchResult.Hits.Hits {
		var src models.Book
		if err := json.Unmarshal(s.Source, &src); err != nil {
			return nil, errors.Wrap(err, err.Error())
		}
		src.Id = s.Id
		res = append(res, src)
	}
	return res, err
}

func (e *ElasticHandler) StoreInfo() (int64, int, error) {
	all := elastic.NewMatchAllQuery()
	cardinalityAgg := elastic.NewCardinalityAggregation().Field("author.keyword")
	builder := e.Client.Search().Index(consts.Index).Query(all).Pretty(true)
	builder = builder.Aggregation("distinctAuthors", cardinalityAgg)
	searchResult, err := builder.Pretty(true).Do(context.TODO())
	if err != nil {
		return 0, 0, errors.Wrap(err, err.Error())
	}
	agg := searchResult.Aggregations
	if agg == nil {
		return 0, 0, errors.New("agg returned nil")
	}
	agg2, found := agg.Cardinality("distinctAuthors")
	if !found {
		return 0, 0, errors.New("not found")
	}
	if agg2 == nil || agg2.Value == nil {
		return 0, 0, errors.New("expected != nil; got: nil")
	}
	booksNum := searchResult.Hits.TotalHits.Value
	return booksNum, int(*agg2.Value), nil
}
