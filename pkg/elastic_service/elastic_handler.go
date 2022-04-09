package elastic_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/daniyakubov/book_service_n/pkg/book_service/models"
	errors "github.com/fiverr/go_errors"
	"github.com/olivere/elastic/v7"
)

type ElasticHandler struct {
	Ctx          *context.Context
	Url          string
	Client       *elastic.Client
	maxSizeQuery int
}

func NewElasticHandler(ctx *context.Context, url string, client *elastic.Client, maxSizeQuery int) ElasticHandler {
	return ElasticHandler{
		ctx,
		url,
		client,
		maxSizeQuery,
	}
}

func (e *ElasticHandler) Post(title string, id string) (err error) {
	_, err = e.Client.Update().Index("books").Id(id).Doc(map[string]interface{}{"title": title}).Do(*e.Ctx)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	return nil
}

func (e *ElasticHandler) Put(postBody []byte) (string, error) {
	put, err := e.Client.Index().
		Index("books").
		BodyString(string(postBody)).
		Do(*e.Ctx)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	return put.Id, err
}
func (e *ElasticHandler) Get(id string) (src *models.Book, err error) {
	get, err := e.Client.Get().
		Index("books").
		Id(id).
		Do(*e.Ctx)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	var getResp models.Book
	if err = json.Unmarshal(get.Source, &getResp); err != nil {
		return nil, err
	}
	return &getResp, err
}

func (e *ElasticHandler) Delete(id string) error {
	_, err := e.Client.Delete().
		Index("books").
		Id(id).
		Do(*e.Ctx)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	return err
}

func (e *ElasticHandler) Search(title string, author string, priceStart float64, priceEnd float64) (res []models.Book, err error) {
	s := ""
	if title == "" && author == "" && priceEnd == 0 {
		s = fmt.Sprintf(`{"match_all": {}}`)
	} else if title == "" && author == "" {
		s = fmt.Sprintf(`{"range": {"price": {"gte": %f, "lte": %f}}}`, priceStart, priceEnd)
	} else if title == "" && priceEnd == 0 {
		s = fmt.Sprintf(`{"constant_score": {"filter": {"bool": {"must":[{"match": {"author": "%s"}}]}}}}`, author)
	} else if author == "" && priceEnd == 0 {
		s = fmt.Sprintf(`{"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}}]}}}}`, title)
	} else if priceEnd == 0 {
		s = fmt.Sprintf(`{"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"match": {"author": "%s"}}]}}}}`, title, author)
	} else if title == "" {
		s = fmt.Sprintf(`{"constant_score": {"filter": {"bool": {"must":[{"match": {"author": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}`, author, priceStart, priceEnd)
	} else if author == "" {
		s = fmt.Sprintf(`{"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}`, title, priceStart, priceEnd)
	} else {
		s = fmt.Sprintf(`{"constant_score": {"filter": {"bool": {"must":[{"match": {"title": "%s"}},{"match": {"author": "%s"}},{"range": {"price": {"gte": %f, "lte": %f} }}]}}}}`, title, author, priceStart, priceEnd)
	}

	query := elastic.RawStringQuery(s)
	searchResult, err := e.Client.Search().
		Index("books").
		Query(query).
		Size(e.maxSizeQuery).
		Do(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	if searchResult.Hits == nil {
		return nil, errors.New("expected SearchResult.Hits != nil; got nil")
	}

	length := len(searchResult.Hits.Hits)
	res = make([]models.Book, length)
	for i := 0; i < length; i++ {
		var src models.Book
		if err := json.Unmarshal(searchResult.Hits.Hits[i].Source, &src); err != nil {
			return nil, errors.Wrap(err, err.Error())
		}
		res[i] = src
		res[i].Id = searchResult.Hits.Hits[i].Id
	}
	return res, err
}

func (e *ElasticHandler) Store() (int64, int, error) {
	count, err := e.Client.Count("books").Do(context.TODO())
	if err != nil {
		return 0, 0, errors.Wrap(err, err.Error())
	}

	all := elastic.NewMatchAllQuery()
	cardinalityAgg := elastic.NewCardinalityAggregation().Field("author.keyword")
	builder := e.Client.Search().Index("books").Query(all).Pretty(true)
	builder = builder.Aggregation("distinctAuthors", cardinalityAgg)
	searchResult, err := builder.Pretty(true).Do(context.TODO())
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
	return count, int(*agg2.Value), nil
}
