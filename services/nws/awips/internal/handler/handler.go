package handler

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/logger"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/surrealdb/surrealdb.go"
)

type Handler struct {
	Logger  *logger.Logger
	DB      *surrealdb.DB
	UGCData map[string]db.UGC
	ctx     context.Context
}

func New(db *surrealdb.DB, ugcData map[string]db.UGC, minLog int) (*Handler, error) {

	l := logger.New(db, slog.Level(minLog))

	handler := Handler{
		Logger:  &l,
		DB:      db,
		UGCData: ugcData,
		ctx:     context.Background(),
	}

	return &handler, nil
}

// This is very similar to the go-nws AWIPS product parser. Duplicated since we need a slightly different workflow
func (handler *Handler) Handle(text string, receivedAt time.Time) error {
	// Get the WMO header
	wmo, err := awips.ParseWMO(text)
	if err != nil {
		return err
	}

	handler.Logger.SetWMO(wmo.Original)

	// TODO: Handle test communications, we can probably utilise them (WOUS99)
	if wmo.Datatype == "WOUS99" || wmo.Datatype == "NTXX98" {
		handler.Logger.Debug("Communication test message received. Ignoring")
		return nil
	}

	// Find the issue time
	issued, err := awips.GetIssuedTime(text)
	if err != nil {
		return err
	}
	if issued.IsZero() {
		handler.Logger.Info("Product does not contain issue date. Defaulting to now (UTC)")
		issued = time.Now().UTC()
	}

	// Get the AWIPS header
	awipsHeader, _ := awips.ParseAWIPS(text)

	if awipsHeader.Original == "" {
		handler.Logger.Info("AWIPS header not found. Product will not be stored.")
		return nil
	} else {
		handler.Logger.SetAWIPS(awipsHeader.Original)
	}

	// Segment the product
	splits := strings.Split(text, "$$")

	segments := []awips.TextProductSegment{}

	for _, segment := range splits {
		segment = strings.TrimSpace(segment)

		// Assume the segment is the end of the product if it is shorter than 10 characters
		if len(segment) < 20 {
			continue
		}

		ugc, err := awips.ParseUGC(segment)
		if err != nil {
			handler.Logger.Error(err.Error())
			continue
		}
		expires := time.Now().UTC()
		if ugc != nil {
			expires = time.Date(issued.Year(), issued.Month(), ugc.Expires.Day(), ugc.Expires.Hour(), ugc.Expires.Minute(), 0, 0, time.UTC)
			if ugc.Expires.Day() > wmo.Issued.Day() && ugc.Expires.Day() == 1 {
				expires = expires.AddDate(0, 1, 0)
			}
			ugc.Merge(issued)
		}

		// Find any VTECs that the segment may have
		vtec, e := awips.ParseVTEC(segment)
		if len(e) != 0 {
			for _, er := range e {
				handler.Logger.Error(er.Error())
			}
			continue
		}

		latlon, err := awips.ParseLatLon(text)
		if err != nil {
			handler.Logger.Error(err.Error())
			continue
		}

		tags, e := awips.ParseTags(text)
		if len(e) != 0 {
			for _, er := range e {
				handler.Logger.Error(er.Error())
			}
		}

		segments = append(segments, awips.TextProductSegment{
			Text:    segment,
			VTEC:    vtec,
			UGC:     ugc,
			Expires: expires,
			LatLon:  latlon,
			Tags:    tags,
		})

	}

	product := &awips.TextProduct{
		Text:     text,
		WMO:      wmo,
		AWIPS:    awipsHeader,
		Issued:   issued,
		Office:   wmo.Office,
		Product:  awipsHeader.Product,
		Segments: segments,
	}

	dbProduct, err := handler.TextProduct(product, receivedAt)
	if err != nil {
		return err
	}

	handler.Logger.SetProduct(*dbProduct.ID)

	if product.AWIPS.Product == "CAP" || product.AWIPS.Product == "WOU" {
		return nil
	}

	if product.AWIPS.Original == "SWOMCD" {
		err := handler.mcd(product, dbProduct)
		if err != nil {
			return err
		}
	}

	if product.HasVTEC() {
		handler.vtec(product, dbProduct.ID)
	}

	err = handler.Logger.Save()
	if err != nil {
		slog.Error("error saving logs to database", "error", err.Error())
	}
	return err
}

// func appendContext(parent context.Context, attr slog.Attr) context.Context {
// 	if parent == nil {
// 		parent = context.Background()
// 	}

// 	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
// 		v = append(v, attr)
// 		return context.WithValue(parent, slogFields, v)
// 	} else {
// 		v := []slog.Attr{}
// 		v = append(v, attr)
// 		return context.WithValue(parent, slogFields, v)
// 	}
// }
