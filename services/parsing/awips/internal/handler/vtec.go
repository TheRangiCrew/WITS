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

// VTEC UGC Relation
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

// Raw VTEC
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
	*Handler                          // The parent handler
	event    *VTECEvent               // The event
	product  TextProduct              // The event's product
	segment  awips.TextProductSegment // The product segment
	vtec     awips.VTEC               // The raw VTEC
	ugc      []struct {
		ID  int
		UGC string
	}
}

// Parse and upload or update a VTEC product
func (handler *Handler) vtec(product *awips.TextProduct, receivedAt time.Time) {
	// Parse the text product
	textProduct, err := handler.TextProduct(product, receivedAt)
	if err != nil {
		handler.logger.Error(err.Error())
		return
	}

	// Set the text product's ID in the logger
	handler.logger.SetProduct(textProduct.ProductID)

	// Process each segment separately since they reference different UGC areas
	for i, segment := range product.Segments {

		// This segment does not have a VTEC so we can skip it
		if len(segment.VTEC) == 0 {
			handler.logger.Info(fmt.Sprintf("Product %s segment %d does not have VTECs. Skipping...", textProduct.ProductID, i))
			continue
		}

		// Go through each VTEC in the segment and process it
		for _, vtec := range segment.VTEC {
			// Skip test and routine products to save on resources
			if vtec.Class == "T" || vtec.Action == "ROU" {
				continue
			}

			// 	TODO: Rethink this because it is not going to work during the new year.
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Lets check if the VTEC Event is already in the database
			rows, err := handler.db.Query(ctx, `
			SELECT * FROM vtec.events WHERE
			wfo = $1 AND phenomena = $2 AND significance = $3 AND event_number = $4 AND year = $5
			`, vtec.WFO, vtec.Phenomena, vtec.Significance, vtec.EventNumber, year)
			if err != nil {
				handler.logger.Error("failed to get vtec_event: " + err.Error())
				continue
			}
			defer rows.Close()

			// If the event is not already there then create the event for the first time
			if !rows.Next() {
				if rows.Err() != nil {
					handler.logger.Error("failed to get vtec_event: " + rows.Err().Error())
					continue
				}

				// The database needs a start and end time but VTECs may not have one.
				if vtec.Start == nil {
					// Use the issue time for the start time
					vtec.Start = &product.Issued
				}
				if vtec.End == nil {
					// Use the expiry of the product for the end time
					vtec.End = &segment.Expires
				}

				// Create the polygon if there is one.
				var polygon *geos.Geom
				if segment.LatLon != nil {
					p := util.PolygonFromAwips(*segment.LatLon.Polygon)
					polygon = p
				}

				// Build the event
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

				// Insert the event
				_, err := handler.db.Exec(context.Background(), `
				INSERT INTO vtec.events(issued, starts, expires, ends, end_initial, class, phenomena, wfo, 
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

				err = vh.handle()
				if err != nil {
					handler.logger.Error(err.Error())
					return
				}
				vh.update()
			} else {
				// Get the event from the row
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

				err = vh.handle()
				if err != nil {
					handler.logger.Error(err.Error())
					return
				}
				vh.update()
			}
		}
	}
}

func (handler *vtecHandler) handle() error {
	handler.updateTimes()

	err := handler.segmentUGC()
	if err != nil {
		handler.logger.Warn("failed to retrieve UGCs for VTEC segment: " + err.Error())
		return err
	}

	err = handler.createUpdates()
	if err != nil {
		handler.logger.Warn(fmt.Sprintf("failed to create updates for VTEC event for %s: %s", handler.vtec.Original, err.Error()))
		return err
	}

	switch handler.vtec.Action {
	case "NEW", "EXB", "EXA":
		err = handler.ugcNEW()
	case "COR":
		err = handler.ugcCOR()
	case "CAN", "UPG", "EXT":
		err = handler.ugcCAN()
	case "CON", "EXP", "ROU":
		err = handler.ugcCON()
	}

	if err != nil {
		handler.logger.Error("failed to handle UGC relations: "+err.Error(), "vtec")
	}

	return nil
}

func (handler *vtecHandler) update() {
	vtec := handler.vtec
	segment := handler.segment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := handler.db.Exec(ctx, `
	UPDATE vtec.events SET updated_at = CURRENT_TIMESTAMP, is_emergency = $6, is_pds = $7 WHERE
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
		polygon = p
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
		pgx.Identifier{"vtec", "updates"},
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

	err = handler.warning(rows[0])
	if err != nil {
		handler.logger.Error("failed to insert/update warning to the database: "+err.Error(), "vtec", vtec)
	}

	return nil
}

func (handler *vtecHandler) ugcNEW() error {
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

	relations := []VTECUGC{}
	for _, ugc := range handler.ugc {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := handler.db.Query(ctx, `
		SELECT * FROM vtec.ugcs WHERE year = $1 AND wfo = $2 AND phenomena = $3 AND
		significance = $4 AND event_number = $5 AND ugc = $6 AND
		action NOT IN ('CAN', 'UPG') AND expires > $7`,
			event.Year, event.WFO, event.Phenomena, event.Significance, event.EventNumber, ugc.ID, expires)

		if err != nil {
			handler.logger.Error("failed to get existing ugc relation: "+err.Error(), "vtec", vtec.Original)
			continue
		}

		duplicates := 0
		deleted := 0
		for rows.Next() {
			if product.isCorrection() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				// Delete the old UGC record
				_, err = handler.db.Exec(ctx, `
				DELETE vtec.ugcs WHERE year = $1 AND wfo = $2 AND phenomena = $3 AND
		significance = $4 AND event_number = $5 AND ugc = $6 AND
		action NOT IN ('NEW', 'EXB', 'EXA')`,
					event.Year, event.WFO, event.Phenomena, event.Significance, event.EventNumber, ugc.ID)

				if err != nil {
					handler.logger.Error("failed to delete duplicate UGC: "+err.Error(), "vtec", vtec.Original, "ugc", ugc.UGC)
					continue
				}

				deleted++
			} else {
				duplicates++
			}
		}

		if deleted > 0 {
			handler.logger.Warn(fmt.Sprintf("Deleted %d deuplicate UGC relations", deleted), "vtec", vtec.Original, "ugc", ugc.UGC)
		}
		if duplicates > 0 {
			handler.logger.Warn(fmt.Sprintf("%d duplicate UGC relations found", duplicates), "vtec", vtec.Original, "ugc", ugc.UGC)
		}

		relations = append(relations, VTECUGC{
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := handler.db.CopyFrom(
		ctx,
		pgx.Identifier{"vtec", "ugcs"},
		[]string{"wfo", "phenomena", "significance", "event_number", "ugc", "issued", "starts", "expires", "ends",
			"end_initial", "action", "year"},
		pgx.CopyFromSlice(len(relations), func(i int) ([]any, error) {
			return []any{relations[i].WFO, relations[i].Phenomena, relations[i].Significance, relations[i].EventNumber, relations[i].UGC,
				relations[i].Issued, relations[i].Starts, relations[i].Expires, relations[i].Ends, relations[i].EndInitial, relations[i].Action,
				relations[i].Year}, nil
		}))
	if err != nil {
		return err
	}

	return nil
}

func (handler *vtecHandler) ugcCOR() error {
	return handler.updateUGC()
}

func (handler *vtecHandler) ugcCAN() error {
	return handler.updateUGC()
}

func (handler *vtecHandler) ugcCON() error {
	return handler.updateUGC()
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := handler.db.Exec(ctx, `
	UPDATE vtec.ugcs SET expires = $1, ends = $2, action = $3 WHERE
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
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				// Find all UGC codes from the state
				rows, err := handler.db.Query(ctx, `
				SELECT ugc FROM postgis.ugcs WHERE state = $1 AND type = $2 AND is_fire = $3 AND valid_to IS NULL;
				`, state.ID, state.Type, isFire)
				if err != nil {
					return fmt.Errorf("error retrieving ALL ugc: %s", err.Error())
				}
				defer rows.Close()

				for rows.Next() {
					var id string
					err := rows.Scan(&id)
					if err != nil {
						return err
					}
					ids = append(ids, id)
				}
				if rows.Err() != nil {
					return fmt.Errorf("error retrieving ALL ugc: %s", rows.Err().Error())
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := handler.db.Query(ctx, `
	SELECT id, ugc FROM postgis.ugcs WHERE ugc = ANY($1) AND valid_to IS NULL
	`, ids)
	if err != nil {
		return fmt.Errorf("error retrieving listed ugcs: %s", err.Error())
	}
	defer rows.Close()

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
	if rows.Err() != nil {
		return fmt.Errorf("error scanning ugcs: %s", rows.Err().Error())
	}

	handler.ugc = ugcs

	return nil
}
