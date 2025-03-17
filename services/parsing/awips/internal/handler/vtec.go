package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/handler/util"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/jackc/pgx/v5"
	"github.com/twpayne/go-geos"
)

// VTEC Event
type VTECEvent struct {
	ID           int       `json:"id,omitempty"`
	EventID      string    `json:"event_id"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	Issued       time.Time `json:"issued"`
	Starts       time.Time `json:"starts,omitempty"`
	Expires      time.Time `json:"expires"`
	Ends         time.Time `json:"ends,omitempty"`
	EndInitial   time.Time `json:"end_initial,omitempty"`
	Class        string    `json:"class"`
	Phenomena    string    `json:"phenomena"`
	WFO          string    `json:"wfo"`
	Significance string    `json:"significance"`
	EventNumber  int       `json:"event_number"`
	Year         int       `json:"year"`
	Title        string    `json:"title"`
	IsEmergency  bool      `json:"is_emergency"`
	IsPDS        bool      `json:"is_pds"`
	PolygonStart *geos.Geom
}

type VTECUGC struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	WFO          string    `json:"wfo"`
	Phenomena    string    `json:"phenomena"`
	Significance string    `json:"significance"`
	EventNumber  int       `json:"event_number"`
	UGC          int       `json:"ugc"`
	Issued       time.Time `json:"issued"`
	Starts       time.Time `json:"starts,omitempty"`
	Expires      time.Time `json:"expires"`
	Ends         time.Time `json:"ends,omitempty"`
	EndInitial   time.Time `json:"end_initial,omitempty"`
	Action       string    `json:"action"`
	Latest       int       `json:"latest"`
	Year         int       `json:"year"`
}

// VTEC Event Update
type VTECUpdate struct {
	ID            int        `json:"id"`
	CreatedAt     time.Time  `json:"created_at,omitempty"`
	Issued        time.Time  `json:"issued"`
	Starts        time.Time  `json:"starts,omitempty"`
	Expires       time.Time  `json:"expires"`
	Ends          time.Time  `json:"ends,omitempty"`
	Text          string     `json:"text"`
	Product       string     `json:"product"`
	WFO           string     `json:"wfo"`
	Action        string     `json:"action"`
	Class         string     `json:"class"`
	Phenomena     string     `json:"phenomena"`
	Significance  string     `json:"significance"`
	EventNumber   int        `json:"event_number"`
	Year          int        `json:"year"`
	Title         string     `json:"title"`
	IsEmergency   bool       `json:"is_emergency"`
	IsPDS         bool       `json:"is_pds"`
	Polygon       *geos.Geom `json:"polygon,omitempty"`
	Direction     *int       `json:"direction"`
	Location      *geos.Geom `json:"location"`
	Speed         *int       `json:"speed"`
	SpeedText     *string    `json:"speed_text"`
	TMLTime       *time.Time `json:"tml_time"`
	UGC           []string   `json:"ugc"`
	Tornado       string
	Damage        string
	HailThreat    string
	HailTag       string
	WindThreat    string
	WindTag       string
	FlashFlood    string
	RainfallTag   string
	FloodTagDam   string
	SpoutTag      string
	SnowSquall    string
	SnowSquallTag string
}

type VTEC struct {
	Class        string `json:"class"`
	Action       string `json:"action"`
	WFO          string `json:"wfo"`
	Phenomena    string `json:"phenomena"`
	Significance string `json:"significance"`
	EventNumber  int    `json:"event_number"`
	Start        string `json:"start"`
	End          string `json:"end"`
}

type vtecHandler struct {
	*Handler
	event   *VTECEvent
	product TextProduct
	segment awips.TextProductSegment
	vtec    awips.VTEC
	ugc     []struct {
		ID  int
		UGC string
	}
}

func (handler *Handler) vtec(product *awips.TextProduct, receivedAt time.Time) {

	textProduct, err := handler.TextProduct(product, receivedAt)
	if err != nil {
		handler.logger.Error(err.Error())
		return
	}

	handler.logger.SetProduct(textProduct.ProductID)

	// Process each segment separately since they reference different UGC areas
	for i, segment := range product.Segments {

		if len(segment.VTEC) == 0 {
			handler.logger.Info(fmt.Sprintf("Product %s segment %d does not have VTECs. Skipping...", textProduct.ProductID, i))
			continue
		}

		// Go through each VTEC in the segment and process it
		for _, vtec := range segment.VTEC {
			if vtec.Class == "T" || vtec.Action == "ROU" {
				continue
			}

			var year int
			if vtec.Start != nil {
				year = vtec.Start.Year()
			} else {
				year = product.Issued.Year()
			}

			var event *VTECEvent

			vh := vtecHandler{
				Handler: handler,
				product: *textProduct,
				segment: segment,
				event:   event,
				vtec:    vtec,
			}

			// Lets check if the VTEC Event is already in the database
			// Find the event
			rows, err := handler.db.Query(handler.db.CTX, `
			SELECT * FROM vtec_event WHERE
			wfo = $1 AND phenomena = $2 AND significance = $3 AND event_number = $4 AND year = $5
			`, vtec.WFO[1:], vtec.Phenomena, vtec.Significance, vtec.EventNumber, year)
			if err != nil {
				handler.logger.Error("failed to get vtec_event: " + err.Error())
				continue
			}

			if !rows.Next() {
				// The database needs a start and end time but VTECs may not have one.
				if vtec.Start == nil {
					// Use the issue time for the start time
					vtec.Start = &product.Issued
				}
				if vtec.End == nil {
					// Use the expiry of the product for the end time
					vtec.End = &segment.Expires
				}

				var polygon *geos.Geom
				if segment.LatLon != nil {
					p := util.PolygonFromAwips(*segment.LatLon.Polygon)
					polygon = &p
				}

				event = &VTECEvent{
					Issued:       product.Issued,
					Starts:       *vtec.Start,
					Expires:      segment.UGC.Expires,
					Ends:         *vtec.End,
					EndInitial:   *vtec.End,
					Class:        vtec.Class,
					Phenomena:    vtec.Phenomena,
					WFO:          vtec.WFO,
					Significance: vtec.Significance,
					EventNumber:  vtec.EventNumber,
					Year:         year,
					Title:        vtec.Title(segment.IsEmergency()),
					IsEmergency:  segment.IsEmergency(),
					IsPDS:        segment.IsPDS(),
					PolygonStart: polygon,
				}

				_, err := handler.db.Exec(context.Background(), `
				INSERT INTO vtec_event(issued, starts, expires, ends, end_initial, class, phenomena, wfo, 
				significance, event_number, year, title, is_emergency, is_pds, polygon_start) VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);
				`, event.Issued, event.Starts, event.Expires, event.Ends, event.EndInitial, event.Class,
					event.Phenomena, event.WFO, event.Significance, event.EventNumber, event.Year, event.Title,
					event.IsEmergency, event.IsPDS, event.PolygonStart)

				if err != nil {
					handler.logger.Error(fmt.Sprintf("failed to create vtec_event: %s", err.Error()))
					continue
				}

				vh.event = event

				vh.create()

			} else {
				event = &VTECEvent{}
				err = rows.Scan(
					&event.ID,
					&event.CreatedAt,
					&event.UpdatedAt,
					&event.Issued,
					&event.Starts,
					&event.Expires,
					&event.Ends,
					&event.EndInitial,
					&event.Class,
					&event.Phenomena,
					&event.WFO,
					&event.Significance,
					&event.EventNumber,
					&event.Year,
					&event.Title,
					&event.IsEmergency,
					&event.IsPDS,
					&event.PolygonStart)

				if err != nil {
					handler.logger.Error("failed to scan vtec_event: " + err.Error())
				}

				vh.event = event

				vh.update()
			}
		}
	}
}

func (handler *vtecHandler) create() {
	handler.updateTimes()

	err := handler.segmentUGC()
	if err != nil {
		handler.logger.Warn("failed to retrieve UGCs for VTEC segment: " + err.Error())
		return
	}

	err = handler.createUpdates()
	if err != nil {
		handler.logger.Warn(fmt.Sprintf("failed to create updates for VTEC event for %s: %s", handler.vtec.Original, err.Error()))
		return
	}

	err = handler.relateUGC()
	if err != nil {
		handler.logger.Warn(fmt.Sprintf("failed to create ugc relation for VTEC event for %s: %s", handler.vtec.Original, err.Error()))
		return
	}
}

func (handler *vtecHandler) update() {
	vtec := handler.vtec
	segment := handler.segment

	handler.updateTimes()

	err := handler.segmentUGC()
	if err != nil {
		handler.logger.Warn("failed to retrieve UGCs for VTEC segment: " + err.Error())
		return
	}

	err = handler.createUpdates()
	if err != nil {
		handler.logger.Warn(fmt.Sprintf("failed to create updates for VTEC event for %s: %s", handler.vtec.Original, err.Error()))
		return
	}

	if vtec.Action == "EXA" || vtec.Action == "EXB" {
		err := handler.relateUGC()
		if err != nil {
			handler.logger.Warn(fmt.Sprintf("failed to create ugc relation for VTEC event for %s: %s", handler.vtec.Original, err.Error()))
			return
		}
	} else {
		err := handler.updateUGC()
		if err != nil {
			handler.logger.Warn(fmt.Sprintf("failed to update ugc relation for VTEC event for %s: %s", handler.vtec.Original, err.Error()))
			return
		}
	}

	_, err = handler.db.Exec(handler.db.CTX, `
	UPDATE vtec_event SET updated_at = CURRENT_TIMESTAMP, is_emergency = $6, is_pds = $7 WHERE
			wfo = $1 AND phenomena = $2 AND significance = $3 AND event_number = $4 AND year = $5
			`, vtec.WFO, vtec.Phenomena, vtec.Significance, vtec.EventNumber, handler.event.Year, segment.IsEmergency(), segment.IsPDS())
	if err != nil {
		handler.logger.Warn(fmt.Sprintf("failed to update VTEC event %s: %s", handler.vtec.Original, err.Error()))
		return
	}
}

// Update the times of the event relative to the action
func (handler *vtecHandler) updateTimes() {
	event := handler.event
	product := handler.product
	segment := handler.segment
	vtec := handler.vtec

	// The product expires at the UGC expiry time
	var end time.Time
	if vtec.End == nil {
		end = segment.UGC.Expires
		handler.logger.Info("VTEC end time is nil. Defaulting to UGC expiry time.")
	} else {
		end = *vtec.End
	}

	switch vtec.Action {
	case "CAN":
		fallthrough
	case "UPG":
		event.Expires = segment.UGC.Expires
		event.Ends = product.Issued.UTC()
	case "EXP":
		event.Expires = end
		event.Ends = end
	case "EXT":
		fallthrough
	case "EXB":
		event.Ends = end
		event.Expires = segment.UGC.Expires
	default:
		// NEW and CON
		if event.Ends.Before(end) {
			event.Ends = end
		}
		if event.Expires.Before(segment.Expires) {
			event.Expires = segment.Expires
		}
	}
}

func (handler *vtecHandler) createUpdates() error {
	event := handler.event
	product := handler.product
	segment := handler.segment
	vtec := handler.vtec

	ugcs := []string{}
	for _, ugc := range handler.ugc {
		ugcs = append(ugcs, ugc.UGC)
	}

	var polygon *geos.Geom
	if segment.LatLon != nil {
		p := util.PolygonFromAwips(*segment.LatLon.Polygon)
		polygon = &p
	}

	var direction *int
	var location *geos.Geom
	var speed *int
	var speedText *string
	var tmlTime *time.Time
	if segment.TML != nil {
		direction = &segment.TML.Direction
		point := geos.NewPoint(segment.TML.Location[:])
		location = point
		speed = &segment.TML.Speed
		speedText = &segment.TML.SpeedString
		tmlTime = &segment.TML.Time
	}

	rows := []VTECUpdate{
		{
			Issued:        *product.Issued,
			Starts:        event.Starts,
			Expires:       segment.UGC.Expires,
			Ends:          event.Ends,
			Text:          segment.Text,
			Product:       product.ProductID,
			WFO:           vtec.WFO,
			Action:        vtec.Action,
			Class:         vtec.Class,
			Phenomena:     vtec.Phenomena,
			Significance:  vtec.Significance,
			EventNumber:   vtec.EventNumber,
			Year:          event.Year,
			Title:         vtec.Title(segment.IsEmergency()),
			IsEmergency:   segment.IsEmergency(),
			IsPDS:         segment.IsPDS(),
			Polygon:       polygon,
			Direction:     direction,
			Location:      location,
			Speed:         speed,
			SpeedText:     speedText,
			TMLTime:       tmlTime,
			UGC:           ugcs,
			Tornado:       segment.Tags["tornado"],
			Damage:        segment.Tags["damage"],
			HailThreat:    segment.Tags["hailThreat"],
			HailTag:       segment.Tags["hail"],
			WindThreat:    segment.Tags["windThreat"],
			WindTag:       segment.Tags["wind"],
			FlashFlood:    segment.Tags["flashFlood"],
			RainfallTag:   segment.Tags["expectedRainfall"],
			FloodTagDam:   segment.Tags["damFailure"],
			SpoutTag:      segment.Tags["spout"],
			SnowSquall:    segment.Tags["snowSquall"],
			SnowSquallTag: segment.Tags["snowSquallImpact"],
		},
	}

	_, err := handler.db.CopyFrom(
		context.Background(),
		pgx.Identifier{"vtec_update"},
		[]string{"issued", "starts", "expires", "ends", "text", "product", "wfo", "action", "class", "phenomena", "significance", "event_number", "year", "title", "is_emergency", "is_pds", "polygon", "direction", "location", "speed", "speed_text", "tml_time", "ugc", "tornado", "damage", "hail_threat", "hail_tag", "wind_threat", "wind_tag", "flash_flood", "rainfall_tag", "flood_tag_dam", "spout_tag", "snow_squall", "snow_squall_tag"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			return []any{
				rows[i].Issued, rows[i].Starts, rows[i].Expires, rows[i].Ends, rows[i].Text, rows[i].Product, rows[i].WFO, rows[i].Action, rows[i].Class, rows[i].Phenomena, rows[i].Significance, rows[i].EventNumber, rows[i].Year, rows[i].Title, rows[i].IsEmergency, rows[i].IsPDS, rows[i].Polygon, rows[i].Direction, rows[i].Location, rows[i].Speed, rows[i].SpeedText, rows[i].TMLTime, rows[i].UGC, rows[i].Tornado, rows[i].Damage, rows[i].HailThreat, rows[i].HailTag, rows[i].WindThreat, rows[i].WindTag, rows[i].FlashFlood, rows[i].RainfallTag, rows[i].FloodTagDam, rows[i].SpoutTag, rows[i].SnowSquall, rows[i].SnowSquallTag,
			}, nil
		}),
	)
	if err != nil {
		return err
	}

	handler.warning(rows[0])

	return nil
}

func (handler *vtecHandler) relateUGC() error {
	event := handler.event
	product := handler.product
	vtec := handler.vtec
	segment := handler.segment

	start := event.Starts
	if start.Equal(*product.Issued) {
		start = *product.Issued
	}

	// The product expires at the UGC expiry time
	expires := segment.UGC.Expires
	var end time.Time
	if vtec.End == nil {
		end = expires
		handler.logger.Info("VTEC end time is nil. Defaulting to UGC expiry time.")
	} else {
		end = *vtec.End
	}

	rows := []VTECUGC{}
	for _, ugc := range handler.ugc {
		rows = append(rows, VTECUGC{
			WFO:          event.WFO,
			Phenomena:    event.Phenomena,
			Significance: event.Significance,
			EventNumber:  event.EventNumber,
			UGC:          ugc.ID,
			Issued:       event.Issued,
			Starts:       start,
			Expires:      expires,
			Ends:         end,
			EndInitial:   end,
			Action:       vtec.Action,
			Year:         event.Year,
		})
	}

	_, err := handler.db.CopyFrom(
		handler.db.CTX,
		pgx.Identifier{"vtec_ugc"},
		[]string{"wfo", "phenomena", "significance", "event_number", "ugc", "issued", "starts", "expires", "ends",
			"end_initial", "action", "year"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			return []any{rows[i].WFO, rows[i].Phenomena, rows[i].Significance, rows[i].EventNumber, rows[i].UGC,
				rows[i].Issued, rows[i].Starts, rows[i].Expires, rows[i].Ends, rows[i].EndInitial, rows[i].Action,
				rows[i].Year}, nil
		}))
	if err != nil {
		return err
	}

	return nil
}

func (handler *vtecHandler) updateUGC() error {
	event := handler.event
	vtec := handler.vtec
	segment := handler.segment

	expires := segment.UGC.Expires
	end := event.Ends

	ugcs := []int{}
	for _, ugc := range handler.ugc {
		ugcs = append(ugcs, ugc.ID)
	}

	_, err := handler.db.Exec(handler.db.CTX, `
	UPDATE vtec_ugc SET expires = $1, ends = $2, action = $3 WHERE
	wfo = $4 AND phenomena = $5 AND significance = $6 AND event_number = $7 AND year = $8
	AND ugc = ANY($9)
	`, expires, end, vtec.Action, event.WFO, event.Phenomena, event.Significance, event.EventNumber,
		event.Year, ugcs)
	return err
}

// Creates an array of UGC record IDs for the database
func (handler *vtecHandler) segmentUGC() error {
	segment := handler.segment

	ids := []string{}
	// For each state...
	for _, state := range segment.UGC.States {
		// ...and for each area...
		for _, area := range state.Areas {
			ugcType := state.Type
			isFire := false
			// Fire weather (FW) events have different zones
			if handler.vtec.Phenomena == "FW" {
				ugcType = "F"
				isFire = true
			}
			if area == "000" || area == "ALL" {
				// Find all UGC codes from the state
				rows, err := handler.db.Query(handler.db.CTX, `
				SELECT ugc FROM ugc WHERE state = $1 AND type = $2 AND is_fire = $3 AND valid_to IS NULL;
				`, state.ID, state.Type, isFire)
				if err != nil {
					return fmt.Errorf("error retrieving ALL ugc: %s", err.Error())
				}

				for rows.Next() {
					var id string
					err := rows.Scan(&id)
					if err != nil {
						return err
					}
					ids = append(ids, id)
				}

				if len(ids) == 0 {
					return fmt.Errorf("got 0 UGC for %s. Expected ALL", state.ID+ugcType+area)
				}
			} else {
				// Get the needed UGCs
				id := state.ID + ugcType + area
				ids = append(ids, id)
			}
		}
	}

	rows, err := handler.db.Query(handler.db.CTX, `
	SELECT id, ugc FROM ugc WHERE ugc = ANY($1) AND valid_to IS NULL
	`, ids)
	if err != nil {
		return fmt.Errorf("error retrieving listed ugcs: %s", err.Error())
	}

	ugcs := []struct {
		ID  int
		UGC string
	}{}
	for rows.Next() {
		ugc := struct {
			ID  int
			UGC string
		}{}
		err := rows.Scan(&ugc.ID, &ugc.UGC)
		if err != nil {
			return fmt.Errorf("error scanning ugcs: %s", err.Error())
		}
		ugcs = append(ugcs, ugc)
	}

	handler.ugc = ugcs

	return nil
}
