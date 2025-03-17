package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheRangiCrew/WITS/services/ingest/nwws/internal/nwws"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Run() {
	nwwsoi, err := nwws.New(&nwws.Config{
		Server:   os.Getenv("NWWSOI_SERVER") + ":5222",
		Room:     os.Getenv("NWWSOI_ROOM"),
		User:     os.Getenv("NWWSOI_USER"),
		Pass:     os.Getenv("NWWSOI_PASS"),
		Resource: os.Getenv("NWWSOI_RESOURCE"),
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := make(chan bool, 1)

	conn, err := amqp.Dial(os.Getenv("RABBIT"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"nwwsoi.exchange", // name
		"direct",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		slog.Error("failed to declare exchange: " + err.Error())
		return
	}

	_, err = ch.QueueDeclare(
		"nwwsoi.queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		slog.Error("failed to declare queue: " + err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("\033[32m *** RabbitMQ connected *** \033[m")

	messages := make(chan *nwws.Message)
	go nwwsoi.Start(messages)

	go func() {
		for message := range messages {
			body, err := json.Marshal(message)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			err = ch.PublishWithContext(ctx,
				"nwwsoi.exchange", // exchange
				"nwwsoi.awips",    // routing key
				false,             // mandatory
				false,             // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				})
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}
	}()

	go func() {
		<-sigs
		stop <- true
	}()

	<-stop
	slog.Info("Shutting down...")
	nwwsoi.Stop()
}
