package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"klikdaily-databoard/definition"
	"klikdaily-databoard/models"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	amqp "github.com/rabbitmq/amqp091-go"
)

func LoadIndexerService() (<-chan error, error) {
	log.Println("Indexer Service mode!")
	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Printf("Error creating Elasticsearch client: %s", err)
		return nil, err
	}

	rmq, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %s", err)
		return nil, err
	}

	ch, err := rmq.Channel()
	if err != nil {
		log.Fatalf("Error creating RabbitMQ channel: %s", err)
		return nil, err
	}

	srv := &IndexerService{
		esClient: esClient,
		rmq:      rmq,
		ch:       ch,
		done:     make(chan struct{}),
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Println("Shutting down...")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			rmq.Close()
			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Println("Shutted down!")
	}()

	go func() {
		log.Println("Listening...")

		if err := srv.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}

type IndexerService struct {
	esClient *elasticsearch.Client
	rmq      *amqp.Connection
	ch       *amqp.Channel
	done     chan struct{}
}

func (is *IndexerService) ListenAndServe() error {
	err := is.ch.ExchangeDeclare(
		definition.ExchangeName, // name
		"topic",                 // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring exchage RabbitMQ channel: %s", err)
		return err
	}

	q, err := is.ch.QueueDeclare(
		definition.QueueName, // name
		true,                 // durable
		false,                // delete when unused
		true,                 // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatalf("Error queue declaration from RabbitMQ: %s", err)
		return err
	}

	err = is.ch.QueueBind(
		q.Name,                  // queue name
		"product.*",             // routing key
		definition.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error queue bind from RabbitMQ: %s", err)
		return err
	}

	msgs, err := is.ch.Consume(
		q.Name,
		definition.ConsumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error consuming messages from RabbitMQ: %s", err)
		return err
	}

	go func() {
		log.Printf("Waiting for messages...")
		for msg := range msgs {
			var nack bool

			log.Printf("Received message for \"%s\": %s", msg.RoutingKey, msg.Body)

			switch msg.RoutingKey {
			case "product.created":
				syncProductIndex(is.esClient, msg.Body)

				if err != nil {
					nack = true
					return
				}

			default:
				nack = true
			}

			if nack {
				log.Printf("NAck!")
				msg.Nack(false, nack)
			} else {
				log.Printf("Ack!")
				msg.Ack(false)
			}
		}

		log.Printf("No more message received!")

		is.done <- struct{}{}
	}()

	return nil
}

func (is *IndexerService) Shutdown(ctx context.Context) error {
	is.ch.Cancel(definition.ConsumerName, false)

	for {
		select {
		case <-ctx.Done():
			log.Println("Done by context...")
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-is.done:
			log.Println("Done by service...")
			return nil
		}
	}
}

func syncProductIndex(esClient *elasticsearch.Client, rawProduct []byte) error {
	var product models.Product
	// Unmarshal pesan JSON menjadi objek Product
	if err := json.Unmarshal(rawProduct, &product); err != nil {
		log.Printf("Error unmarshaling product from JSON: %s", err)
		return err
	}
	// Prepare the product document to be indexed in Elasticsearch
	productDoc := map[string]interface{}{
		"id":   product.ID,
		"name": product.Name, // Replace with the actual product name
		// Add more fields as needed
	}

	productJSON, err := json.Marshal(productDoc)
	if err != nil {
		log.Printf("Error marshaling product document: %s", err)
		return err
	}

	// Sync product data to Elasticsearch
	req := esapi.IndexRequest{
		Index:      "products", // Replace with the actual index name
		DocumentID: product.ID,
		Body:       bytes.NewReader(productJSON),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Printf("Error indexing document in Elasticsearch: %s", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing document in Elasticsearch: [%s] %s", res.Status(), res.String())
		return err
	}

	log.Println("Product data synchronized to Elasticsearch successfully.")
	return nil
}
