package main

import (
	"context"
	"klikdaily-databoard/config"
	"klikdaily-databoard/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	router := gin.Default()

	// router.POST("/users", authMiddleware, userController.RegisterUser)
	// router.POST("/login", userController.LoginUser)
	// router.GET("/users", authMiddleware, userController.GetAllUsers)
	// router.GET("/users/:user_id", authMiddleware, userController.GetUsersByID)
	// router.PUT("/users/:user_id", authMiddleware, userController.UpdateUser)
	// router.DELETE("/users/:user_id", authMiddleware, userController.DeleteUser)
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
	// RabbitMQ connection details
	rabbitMQURL := "amqp://guest:guest@localhost:5672/"
	rabbitMQQueue := "postgre_sync_queue"
	// Connect to RabbitMQ
	rabbitMQConnection, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQConnection.Close()

	// Create a RabbitMQ channel
	ch, err := rabbitMQConnection.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

	// Declare the RabbitMQ queue
	_, err = ch.QueueDeclare(
		rabbitMQQueue, // Queue name
		true,          // Durable
		false,         // Auto-delete
		false,         // Exclusive
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		log.Fatal("Failed to declare RabbitMQ queue:", err)
	}
	routes.RouteAPI(router, context.Background(), config.NewConnection(), config.NewConnectionRedis(), es)
	router.Run("localhost:9000")
}
