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
func (e *ElasticHandler) Get(id string) (get *elastic.GetResult, err error) {

	get, err = e.Client.Get().
		Index("books").
		Id(id).
		Do(*e.Ctx)
	if err != nil {
		return get, errors.Wrap(err, err.Error())
	}

	return get, err
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

func (e *ElasticHandler) Search(title string, author string, priceStart float64, priceEnd float64) (res []models.Source, err error) {
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
	res = make([]models.Source, length)

	for i := 0; i < length; i++ {
		var src models.Source
		if err := json.Unmarshal(searchResult.Hits.Hits[i].Source, &src); err != nil {
			return nil, errors.Wrap(err, err.Error())
		}
		res[i] = src
		res[i].Id = searchResult.Hits.Hits[i].Id
	}

	return res, err
}

/*
func (e *ElasticHandler) Store() (resp1 *http.Response, resp2 *http.Response, err error) {
	s1 := fmt.Sprintf(`{"query": {"match_all": {}}}`)
	myJson := bytes.NewBuffer([]byte(s1))

	req, err := http.NewRequest(consts.GetMethod, e.Url+"_count/", myJson)
	if err != nil {
		return nil, nil, errors.Wrap(err, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp1, err = e.Client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, err.Error())
	}

	s2 := fmt.Sprintf(`{"aggs" : {"authors_count" : {"cardinality" : {"field" : "authorsName.keyword"}}}}`)
	myJson2 := bytes.NewBuffer([]byte(s2))

	resp2, err = e.Client.Post(e.Url+"_search/", "application/json", myJson2)
	if err != nil {
		return nil, nil, errors.Wrap(err, err.Error())
	}
	return resp1, resp2, nil
}
*/
