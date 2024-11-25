package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func (handler *Handler) TextProduct(product *awips.TextProduct, receivedAt time.Time) (*db.Product, error) {

	issued := models.CustomDateTime{Time: product.Issued}

	data := db.Product{
		Product:    product.Product,
		Issuer:     &models.RecordID{Table: "office", ID: product.Office},
		Text:       strings.ReplaceAll(product.Text, `"`, `\"`),
		ReceivedAt: &models.CustomDateTime{Time: receivedAt},
		WMO: db.WMO{
			Datatype: product.WMO.Datatype,
			Issued:   models.CustomDateTime{Time: product.Issued},
			Original: product.WMO.Original,
			WFO:      product.WMO.Office,
			BBB:      product.WMO.BBB,
		},
	}

	query := fmt.Sprintf(`CREATE ONLY product CONTENT {
	id: product:[
		"%s",
		%s,
		(SELECT count() FROM product:["%s", %s, 0]..["%s", %s, 9999] GROUP ALL)[0].count OR 0
	],
	issuer: %s,
	product: "%s",
	received_at: %s,
	text: "%s",
	wmo: {
		bbb: "%s",
		datatype: "%s",
		issued: %s,
		original: "%s",
		wfo: "%s"
	}
}`, product.AWIPS.Original, issued.SurrealString(), product.AWIPS.Original, issued.SurrealString(), product.AWIPS.Original, issued.SurrealString(), data.Issuer.String(), data.Product, data.ReceivedAt.SurrealString(), data.Text, data.WMO.BBB, data.WMO.Datatype, data.WMO.Issued.SurrealString(), data.WMO.Original, data.WMO.WFO)

	res, err := surrealdb.Query[db.Product](handler.DB, query, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	result := (*res)[0].Result

	return &result, nil
}
