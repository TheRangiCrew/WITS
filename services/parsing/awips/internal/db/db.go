package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twpayne/go-geos"
	pgxgeos "github.com/twpayne/pgx-geos"
)

func New() (*pgxpool.Pool, error) {
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
