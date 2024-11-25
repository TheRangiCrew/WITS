package handler

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/surrealdb/surrealdb.go"
)

type Data struct {
	UGC map[string]db.UGC
}

type Handler struct {
	DB      *surrealdb.DB
	Data    *Data
	ErrChan chan error
}

func New(config db.DBConfig) (*Handler, error) {
	db, err := db.New(config)
	if err != nil {
		return nil, err
	}

	handler := Handler{
		DB:      db,
		Data:    &Data{},
		ErrChan: make(chan error),
	}

	err = handler.LoadUGC()
	if err != nil {
		return nil, err
	}

	return &handler, nil
}

func (handler *Handler) LoadUGC() error {
	slog.Info("Getting UGC data")

	// Get the latest UGC data
	queryResult, err := surrealdb.Query[[]db.UGC](handler.DB, "SELECT * OMIT geometry, centre FROM ugc WHERE valid_to == null", map[string]interface{}{})
	if err != nil {
		return err
	}

	result := *queryResult

	if len(result[0].Result) == 0 {
		return fmt.Errorf("Received 0 UGC records")
	}

	data := map[string]db.UGC{}
	for _, ugc := range result[0].Result {
		data[ugc.ID.ID.(string)] = ugc
	}

	handler.Data.UGC = data

	slog.Info("Retrieved UGC data")

	return nil
}

func (handler *Handler) Handle(text string, receivedAt time.Time) error {
	product, err := awips.New(text)
	if err != nil {
		return err
	}

	if product.AWIPS.Product == "CAP" || product.AWIPS.Product == "WOU" {
		return nil
	}

	dbProduct, err := handler.TextProduct(product, receivedAt)
	if err != nil {
		return err
	}

	if product.HasVTEC() {
		err := handler.vtec(product, dbProduct)
		if err != nil {
			return err
		}
	}

	return err
}
