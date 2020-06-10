package controllers

import (
	"context"
	"encoding/json"
	"github.com/avigailbr/book_service/config"
	"github.com/avigailbr/book_service/connectors"
	"github.com/avigailbr/book_service/models"
	"gopkg.in/olivere/elastic.v5"
	"net/http"
	"strconv"
	"strings"
)

var (
	BookController IBookController
)

type IBookController interface {
	GetBook(id string) (*models.Book, error)
	InsertBook(book *models.Book) (string, error)
	UpdateBook(id string, fields map[string]interface{}) error
	DeleteBook(id string) error
	Search(fields map[string]string) ([]models.Book, error)
	Info() (map[string]interface{}, error)
}

type ElasticBookController struct {
}

func NewElasticBookController() IBookController {
	return &ElasticBookController{}
}

func (e ElasticBookController) GetBook(id string) (*models.Book, error) {
	es := connectors.BookConnectors.ElasticClient
	ctx := context.Background()
	get, err := es.Client.Get().
		Index(es.Index).
		Type(es.DocType).
		Id(id).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var book models.Book
	if data, err := json.Marshal(get.Source); err != nil {
		return nil, err
	} else {
		json.Unmarshal(data, &book)
		return &book, nil
	}

}

func (e ElasticBookController) InsertBook(book *models.Book) (string, error) {
	es := connectors.BookConnectors.ElasticClient
	ctx := context.Background()

	data, err := json.Marshal(book)
	if err != nil {
		return "", err
	}

	put, err := es.Client.Index().
		Index(es.Index).
		Type(es.DocType).
		BodyJson(string(data)).
		Do(ctx)

	if err != nil {
		return "", err
	}

	return put.Id, nil

}

func (e ElasticBookController) UpdateBook(id string, fields map[string]interface{}) error {
	es := connectors.BookConnectors.ElasticClient
	ctx := context.Background()

	_, err := es.Client.Update().
		Index(es.Index).
		Type(es.DocType).
		Id(id).
		Doc(fields).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (e ElasticBookController) DeleteBook(id string) error {
	es := connectors.BookConnectors.ElasticClient
	// must pass a context to execute each service
	ctx := context.Background()
	_, err := es.Client.Delete().
		Index(es.Index).
		Type(es.DocType).
		Id(id).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e ElasticBookController) Search(fields map[string]string) ([]models.Book, error) {

	query := elastic.NewBoolQuery()

	if title, ok := fields["title"]; ok && title != "" {
		titleMatchQuery := elastic.NewMatchQuery("title", title)

		if author, ok := fields["author_name"]; ok && author != "" {
			authorMatchQuery := elastic.NewMatchQuery("author_name", author)
			query.Must(titleMatchQuery, authorMatchQuery)

		} else {
			query.Must(titleMatchQuery)
		}

	} else if author, ok := fields["author_name"]; ok && author != "" {
		authorMatchQuery := elastic.NewMatchQuery("author_name", author)
		query.Must(authorMatchQuery)
	}
	if price, ok := fields["price_range"]; ok && price != "" {
		var from, to int
		var err error
		s := strings.Split(price, "-")
		if len(s) != 2 {
			return nil, models.NewStringError("Conversion failed for `Price_range` field")
		}
		if from, err = strconv.Atoi(s[0]); err != nil {
			return nil, models.NewStringError("Conversion failed for `Price_range` field")
		}
		if to, err = strconv.Atoi(s[1]); err != nil {
			return nil, models.NewStringError("Conversion failed for `Price_range` field")
		}
		priceFilterQuery := elastic.NewRangeQuery("price").From(from).To(to)
		query.Filter(priceFilterQuery)
	}

	es := connectors.BookConnectors.ElasticClient
	ctx := context.Background()

	searchResult, err := es.Client.Search().
		Index(es.Index).
		Type(es.DocType).
		Query(query).
		From(config.SearchResultFrom).
		Size(config.SearchResultSize).
		Pretty(true).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var books []models.Book

	for _, hit := range searchResult.Hits.Hits {
		var book models.Book
		data, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(data, &book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil

}

func (e ElasticBookController) Info() (map[string]interface{}, error) {

	info := make(map[string]interface{})

	es := connectors.BookConnectors.ElasticClient
	ctx := context.Background()

	sr, err := es.Client.Search().Aggregation("author_count", elastic.NewCardinalityAggregation().Field("author_name.keyword")).
		Index(es.Index).
		Type(es.DocType).
		Size(0).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var authorsCount int
	authors, found := sr.Aggregations.Terms("author_count")
	if found {
		for _, count := range authors.Aggregations {
			if i, err := json.Marshal(&count); err != nil {
				return nil, err
			} else if authorsCount, err = strconv.Atoi(string(i)); err != nil {
				return nil, err
			}
		}
		info["author_count"] = authorsCount
		info["books_count"] = sr.Hits.TotalHits

	} else {
		return nil, &elastic.Error{Status: http.StatusNotFound}
	}

	return info, nil
}
