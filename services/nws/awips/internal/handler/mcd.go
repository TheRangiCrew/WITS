package handler

import (
	"fmt"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/handler/util"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/TheRangiCrew/go-nws/pkg/awips/products"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func (handler *Handler) mcd(product *awips.TextProduct, dbProduct *db.Product) error {

	for _, segment := range product.Segments {
		mcdProduct, err := products.ParseMCD(product.Text)
		if err != nil {
			return err
		}

		id := models.NewRecordID("mcd", db.MCDID{Number: mcdProduct.Number, Year: product.Issued.Year()})

		issued := models.CustomDateTime{Time: time.Date(product.Issued.Year(), product.Issued.Month(), mcdProduct.Issued.Day(), mcdProduct.Issued.Hour(), mcdProduct.Issued.Minute(), 0, 0, time.UTC)}
		expires := models.CustomDateTime{Time: time.Date(segment.Expires.Year(), segment.Expires.Month(), mcdProduct.Expires.Day(), mcdProduct.Expires.Hour(), mcdProduct.Expires.Minute(), 0, 0, time.UTC)}

		polygon := util.PolygonFromAwips(mcdProduct.Polygon)

		_, err = surrealdb.Query[[]db.MCD](handler.DB, `UPSERT $id MERGE {
			original: $original,
			issued: $issued,
			expires: $expires,
			concerning: $concerning,
			polygon: $polygon,
			watch_probability: $probability
		}`,
			map[string]interface{}{
				"id":          id,
				"original":    dbProduct.ID,
				"issued":      issued,
				"expires":     expires,
				"concerning":  mcdProduct.Concerning,
				"polygon":     polygon,
				"probability": mcdProduct.WatchProbability,
			})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}
