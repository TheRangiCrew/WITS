package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RealtimeExchange = "realtime.exchange"
)

type ServerConfig struct {
	MinLog int
}

type Server struct {
	DB         *pgxpool.Pool
	Rabbit     *amqp.Connection
	Publishers map[string]*amqp.Channel
	config     ServerConfig
}

type Message struct {
	Text       string
	ReceivedAt time.Time
}

func New(config ServerConfig) (*Server, error) {
	// Create a new database connection pool
	db, err := db.New()
	if err != nil {
		return nil, err
	}

	rabbit, err := newRabbitConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	server := Server{
		DB:     db,
		Rabbit: rabbit,
		config: config,
	}

	return &server, nil
}

func (server *Server) Start() {

	// Set up interrupt signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := make(chan bool, 1)

	// Perform health check, not sure if this does anything useful
	err := server.HeathCheck()
	if err != nil {
		slog.Error(err.Error())
		slog.Info("Stopping...")
		return
	}

	// Initialis
	err = server.initialiseRabbit()
	if err != nil {
		slog.Error(err.Error())
		slog.Info("Stopping...")
		return
	}

	channel, err := server.Rabbit.Channel()
	if err != nil {
		slog.Error("failed to create RabbitMQ channel: " + err.Error())
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"nwwsoi.queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		slog.Error("failed to register queue: " + err.Error())
		return
	}

	err = channel.QueueBind(
		queue.Name,        // queue name
		"nwwsoi.awips",    // routing key
		"nwwsoi.exchange", // exchange
		false,
		nil)
	if err != nil {
		slog.Error("failed to register queue binding: " + err.Error())
		return
	}

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		slog.Error("failed to register consumer: " + err.Error())
		return
	}

	slog.Info("\033[32m *** RabbitMQ connected *** \033[m")

	go func() {
		for message := range msgs {
			m := &Message{}
			err := json.Unmarshal(message.Body, m)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			go server.HandleMessage(m.Text, m.ReceivedAt.UTC())
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				err := server.HeathCheck()
				if err != nil {
					slog.Error(err.Error())
					stop <- true
					return
				}
			}
		}
	}()

	go func() {
		<-sigs
		stop <- true
	}()

	<-stop
	slog.Info("Shutting down...")
	ticker.Stop()
	defer server.Rabbit.Close()
}

// Makes a Rabbit
func (server *Server) initialiseRabbit() error {

	conn, err := amqp.Dial(os.Getenv("RABBIT"))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	server.Rabbit = conn

	server.Publishers = make(map[string]*amqp.Channel)
	publisherConfigs := []struct {
		Name, Type string
	}{
		{RealtimeExchange, "topic"},
	}

	for _, cfg := range publisherConfigs {
		ch, err := newPublishingChannel(server.Rabbit, cfg.Name, cfg.Type)
		if err != nil {
			slog.Error("Failed to create RabbitMQ publisher", "name", cfg.Name, "error", err)
			continue
		}
		server.Publishers[cfg.Name] = ch
	}

	return nil
}

func (server *Server) GetPublisher(name string) (*amqp.Channel, error) {
	publisher, ok := server.Publishers[name]
	if !ok {
		return nil, fmt.Errorf("publisher %s not found", name)
	}
	return publisher, nil
}

func (server *Server) HeathCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.DB.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
