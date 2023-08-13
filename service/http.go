package service

import (
	"context"
	"klikdaily-databoard/config"
	"klikdaily-databoard/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func LoadHttpService() *gin.Engine {
	router := gin.Default()

	es, err := config.SetupElasticsearch()
	if err != nil {
		panic(err)
	}
	// Create the Elasticsearch index
	indexName := "elasticsearch_index"
	// Check if the index exists
	exists, err := config.IndexExists(es, indexName)
	if err != nil {
		log.Fatalf("Error checking index existence: %s", err)
	}

	if !exists {
		// Create the Elasticsearch index
		err = config.CreateIndex(es, indexName)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
		log.Printf("Index %s created", indexName)
	} else {
		log.Printf("Index %s already exists", indexName)
	}

	mb, err := config.NewConnectionRabbitMQ()

	if err != nil {
		log.Printf("%q", err)
	}

	routes.RouteAPI(router, context.Background(), config.NewConnection(), config.NewConnectionRedis(), es, mb)
	router.Run("localhost:9000")

	return router
}
