package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"klikdaily-databoard/config"
	"klikdaily-databoard/definition"
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	CreateProduct(pr models.ProductRequest, tokenString string) chan RepositoryResult[any]
	GetProductById(id string) chan RepositoryResult[any]
	GetProducts(productReq models.ProductRequest, searchKeyword string) chan RepositoryResult[any]
	UpdateProduct(tokenString string, productReq models.ProductRequest) chan RepositoryResult[any]
	GetPreviousId() string
}

type productRepository struct {
	db *gorm.DB
	es *elasticsearch.Client
	mb *config.MessageBroker
}

func InitProductRepository(db *gorm.DB, es *elasticsearch.Client, mb *config.MessageBroker) ProductRepositoryInterface {
	return &productRepository{
		db,
		es,
		mb,
	}
}

func (r *productRepository) getBodyBytes(query map[string]interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	return &buf
}

func (r *productRepository) GetProducts(productReq models.ProductRequest, searchKeyword string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])

	go func() {
		// Build the Elasticsearch query
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"name": searchKeyword,
				},
			},
		}

		// Execute the Elasticsearch search
		res, err := r.es.Search(
			r.es.Search.WithIndex("products"),
			r.es.Search.WithBody(r.getBodyBytes(query)),
		)
		if err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		defer res.Body.Close()

		// Parse the search response and extract the products
		products := []models.ProductResponse{}
		if res.IsError() {
			result <- RepositoryResult[any]{
				Data:       products,
				Error:      fmt.Errorf("Elasticsearch search error: %s", res.Status()),
				Message:    "Failed to retrieve products",
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		var response map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			result <- RepositoryResult[any]{
				Data:       products,
				Error:      fmt.Errorf("Error parsing Elasticsearch response: %s", err),
				Message:    "Failed to retrieve products: 12",
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		hits, ok := response["hits"].(map[string]interface{})["hits"].([]interface{})
		if !ok {
			result <- RepositoryResult[any]{
				Data:       products,
				Error:      fmt.Errorf("Unexpected Elasticsearch response structure"),
				Message:    "Failed to retrieve products 13",
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		for _, hit := range hits {
			source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
			if !ok {
				continue
			}

			product := models.ProductResponse{
				ID:   "test",
				Name: source["name"].(string),
				// Extract other fields as needed
			}

			products = append(products, product)
		}

		result <- RepositoryResult[any]{
			Data:       products,
			Error:      nil,
			Message:    "Success Get Products",
			StatusCode: http.StatusOK,
		}
	}()
	return result
}

func (r *productRepository) GetPreviousId() string {
	latestID := ""
	if err := r.db.Model(&models.Product{}).Select("id").Order("id desc").Limit(1).Scan(&latestID).Error; err != nil {
		return "error " + err.Error()
	}
	return latestID
}

func (r *productRepository) GetProductById(id string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		product := models.Product{}
		if err := r.db.Where("id = ?", id).Find(&product).Error; err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		if product.Name == "" {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      errors.New("product not found"),
				Message:    "Product Not found",
				StatusCode: http.StatusNotFound,
			}
		}
		result <- RepositoryResult[any]{
			Data:       product,
			Error:      nil,
			Message:    "Success get product by id: " + id,
			StatusCode: http.StatusOK,
		}
	}()
	return result
}

func (r *productRepository) UpdateProduct(tokenString string, productReq models.ProductRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		userName := helper.ExtractUserIDFromToken(tokenString)
		product := models.Product{
			Name:       productReq.Name,
			ID:         productReq.ID,
			BrandId:    productReq.BrandId,
			CategoryId: productReq.CategoryId,
			UpdatedBy:  userName,
			Status:     "active",
			CreatedBy:  productReq.CreatedBy,
			CreatedAt:  productReq.CreatedAt,
			SKU:        productReq.SKU,
		}
		tx := r.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      errors.New("panic occured"),
					StatusCode: http.StatusInternalServerError,
					Message:    "An unexpected token",
				}
				return
			}
		}()

		err := tx.Transaction(func(tx *gorm.DB) error {
			if err := r.db.Save(&product).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			tx.Rollback()
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       product,
			Error:      nil,
			Message:    "Success Update Product",
			StatusCode: http.StatusOK,
		}
	}()
	return result
}

func (r *productRepository) CreateProduct(pr models.ProductRequest, tokenString string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		tx := r.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      errors.New("panic occurred"),
					Message:    "An unexpected error occurred",
					StatusCode: http.StatusInternalServerError,
				}
			}
		}()

		latestID := r.GetPreviousId()
		if strings.Contains(latestID, "error") {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      errors.New(latestID),
				StatusCode: http.StatusInternalServerError,
				Message:    latestID,
			}
			return
		}
		latestOnlyId := helper.SplitProductID(latestID)
		latestIdInt, err := strconv.Atoi(latestOnlyId)
		if err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		productId, latestIdString := helper.GenerateProductID(latestIdInt)
		userName := helper.ExtractUserIDFromToken(tokenString)
		product := models.Product{
			ID:         productId,
			Name:       pr.Name,
			BrandId:    pr.BrandId,
			Status:     "active",
			CategoryId: pr.CategoryId,
			CreatedBy:  userName,
			UpdatedBy:  userName,
			SKU:        pr.BrandId + pr.CategoryId + latestIdString,
		}

		// Create the product within the transaction
		if err := r.db.Create(&product).Error; err != nil {
			tx.Rollback()
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    "Error: " + err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		// Start the Elasticsearch consumer
		// go startElasticsearchConsumer()

		if err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    "Error declaring exchange to RabbitMQ: " + err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		productJSON, err := json.Marshal(product)
		if err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    "Error marshaling product to JSON: " + err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		confirmations := r.mb.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer r.mb.HandleConfirmation(confirmations)

		// Publish the product ID to RabbitMQ
		err = r.mb.Channel.Publish(
			definition.ExchangeName,
			"product.created",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        productJSON,
			},
		)
		if err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    "Error publishing message to RabbitMQ: " + err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    "Error: " + err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		// Send the successful response
		result <- RepositoryResult[any]{
			Data:       product,
			Error:      nil,
			Message:    "Create Product Success",
			StatusCode: http.StatusCreated,
		}
	}()
	return result
}

func startElasticsearchConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating RabbitMQ channel: %s", err)
	}
	defer ch.Close()

	queueName := "product_sync_queue"
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error consuming messages from RabbitMQ: %s", err)
	}

	for msg := range msgs {
		// Sync product data to Elasticsearch
		syncProductToElasticsearch(msg.Body)
		log.Printf("Received message: %s", msg.Body)
	}
}

func syncProductToElasticsearch(msg []byte) {
	var product models.Product

	// Unmarshal pesan JSON menjadi objek Product
	if err := json.Unmarshal(msg, &product); err != nil {
		log.Printf("Error unmarshaling product from JSON: %s", err)
		return
	}

	// Create an Elasticsearch client
	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Printf("Error creating Elasticsearch client: %s", err)
		return
	}

	// Prepare the product document to be indexed in Elasticsearch
	productDoc := map[string]interface{}{
		"id":   product.ID,
		"name": product.Name, // Replace with the actual product name
		// Add more fields as needed
	}

	// Convert the product document to JSON
	productJSON, err := json.Marshal(productDoc)
	if err != nil {
		log.Printf("Error marshaling product document: %s", err)
		return
	}

	// Prepare the Elasticsearch index request
	req := esapi.IndexRequest{
		Index:      "products", // Replace with the actual index name
		DocumentID: product.ID,
		Body:       bytes.NewReader(productJSON),
		Refresh:    "true",
	}

	// Perform the Elasticsearch index request
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Printf("Error indexing document in Elasticsearch: %s", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing document in Elasticsearch: [%s] %s", res.Status(), res.String())
		return
	}

	log.Println("Product data synchronized to Elasticsearch successfully.")
}
