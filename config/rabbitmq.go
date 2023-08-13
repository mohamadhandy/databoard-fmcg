package config

import (
	"klikdaily-databoard/definition"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewConnectionRabbitMQ() (*MessageBroker, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		definition.ExchangeName, // name
		"topic",                 // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return nil, err
	}

	if err = ch.Confirm(false); err != nil {
		return nil, err
	}

	return &MessageBroker{
		Connection: conn,
		Channel:    ch,
	}, err
}

func (m MessageBroker) HandleConfirmation(confirmations <-chan amqp.Confirmation) {
	log.Printf("Waiting for publish confirmation...")

	if confirmed := <-confirmations; confirmed.Ack {
		log.Printf("Confirmed delivery : %d", confirmed.DeliveryTag)
	} else {
		log.Printf("Failed delivery : %d", confirmed.DeliveryTag)
	}
}
