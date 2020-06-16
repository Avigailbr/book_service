package connectors

import (
	"gopkg.in/olivere/elastic.v5"
)

type ElasticClient struct {
	Client  *elastic.Client
	Index   string
	DocType string
}

func NewElasticClient(url, index, doctype string) (*ElasticClient, error) {

	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	elastic := &ElasticClient{
		Client:  client,
		Index:   index,
		DocType: doctype,
	}

	return elastic, nil
}
