package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/avigailbr/book_service/connectors"
	"github.com/avigailbr/book_service/models"
	"gopkg.in/olivere/elastic.v5"
	"net/http"
	"strconv"
	"strings"
)

const (
	SearchResultFrom = 0
	SearchResultSize = 10
)

type ElasticBookStore struct {
	Client *connectors.ElasticClient
}

func NewElasticBookStore(url, index, doctype string) (IBookStorer, error) {
	esClient, err := connectors.NewElasticClient(url, index, doctype)
	if err != nil {
		return nil, err
	}
	esBookStore := &ElasticBookStore{
		Client: esClient,
	}
	return esBookStore, nil
}

func (e ElasticBookStore) GetBook(id string) (*models.Book, error) {
	es := e.Client
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

func (e ElasticBookStore) InsertBook(book *models.Book) (string, error) {
	es := e.Client
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

func (e ElasticBookStore) UpdateBook(id string, fields map[string]interface{}) error {
	es := e.Client
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

func (e ElasticBookStore) DeleteBook(id string) error {
	es := e.Client
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

func (e ElasticBookStore) Search(fields map[string]string) ([]models.Book, error) {

	query := elastic.NewBoolQuery()

	if title, ok := fields["title"]; ok && title != "" {
		titleMatchQuery := elastic.NewMatchQuery("title", title)
		query.Must(titleMatchQuery)
	}

	if author, ok := fields["author_name"]; ok && author != "" {
		authorMatchQuery := elastic.NewMatchQuery("author_name", author)
		query.Must(authorMatchQuery)
	}
	if price, ok := fields["price_range"]; ok && price != "" {
		var from, to int
		var err error
		s := strings.Split(price, "-")
		if len(s) != 2 {
			return nil, errors.New("Conversion failed for `Price_range` field")
		}
		if from, err = strconv.Atoi(s[0]); err != nil {
			return nil, errors.New("Conversion failed for `Price_range` field")
		}
		if to, err = strconv.Atoi(s[1]); err != nil {
			return nil, errors.New("Conversion failed for `Price_range` field")
		}
		priceFilterQuery := elastic.NewRangeQuery("price").From(from).To(to)
		query.Filter(priceFilterQuery)
	}

	es := e.Client
	ctx := context.Background()

	searchResult, err := es.Client.Search().
		Index(es.Index).
		Type(es.DocType).
		Query(query).
		From(SearchResultFrom).
		Size(SearchResultSize).
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

func (e ElasticBookStore) Info() (map[string]interface{}, error) {

	info := make(map[string]interface{})

	es := e.Client
	ctx := context.Background()

	sr, err := es.Client.Search().Aggregation("author_count", elastic.NewCardinalityAggregation().Field("author_name.keyword")).
		Index(es.Index).
		Type(es.DocType).
		Size(0).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	authorsCount, found := sr.Aggregations["author_count"]
	if found {
		info["author_count"] = authorsCount
		info["books_count"] = sr.Hits.TotalHits

	} else {
		return nil, &elastic.Error{Status: http.StatusNotFound}
	}

	return info, nil
}
