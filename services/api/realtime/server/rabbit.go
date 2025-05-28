package server

import (
	"log/slog"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitChannel struct {
	Name       string
	Queue      string
	Exchange   string
	RoutingKey string
	Channel    *amqp.Channel
	Messages   <-chan amqp.Delivery
	Handler    func(msg amqp.Delivery) error
}

func (feed *RabbitChannel) Start() error {
	msgs, err := feed.Channel.Consume(
		feed.Queue, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	feed.Messages = msgs

	go func(feed *RabbitChannel) {
		for msg := range feed.Messages {
			err := feed.Handler(msg)
			if err != nil {
				slog.Error("Error handling message", "feed", feed.Name, "error", err)
				continue
			}
		}
	}(feed)

	return nil
}

func newRabbitConnection() (*amqp.Connection, error) {

	conn, err := amqp.Dial(os.Getenv("RABBIT"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newRabbitChannel(conn *amqp.Connection, name, queue, exchange, routingKey string, handler func(msg amqp.Delivery) error) (*RabbitChannel, error) {
	// New channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitChannel{
		Name:       name,
		Queue:      queue,
		Exchange:   exchange,
		RoutingKey: routingKey,
		Channel:    ch,
		Handler:    handler,
	}, nil
}
