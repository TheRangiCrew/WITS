package handler

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/logger"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	logger *logger.Logger
	db     *pgxpool.Pool
	rabbit *amqp.Channel
	ctx    context.Context
}

func New(db *pgxpool.Pool, minLog int) (*Handler, error) {

	l := logger.New(db, slog.Level(minLog))

	handler := Handler{
		logger: &l,
		db:     db,
		ctx:    context.Background(),
	}

	return &handler, nil
}

// This is very similar to the go-nws AWIPS product parser. Duplicated since we need a slightly different workflow
func (handler *Handler) Handle(text string, receivedAt time.Time) error {
	// Get the WMO header
	wmo, err := awips.ParseWMO(text)
	if err != nil {
		handler.logger.Error(err.Error())
		return err
	}

	handler.logger.SetWMO(wmo.Original)

	// TODO: Handle test communications, we can probably utilise them (WOUS99)
	if wmo.Datatype == "WOUS99" || wmo.Datatype == "NTXX98" {
		handler.logger.Debug("Communication test message received. Ignoring")
		return nil
	}

	// Get the AWIPS header
	awipsHeader, err := awips.ParseAWIPS(text)
	if err != nil {
		handler.logger.Debug(err.Error())
	}

	ignore := []string{"CAP", "HML"}
	if slices.Contains(ignore, awipsHeader.Product) {
		handler.logger.Info(fmt.Sprintf("%s product is flagged. Ignoring", awipsHeader.Product))
		return nil
	}

	if awipsHeader.Original == "" {
		handler.logger.Info("AWIPS header not found. Product will not be stored.")
		return nil
	} else {
		handler.logger.SetAWIPS(awipsHeader.Original)
	}

	// Find the issue time
	issued, err := awips.GetIssuedTime(text)
	if err != nil {
		handler.logger.Error(err.Error())
		return err
	}
	if issued.IsZero() {
		handler.logger.Info("Product does not contain issue date. Defaulting to now (UTC)")
		issued = time.Now().UTC()
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
			handler.logger.Error(err.Error())
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
				handler.logger.Error(er.Error())
			}
			continue
		}

		latlon, err := awips.ParseLatLon(segment)
		if err != nil {
			handler.logger.Error(err.Error())
			continue
		}

		tags, e := awips.ParseTags(segment)
		if len(e) != 0 {
			for _, er := range e {
				handler.logger.Error(er.Error())
			}
		}

		tml, err := awips.ParseTML(segment, issued)
		if err != nil {
			handler.logger.Warn("failed to parse TML: " + err.Error())
		}

		segments = append(segments, awips.TextProductSegment{
			Text:    segment,
			VTEC:    vtec,
			UGC:     ugc,
			Expires: expires,
			LatLon:  latlon,
			Tags:    tags,
			TML:     tml,
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

	if product.AWIPS.Product == "WOU" {
		return nil
	}

	if product.AWIPS.Original == "SWOMCD" {
		err := handler.mcd(product, receivedAt)
		if err != nil {
			handler.logger.Error("error processing mcd: " + err.Error())
			return err
		}
	}

	if product.HasVTEC() {
		handler.vtec(product, receivedAt)
	}

	err = handler.logger.Save()
	if err != nil {
		slog.Error("error saving logs to database", "error", err.Error())
	}
	return err

}

func (handler *Handler) AddRabbit(rabbit *amqp.Channel) {
	handler.rabbit = rabbit
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
