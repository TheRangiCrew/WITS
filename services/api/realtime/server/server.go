package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/twpayne/go-geos"
	pgxgeos "github.com/twpayne/pgx-geos"
)

const (
	RealtimeExchange = "realtime.exchange"
)

type Server struct {
	DB       *pgxpool.Pool
	RDB      *redis.Client
	Rabbit   *amqp.Connection
	Channels map[string]*RabbitChannel
	Router   *gin.Engine
	Hubs     map[string]*Hub
}

// Start the server
func (server *Server) Start() {
	port := os.Getenv("PORT")

	slog.Info(fmt.Sprintf("Starting server on %s...", port))

	// Start consuming messages
	server.Consume()

	// Start websocket hub
	for _, hub := range server.Hubs {
		go hub.run()
	}

	// Run API server
	server.Router.Run(":" + port)
}

// Start consuming message from RabbitMQ channels
func (server *Server) Consume() {
	if len(server.Channels) == 0 {
		slog.Info("No channels to consume")
		return
	}

	for name, feed := range server.Channels {
		err := feed.Start()
		if err != nil {
			slog.Error("Failed to start consuming messages from "+name, "error", err)
			continue
		}
	}

	return
}

func (server *Server) setupChannels() {
	channels := make(map[string]*RabbitChannel)
	channelConfigs := []struct {
		Name, Queue, Exchange, RoutingKey string
		handler                           func(msg amqp.Delivery) error
	}{
		{"warnings", "warnings.queue", RealtimeExchange, "realtime.warnings", server.handleWarning},
		// Add more feeds as needed
	}
	//
	for _, cfg := range channelConfigs {
		channel, err := newRabbitChannel(server.Rabbit, cfg.Name, cfg.Queue, cfg.Exchange, cfg.RoutingKey, cfg.handler)
		if err != nil {
			slog.Error("Failed to setup RabbitMQ feed", "feed", cfg.Name, "error", err)
			continue
		}
		channels[cfg.Name] = channel
	}
	server.Channels = channels
}

// Create a new server instance
func New() *Server {
	// Create Redis connection
	rdb := createRedisConnection()

	// Create database connection pool
	db, err := createDBPool()
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		return nil
	}

	// Create RabbitMQ connection
	rabbit, err := newRabbitConnection()
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return nil
	}

	// Create the server instance
	server := &Server{
		DB:     db,
		RDB:    rdb,
		Rabbit: rabbit,
		Router: gin.Default(),
		Hubs:   map[string]*Hub{},
	}

	// Create RabbitMQ channels
	server.setupChannels()

	err = rdb.FlushDB(context.Background()).Err()
	if err != nil {
		slog.Error("Failed to flush Redis database", "error", err)
		return nil
	}

	// Initialise MCD
	slog.Info("Initialising MCD data...")
	err = server.InitialiseMCD()
	if err != nil {
		slog.Error("Failed to initialise MCD data", "error", err)
		return nil
	}

	// Initialise warnings
	slog.Info("Initialising warning data...")
	err = server.InitialiseWarningData()
	if err != nil {
		slog.Error("Failed to initialise warning data", "error", err)
		return nil
	}

	return server
}

// Setup a database connection pool
func createDBPool() (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		if err := pgxgeos.Register(ctx, conn, geos.NewContext()); err != nil {
			return err
		}
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	slog.Info("\033[32m *** Database connected *** \033[m")

	return pool, err
}

func createRedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
