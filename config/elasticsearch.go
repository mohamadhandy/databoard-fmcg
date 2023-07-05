package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func SetupElasticsearch() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

func CreateIndex(client *elasticsearch.Client, indexName string) error {
	// Create the index request
	req := esapi.IndicesCreateRequest{
		Index: indexName,
	}

	// Execute the request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the response status
	if res.IsError() {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("failed to create index: %s\nError: %s", res.Status(), string(bodyBytes))
	}

	// Index creation successful
	return nil
}

func IndexExists(client *elasticsearch.Client, indexName string) (bool, error) {
	// Create the index exists request
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	// Execute the request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// Check the response status
	if res.IsError() {
		return false, fmt.Errorf("failed to check index existence: %s", res.Status())
	}

	// Parse the response
	var respBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		return false, fmt.Errorf("failed to parse index existence response: %s", err)
	}

	// Check if the index exists
	exists, ok := respBody[indexName].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected response structure while checking index existence")
	}

	return exists, nil
}
