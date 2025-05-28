package server

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func newRabbitConnection() (*amqp.Connection, error) {

	conn, err := amqp.Dial(os.Getenv("RABBIT"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newPublishingChannel(conn *amqp.Connection, name, exchangeType string) (*amqp.Channel, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		name,         // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	return ch, nil

}
