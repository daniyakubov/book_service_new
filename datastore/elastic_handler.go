package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/daniyakubov/book_service_n/config"
	"github.com/daniyakubov/book_service_n/consts"
	"github.com/daniyakubov/book_service_n/consts/elastic_fields"
	"github.com/daniyakubov/book_service_n/models"
	errors "github.com/fiverr/go_errors"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

var _ BookStorer = &ElasticHandler{}

type ElasticHandler struct {
	Client *elastic.Client
}

func NewElasticHandler(client *elastic.Client) ElasticHandler {
	return ElasticHandler{
		client,
	}
}

func (e *ElasticHandler) UpdateBook(ctx context.Context, title string, id string) (err error) {
	_, err = e.Client.Update().
		Index(consts.Index).
		Id(id).
		Doc(map[string]interface{}{elastic_fields.Title: title}).
		Do(ctx)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("couldn't update book with title %s and id %s", title, id))
	}

	return nil
}

func (e *ElasticHandler) InsertBook(ctx context.Context, body []byte) (string, error) {
	res, err := e.Client.Index().
		Index(consts.Index).
		BodyString(string(body)).
		Do(ctx)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("couldn't update book with body: %s", string(body)))
	}

	return res.Id, nil
}
func (e *ElasticHandler) GetBook(ctx context.Context, id string) (src *models.Book, err error) {
	get, err := e.Client.Get().
		Index(consts.Index).
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't get book form server with id: %s", id))
	}

	var book models.Book
	if err = json.Unmarshal(get.Source, &book); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't unmarshal result of book with id: %s, in getBoook function ", id))
	}

	return &book, nil
}

func (e *ElasticHandler) DeleteBook(ctx context.Context, id string) error {
	_, err := e.Client.Delete().
		Index(consts.Index).
		Id(id).
		Do(ctx)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("couldn't delete book with id: %s", id))
	}

	return nil
}

func buildSearchQuery(fields map[string]string) (s *elastic.BoolQuery, err error) {
	var from, to int
	if price, ok := fields["price_range"]; ok && price != "" {
		priceSplit := strings.Split(string(fields["price_range"]), "-")
		if len(priceSplit) != 2 {
			return nil, errors.New("failed to pars `price_range` field")
		}
		if from, err = strconv.Atoi(priceSplit[0]); err != nil {
			return nil, errors.New("failed to pars `price_range` field")
		}
		if to, err = strconv.Atoi(priceSplit[1]); err != nil {
			return nil, errors.New("failed to pars `price_range` field")
		}
	}

	if from < consts.MinPrice {
		return nil, errors.New(fmt.Sprintf("illegal price range, price should be higher than %d", consts.MinPrice))
	}
	if to > consts.MaxPrice {
		return nil, errors.New(fmt.Sprintf("illegal price range, price should be lower than %d", consts.MaxPrice))
	}

	q := elastic.NewBoolQuery()

	if fields["title"] != "" {
		q.Must(elastic.NewMatchQuery(elastic_fields.Title, fields["title"]))
	}
	if fields["author"] != "" {
		q.Must(elastic.NewMatchQuery(elastic_fields.Author, fields["author"]))
	}
	if to != 0 {
		q.Must(elastic.NewRangeQuery(elastic_fields.Price).From(from).To(to))
	}

	return q, nil
}

func (e *ElasticHandler) Search(ctx context.Context, fields map[string]string) ([]models.Book, error) {
	q, err := buildSearchQuery(fields)
	if err != nil {
		return nil, err
	}

	searchResult, err := e.Client.Search().
		Index(consts.Index).
		Query(q).
		From(0).
		Size(config.MaxQueryResults).
		Do(ctx)

	res := []models.Book{}
	for _, hit := range searchResult.Hits.Hits {
		var book models.Book
		if err := json.Unmarshal(hit.Source, &book); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unmarshaling in search failed for title: %s, author: %s, priceRange: %s", fields["title"], fields["author"], fields["priceRange"]))
		}
		book.Id = hit.Id
		res = append(res, book)
	}

	return res, err
}

func (e *ElasticHandler) StoreInfo(ctx context.Context) (info map[string]interface{}, err error) {
	searchResult, err := e.Client.Search().Aggregation("distinctAuthors", elastic.NewCardinalityAggregation().Field("author.keyword")).
		Index(consts.Index).
		Size(0).
		Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "store information retrieval failed ")
	}

	agg := searchResult.Aggregations
	if agg == nil {
		return nil, errors.New("agg returned nil in store information function")
	}

	agg2, found := agg.Cardinality("distinctAuthors")
	if !found {
		return nil, errors.New("aggregation was not found for distinct authors in store information function")
	}
	if agg2 == nil || agg2.Value == nil {
		return nil, errors.New("aggregation was nil for distinct authors in store information function")
	}
	info = make(map[string]interface{})
	info["books_num"] = searchResult.Hits.TotalHits.Value
	info["distinct_authors_num"] = int(*agg2.Value)
	return info, nil
}
