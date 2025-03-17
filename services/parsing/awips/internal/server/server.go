package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/handler"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ServerConfig struct {
	MinLog int
}

type Server struct {
	DB     *db.Pool
	Rabbit *amqp.Channel
	MinLog int
}

type Message struct {
	Text       string
	ReceivedAt time.Time
}

func New(config ServerConfig) (*Server, error) {
	db, err := db.New()
	if err != nil {
		return nil, err
	}

	server := Server{
		DB:     db,
		MinLog: config.MinLog,
	}

	return &server, nil
}

func (server *Server) Start() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := make(chan bool, 1)

	err := server.HeathCheck()
	if err != nil {
		slog.Error(err.Error())
		slog.Info("Stopping...")
		return
	}

	err = server.InitialiseRabbit()
	if err != nil {
		slog.Error(err.Error())
		slog.Info("Stopping...")
		return
	}

	q, err := server.Rabbit.QueueDeclare(
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

	err = server.Rabbit.QueueBind(
		q.Name,            // queue name
		"nwwsoi.awips",    // routing key
		"nwwsoi.exchange", // exchange
		false,
		nil)
	if err != nil {
		slog.Error("failed to register queue binding: " + err.Error())
		return
	}

	msgs, err := server.Rabbit.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
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

			h, err := handler.New(server.DB, server.MinLog)
			h.AddRabbit(server.Rabbit)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			go h.Handle(m.Text, m.ReceivedAt.UTC())
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
func (server *Server) InitialiseRabbit() error {

	conn, err := amqp.Dial(os.Getenv("RABBIT"))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: " + err.Error())
	}

	server.Rabbit, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: " + err.Error())
	}

	err = server.Rabbit.ExchangeDeclare(
		"nwwsoi.exchange", // name
		"direct",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: " + err.Error())
	}

	err = server.Rabbit.ExchangeDeclare(
		"awips.exchange", // name
		"direct",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: " + err.Error())
	}

	return nil
}

func (server *Server) HeathCheck() error {
	err := server.DB.Ping(server.DB.CTX)
	if err != nil {
		return err
	}
	return nil
}
