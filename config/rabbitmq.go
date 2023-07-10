package config

import (
	"log"

	"github.com/streadway/amqp"
)

func SetupRabbitMQ() (*amqp.Channel, error) {
	// RabbitMQ connection details
	rabbitMQURL := "amqp://guest:guest@localhost:5672/"
	rabbitMQQueue := "postgre_sync_queue"
	// Connect to RabbitMQ
	rabbitMQConnection, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
		return nil, err
	}
	defer rabbitMQConnection.Close()

	// Create a RabbitMQ channel
	ch, err := rabbitMQConnection.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
		return nil, err
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
		return nil, err
	}
	return ch, nil
}
