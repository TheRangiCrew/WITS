package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/logger"
	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/util"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/TheRangiCrew/go-nws/pkg/awips/products"
)

func (server *Server) mcd(product *awips.TextProduct, receivedAt time.Time) error {

	log := logger.New(server.DB, slog.Level(server.config.MinLog))

	textProduct, err := server.TextProduct(product, receivedAt)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.With("product", textProduct.ProductID)

	for _, segment := range product.Segments {
		mcdProduct, err := products.ParseMCD(segment.Text)
		if err != nil {
			return err
		}

		polygon := util.PolygonFromAwips(mcdProduct.Polygon)

		issued := time.Date(product.Issued.Year(), product.Issued.Month(), mcdProduct.Issued.Day(), mcdProduct.Issued.Hour(), mcdProduct.Issued.Minute(), 0, 0, time.UTC)
		expires := time.Date(segment.Expires.Year(), segment.Expires.Month(), mcdProduct.Expires.Day(), mcdProduct.Expires.Hour(), mcdProduct.Expires.Minute(), 0, 0, time.UTC)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = server.DB.Exec(ctx, `
		INSERT INTO mcd.mcd (id, product, issued, expires, year, concerning, geom, watch_probability, most_prob_tornado, most_prob_gust, most_prob_hail) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
		`, mcdProduct.Number, textProduct.ProductID, issued, expires, textProduct.Issued.Year(), mcdProduct.Concerning, &polygon, mcdProduct.WatchProbability, mcdProduct.MostProbTornado, mcdProduct.MostProbGust, mcdProduct.MostProbHail)
		if err != nil {
			return err
		}

	}

	return nil
}
