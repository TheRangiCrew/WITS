package handler

import (
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/handler/util"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/TheRangiCrew/go-nws/pkg/awips/products"
)

func (handler *Handler) mcd(product *awips.TextProduct, receivedAt time.Time) error {

	textProduct, err := handler.TextProduct(product, receivedAt)
	if err != nil {
		handler.logger.Error(err.Error())
		return err
	}

	handler.logger.SetProduct(textProduct.ProductID)

	for _, segment := range product.Segments {
		mcdProduct, err := products.ParseMCD(segment.Text)
		if err != nil {
			return err
		}

		polygon := util.PolygonFromAwips(mcdProduct.Polygon)

		issued := time.Date(product.Issued.Year(), product.Issued.Month(), mcdProduct.Issued.Day(), mcdProduct.Issued.Hour(), mcdProduct.Issued.Minute(), 0, 0, time.UTC)
		expires := time.Date(segment.Expires.Year(), segment.Expires.Month(), mcdProduct.Expires.Day(), mcdProduct.Expires.Hour(), mcdProduct.Expires.Minute(), 0, 0, time.UTC)

		_, err = handler.db.Exec(handler.db.CTX, `
		INSERT INTO mcd (id, product, issued, expires, year, concerning, geom, watch_probability) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);
		`, mcdProduct.Number, textProduct.ProductID, issued, expires, textProduct.Issued.Year(), mcdProduct.Concerning, &polygon, mcdProduct.WatchProbability)
		if err != nil {
			return err
		}

	}

	return nil
}
