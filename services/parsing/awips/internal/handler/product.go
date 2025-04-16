package handler

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/TheRangiCrew/go-nws/pkg/awips"
)

// Text Product
type TextProduct struct {
	ID         int        `json:"id"`
	ProductID  string     `json:"product_id"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	ReceivedAt *time.Time `json:"received_at"`
	Issued     *time.Time `json:"issued"`
	Source     string     `json:"source"`
	Data       string     `json:"data"`
	WMO        string     `json:"wmo"`
	AWIPS      string     `json:"awips"`
	BBB        string     `json:"bbb"`
}

func (handler *Handler) TextProduct(product *awips.TextProduct, receivedAt time.Time) (*TextProduct, error) {

	id := fmt.Sprintf("%s-%s-%s-%s", product.Issued.UTC().Format("200601021504"), product.Office, product.WMO.Datatype, product.AWIPS.Original)

	if len(product.WMO.BBB) > 0 {
		id += "-" + product.WMO.BBB
	}

	textProduct := TextProduct{
		ProductID:  id,
		ReceivedAt: &receivedAt,
		Issued:     &product.Issued,
		Source:     product.AWIPS.WFO,
		Data:       product.Text,
		WMO:        product.WMO.Datatype,
		AWIPS:      product.AWIPS.Original,
		BBB:        product.WMO.BBB,
	}

	rows, err := handler.db.Query(context.Background(), `
	INSERT INTO awips.products (product_id, received_at, issued, source, data, wmo, awips, bbb) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at;
	`, id, receivedAt, product.Issued, product.AWIPS.WFO, product.Text, product.WMO.Datatype, product.AWIPS.Original, product.WMO.BBB)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&textProduct.ID, &textProduct.CreatedAt)
		return &textProduct, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return nil, errors.New("no rows returned when creating new text product: " + rows.Err().Error())
}

func (product *TextProduct) isCorrection() bool {
	resent := regexp.MustCompile("...(RESENT|RETRANSMITTED|CORRECTED)")

	if len(resent.FindString(product.Data)) > 0 {
		return true
	}
	if len(product.BBB) > 0 && (string(product.BBB[0]) == "A" || string(product.BBB[0]) == "C") {
		return true
	}

	return false
}
