package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ElrondNetwork/elastic-search-tester/config"
	"github.com/ElrondNetwork/elastic-search-tester/types"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var log = logger.GetOrCreate("client")

type elasticClient struct {
	client *elasticsearch.Client
}

func NewElasticClient(cfg *config.Config) (*elasticClient, error) {
	esConfig := elasticsearch.Config{
		Addresses: []string{cfg.ElasticURL},
		Username:  cfg.User,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create database reader %w", err)
	}

	return &elasticClient{
		client: client,
	}, nil
}

func (ec *elasticClient) DoGetRequest(query *bytes.Buffer, index string, sort ...string) (types.ObjectMap, error) {
	res, err := ec.client.Search(
		ec.client.Search.WithIndex(index),
		ec.client.Search.WithBody(query),
		ec.client.Search.WithSort(sort...),
	)
	if err != nil {
		return nil, err
	}

	var decodedBody types.ObjectMap
	err = parseResponse(res, &decodedBody)
	if err != nil {
		return nil, err
	}

	return decodedBody, nil
}

func parseResponse(res *esapi.Response, dest interface{}) error {
	defer func() {
		if res != nil && res.Body != nil {
			err := res.Body.Close()
			log.LogIfError(err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot get response status code %d", res.StatusCode)
	}

	err := loadResponseBody(res.Body, dest)
	if err != nil {
		return err
	}

	return nil
}

func loadResponseBody(body io.ReadCloser, dest interface{}) error {
	if body == nil {
		return nil
	}
	if dest == nil {
		_, err := io.Copy(ioutil.Discard, body)
		return err
	}

	err := json.NewDecoder(body).Decode(dest)
	return err
}
