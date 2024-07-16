package services

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQPublisher(uri string) (*RabbitMQPublisher, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ { // Try to connect 10 times
		conn, err = amqp.Dial(uri)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 5 seconds... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQPublisher{conn: conn, channel: channel}, nil
}

func (p *RabbitMQPublisher) Publish(queueName string, message string) error {
	_, err := p.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *RabbitMQPublisher) Close() {
	if err := p.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := p.conn.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}
