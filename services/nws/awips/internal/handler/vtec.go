package handler

import (
	"fmt"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/handler/util"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type vtecHandler struct {
	*Handler
	eventID   db.VTECEventID
	event     *db.VTECEvent
	warningID db.WarningID
	warning   *db.Warning
	productID *models.RecordID
	product   *awips.TextProduct
	segment   *awips.TextProductSegment
	ugc       map[string]db.UGC
	vtec      awips.VTEC
}

func (handler *Handler) vtec(product *awips.TextProduct, productID *models.RecordID) {

	// Process each segment separately since they reference different UGC areas
	for i, segment := range product.Segments {

		if len(segment.VTEC) == 0 {
			handler.Logger.Info(fmt.Sprintf("Product %s segment %d does not have VTECs. Skipping...", productID, i))
			continue
		}

		// Go through each VTEC in the segment and process it
		for _, vtec := range segment.VTEC {
			var year int
			if vtec.Start != nil {
				year = vtec.Start.Year()
			} else {
				year = product.Issued.Year()
			}

			// Lets check if the VTEC Event is already in the database
			// We need the ID for that
			eventID := db.VTECEventID{
				EventNumber:  vtec.EventNumber,
				Phenomena:    vtec.Phenomena,
				Office:       vtec.WFO,
				Significance: vtec.Significance,
				Year:         year,
			}

			// Find the event
			event, err := surrealdb.Select[db.VTECEvent, models.RecordID](handler.DB, eventID.RecordID())
			if err != nil {
				handler.Logger.Error(fmt.Sprintf("error retrieving current VTEC event: %s", err.Error()))
				continue
			}

			// Events and warnings are basically the same thing but different
			warningID := db.WarningID{
				EventNumber:  vtec.EventNumber,
				Phenomena:    vtec.Phenomena,
				Office:       vtec.WFO,
				Significance: vtec.Significance,
				Year:         year,
			}

			// Find the warning too
			warning, err := surrealdb.Select[db.Warning, models.RecordID](handler.DB, warningID.RecordID())
			if err != nil {
				handler.Logger.Error(fmt.Sprintf("error retrieving current warning: %s", err.Error()))
				continue
			}

			if event.ID == nil {
				// The record does not exist so we will create it
				recordID := models.NewRecordID("vtec_event", eventID)
				phenomena := models.NewRecordID("vtec_phenomena", vtec.Phenomena)
				significance := models.NewRecordID("vtec_significance", vtec.Significance)
				office := models.NewRecordID("office", product.AWIPS.WFO)

				// The database needs a start and end time but VTECs may not have one.
				if vtec.Start == nil {
					// Use the issue time for the start time
					vtec.Start = &product.Issued
				}
				if vtec.End == nil {
					// Use the expiry of the product for the end time
					vtec.End = &segment.Expires
				}

				event, err = surrealdb.Create[db.VTECEvent](handler.DB, models.Table("vtec_event"), db.VTECEvent{
					ID:           &recordID,
					Issued:       &models.CustomDateTime{Time: product.Issued},
					Start:        &models.CustomDateTime{Time: *vtec.Start},
					Expires:      &models.CustomDateTime{Time: segment.UGC.Expires},
					End:          &models.CustomDateTime{Time: *vtec.End},
					EndInitial:   &models.CustomDateTime{Time: *vtec.End},
					Phenomena:    &phenomena,
					Significance: &significance,
					EventNumber:  vtec.EventNumber,
					Office:       &office,
					Title:        vtec.Title(segment.IsEmergency()),
					IsEmergency:  segment.IsEmergency(),
					IsPDS:        segment.IsPDS(),
				})
				if err != nil {
					handler.Logger.Error(fmt.Sprintf("failed to create VTEC event: %s", err.Error()))
					continue
				}
			}

			if warning.ID == nil {
				// The warning does not exists so create it
				recordID := warningID.RecordID()

				warning, err = surrealdb.Create[db.Warning](handler.DB, models.Table("warning"), db.Warning{
					ID:           &recordID,
					Issued:       &models.CustomDateTime{Time: product.Issued},
					Start:        &models.CustomDateTime{Time: *vtec.Start},
					Expires:      &models.CustomDateTime{Time: segment.UGC.Expires},
					End:          &models.CustomDateTime{Time: *vtec.End},
					Phenomena:    event.Phenomena,
					Significance: event.Significance,
					EventNumber:  vtec.EventNumber,
					Office:       event.Office,
					Title:        vtec.Title(segment.IsEmergency()),
					IsEmergency:  segment.IsEmergency(),
					IsPDS:        segment.IsPDS(),
				})
				if err != nil {
					handler.Logger.Error(fmt.Sprintf("failed to create warning: %s", err.Error()))
					continue
				}
			}

			// Update flags if necessary
			if segment.IsEmergency() {
				event.IsEmergency = segment.IsEmergency()
				warning.IsEmergency = segment.IsEmergency()
			}
			if segment.IsPDS() {
				event.IsPDS = segment.IsPDS()
				warning.IsPDS = segment.IsPDS()
			}

			vh := vtecHandler{
				Handler:   handler,
				eventID:   eventID,
				event:     event,
				warningID: warningID,
				warning:   warning,
				productID: productID,
				product:   product,
				segment:   &segment,
				vtec:      vtec,
			}

			err = vh.segmentUGC()
			if err != nil {
				handler.Logger.Error(err.Error())
				continue
			}

			vh.handle()
		}

	}
}

// Handle the VTEC product
func (handler *vtecHandler) handle() {
	event := handler.event
	warning := handler.warning
	product := handler.product
	segment := handler.segment
	vtec := handler.vtec

	// The product expires at the UGC expiry time
	expires := segment.UGC.Expires
	var end time.Time
	if vtec.End == nil {
		end = expires
		handler.Logger.Info("VTEC end time is nil. Defaulting to UGC expiry time.")
	} else {
		end = *vtec.End
	}

	switch vtec.Action {
	case "CAN":
		fallthrough
	case "UPG":
		event.Expires.Time = segment.UGC.Expires
		warning.Expires.Time = segment.UGC.Expires
		event.End.Time = product.Issued.UTC()
		warning.End.Time = product.Issued.UTC()
	case "EXP":
		event.Expires.Time = end
		warning.Expires.Time = end
		event.End.Time = end
		warning.End.Time = end
	case "EXT":
		fallthrough
	case "EXB":
		event.End.Time = end
		warning.End.Time = end
		event.Expires.Time = segment.UGC.Expires
		warning.Expires.Time = segment.UGC.Expires
	default:
		// NEW and CON
		if event.End.Time.Before(end) {
			event.End.Time = end
		}
		if warning.End.Time.Before(end) {
			warning.End.Time = end
		}
		if event.Expires.Time.Before(segment.Expires) {
			event.Expires.Time = segment.Expires
		}
		if warning.Expires.Time.Before(segment.Expires) {
			warning.Expires.Time = segment.Expires
		}
	}

	historyID, warningID, err := handler.createHistoryRecords()
	if err != nil {
		handler.Logger.Error(err.Error())
		return
	}

	if vtec.Action == "NEW" || vtec.Action == "EXA" {
		handler.relateUGC(historyID, warningID)
	} else {
		handler.updateUGC(historyID, warningID)
	}

	_, err = surrealdb.Merge[db.VTECEvent, models.RecordID](handler.DB, *event.ID, event)
	if err != nil {
		handler.Logger.Error(fmt.Sprintf("error updating %s: ", event.ID.String()))
	}

	_, err = surrealdb.Merge[db.Warning, models.RecordID](handler.DB, *warning.ID, warning)
	if err != nil {
		handler.Logger.Error(fmt.Sprintf("error updating %s: ", warning.ID.String()))
	}
}

// Create the historical VTEC and warning records
func (handler *vtecHandler) createHistoryRecords() (*db.VTECHistoryID, *db.WarningHistoryID, error) {
	event := handler.event
	product := handler.product
	segment := handler.segment
	vtec := handler.vtec
	ugcs := handler.ugc

	// Get any polygons in the product
	var latlon *db.LatLon
	var polygon *models.GeometryPolygon
	var bbox models.GeometryMultiPoint
	if segment.LatLon != nil {
		output := util.LatLonFromAwips(*segment.LatLon)
		latlon = &output
		polygon = &latlon.Points
		handler.warning.Polygon = polygon
		bbox = util.BBoxFromAwips(*segment.LatLon)
	}

	// Get any TML data
	var tml *db.TML
	tmlAwips, err := awips.ParseTML(segment.Text, product.Issued)
	if err != nil {
		handler.Logger.Error("error parsing TML: " + err.Error())
	}

	if tmlAwips != nil {
		tml = &db.TML{
			Direction:   tmlAwips.Direction,
			Location:    models.NewGeometryPoint(tmlAwips.Location[0], tmlAwips.Location[1]),
			Speed:       tmlAwips.Speed,
			SpeedString: tmlAwips.SpeedString,
			Time:        models.CustomDateTime{Time: tmlAwips.Time},
			Original:    tmlAwips.Original,
		}
	}

	mpolygons := []models.GeometryMultiPolygon{}
	ugc := []*models.RecordID{}
	for _, i := range ugcs {
		ugc = append(ugc, i.ID)
		if polygon == nil {
			mpolygons = append(mpolygons, *i.Geometry)
		}
	}

	if polygon == nil {
		bbox = util.BBoxFromMultiPolygon(mpolygons)
	}

	action := models.NewRecordID("vtec_action", vtec.Action)

	// The product expires at the UGC expiry time
	expires := segment.UGC.Expires
	end := event.End

	vtecHistoryID := db.VTECHistoryID{
		EventNumber:  vtec.EventNumber,
		Phenomena:    vtec.Phenomena,
		Office:       vtec.WFO,
		Significance: vtec.Significance,
		Year:         handler.eventID.Year,
		Sequence:     handler.event.Updates, // Starting from 0 is fine
	}

	// Create the historical record
	historyRecord, err := surrealdb.Create[db.VTECHistory](handler.DB, models.Table("vtec_history"), db.VTECHistory{
		ID:           vtecHistoryID.RecordID(),
		Issued:       &models.CustomDateTime{Time: product.Issued},
		Start:        event.Start,
		Expires:      &models.CustomDateTime{Time: expires},
		End:          end,
		Original:     segment.Text,
		Title:        vtec.Title(segment.IsEmergency()),
		Action:       &action,
		Phenomena:    event.Phenomena,
		Office:       event.Office,
		Significance: event.Significance,
		EventNumber:  vtec.EventNumber,
		VTEC: db.VTEC{
			Class:        vtec.Class,
			Action:       vtec.Action,
			WFO:          vtec.WFO,
			Phenomena:    vtec.Phenomena,
			Significance: vtec.Significance,
			EventNumber:  vtec.EventNumber,
			Start:        vtec.StartString,
			End:          vtec.EndString,
		},
		IsEmergency: segment.IsEmergency(),
		IsPDS:       segment.IsPDS(),
		LatLon:      latlon,
		Polygon:     polygon,
		BBox:        bbox,
		Tags:        segment.Tags,
		TML:         tml,
		Product:     handler.productID,
		UGC:         ugc,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating %s: %s", vtecHistoryID.String(), err.Error())
	}

	// Relate the event to the historical record
	err = surrealdb.Relate(handler.DB, &surrealdb.Relationship{
		In:       *event.ID,
		Out:      *historyRecord.ID,
		Relation: models.Table("vtec_event_history"),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error relating vtec_event_history %s: %s", vtecHistoryID.String(), err.Error())
	}

	// Update the event's updates
	event.Updates++

	event.Title = vtec.Title(segment.IsEmergency())

	// Now for the warning
	warningHistoryID := db.WarningHistoryID{
		EventNumber:  vtec.EventNumber,
		Phenomena:    vtec.Phenomena,
		Office:       vtec.WFO,
		Significance: vtec.Significance,
		Year:         handler.warningID.Year,
		Sequence:     handler.warning.Updates, // Starting from 0 is fine
	}

	// Create the historical warning
	warningRecord, err := surrealdb.Create[db.WarningHistory](handler.DB, models.Table("warning_history"), db.WarningHistory{
		ID:           warningHistoryID.RecordID(),
		Issued:       &models.CustomDateTime{Time: product.Issued},
		Start:        event.Start,
		Expires:      &models.CustomDateTime{Time: expires},
		End:          end,
		Original:     segment.Text,
		Title:        vtec.Title(segment.IsEmergency()),
		Action:       &action,
		Phenomena:    event.Phenomena,
		Office:       event.Office,
		Significance: event.Significance,
		EventNumber:  vtec.EventNumber,
		IsEmergency:  segment.IsEmergency(),
		IsPDS:        segment.IsPDS(),
		Polygon:      polygon,
		BBox:         bbox,
		Tags:         segment.Tags,
		TML:          tml,
		Product:      handler.productID,
		UGC:          ugc,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating %s: %s", warningHistoryID.String(), err.Error())
	}

	// Relate the event to the warning record
	err = surrealdb.Relate(handler.DB, &surrealdb.Relationship{
		In:       *handler.warning.ID,
		Out:      *warningRecord.ID,
		Relation: models.Table("warning_event_history"),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error relating warning_event_history %s: %s", warningHistoryID.String(), err.Error())
	}

	handler.warning.Updates++

	handler.warning.Title = vtec.Title(segment.IsEmergency())
	handler.warning.BBox = bbox

	return &vtecHistoryID, &warningHistoryID, nil
}

func (handler *vtecHandler) relateUGC(historyID *db.VTECHistoryID, warningID *db.WarningHistoryID) {
	product := handler.product
	event := handler.event
	warning := handler.warning
	segment := handler.segment
	vtec := handler.vtec
	ugcs := handler.ugc

	action := models.NewRecordID("vtec_action", vtec.Action)

	// The product expires at the UGC expiry time
	expires := segment.UGC.Expires
	var end time.Time
	if vtec.End == nil {
		end = expires
		handler.Logger.Info("VTEC end time is nil. Defaulting to UGC expiry time.")
	} else {
		end = *vtec.End
	}

	for _, ugc := range ugcs {
		err := surrealdb.Relate(handler.DB, &surrealdb.Relationship{
			In:       *event.ID,
			Out:      *ugc.ID,
			Relation: models.Table("vtec_ugc"),
			Data: map[string]any{
				"id": models.NewRecordID("vtec_ugc", db.VTECUGCID{
					EventNumber:  vtec.EventNumber,
					Phenomena:    vtec.Phenomena,
					Office:       vtec.WFO,
					Significance: vtec.Significance,
					Year:         handler.eventID.Year,
					UGC:          fmt.Sprintf("%v", ugc.ID.ID),
				}),
				"in":          event.ID,
				"out":         ugc.ID,
				"created_at":  &models.CustomDateTime{Time: time.Now().UTC()},
				"updated_at":  &models.CustomDateTime{Time: time.Now().UTC()},
				"issued":      &models.CustomDateTime{Time: product.Issued},
				"start":       event.Start,
				"expires":     &models.CustomDateTime{Time: expires},
				"end":         &models.CustomDateTime{Time: end},
				"end_initial": &models.CustomDateTime{Time: end},
				"action":      &action,
				"latest":      historyID.RecordID(),
			},
		})
		if err != nil {
			handler.Logger.Error(fmt.Sprintf("error relating %s: %s", ugc.ID.String(), err.Error()))
			continue
		}

		err = surrealdb.Relate(handler.DB, &surrealdb.Relationship{
			In:       *warning.ID,
			Out:      *ugc.ID,
			Relation: models.Table("warning_ugc"),
			Data: map[string]any{
				"id": models.NewRecordID("warning_ugc", db.WarningUGCID{
					EventNumber:  vtec.EventNumber,
					Phenomena:    vtec.Phenomena,
					Office:       vtec.WFO,
					Significance: vtec.Significance,
					Year:         handler.eventID.Year,
					UGC:          fmt.Sprintf("%v", ugc.ID.ID),
				}),
				"in":          warning.ID,
				"out":         ugc.ID,
				"created_at":  &models.CustomDateTime{Time: time.Now().UTC()},
				"updated_at":  &models.CustomDateTime{Time: time.Now().UTC()},
				"issued":      &models.CustomDateTime{Time: product.Issued},
				"start":       event.Start,
				"expires":     &models.CustomDateTime{Time: expires},
				"end":         &models.CustomDateTime{Time: end},
				"end_initial": &models.CustomDateTime{Time: end},
				"action":      &action,
				"latest":      warningID.RecordID(),
			},
		})
		if err != nil {
			handler.Logger.Error(fmt.Sprintf("error relating %s to %s: %s", ugc.ID.String(), warning.ID.String(), err.Error()))
		}
	}
}

func (handler *vtecHandler) updateUGC(historyID *db.VTECHistoryID, warningID *db.WarningHistoryID) {
	segment := handler.segment
	vtec := handler.vtec
	ugc := handler.ugc

	vtecUgcs := []models.RecordID{}
	warningUgcs := []models.RecordID{}

	for _, ugc := range ugc {
		id := db.VTECUGCID{
			EventNumber:  vtec.EventNumber,
			Phenomena:    vtec.Phenomena,
			Office:       vtec.WFO,
			Significance: vtec.Significance,
			Year:         handler.eventID.Year,
			UGC:          fmt.Sprintf("%v", ugc.ID.ID),
		}
		vtecUgcs = append(vtecUgcs, models.NewRecordID("vtec_ugc", id))

		wID := db.WarningUGCID{
			EventNumber:  vtec.EventNumber,
			Phenomena:    vtec.Phenomena,
			Office:       vtec.WFO,
			Significance: vtec.Significance,
			Year:         handler.eventID.Year,
			UGC:          fmt.Sprintf("%v", ugc.ID.ID),
		}
		warningUgcs = append(warningUgcs, models.NewRecordID("warning_ugc", wID))
	}

	action := models.NewRecordID("vtec_action", vtec.Action)

	// The product expires at the UGC expiry time
	expires := segment.UGC.Expires
	end := handler.event.End

	res, err := surrealdb.Merge[[]db.VTECUGC](handler.DB, vtecUgcs, map[string]interface{}{
		"expires": expires.UTC(),
		"end":     end.UTC(),
		"action":  action,
		"latest":  historyID.RecordID(),
	})
	if err != nil {
		handler.Logger.Error("error updating VTEC UGC: " + err.Error())
		return
	}

	if len(*res) == 0 {
		handler.Logger.Info(fmt.Sprintf("Missing UGC relation(s) for %s. Creating now.", handler.eventID.String()))
		handler.relateUGC(historyID, warningID)
		return
	}

	_, err = surrealdb.Merge[[]db.WarningUGC](handler.DB, warningUgcs, map[string]interface{}{
		"expires": expires.UTC(),
		"end":     end.UTC(),
		"action":  action,
		"latest":  warningID.RecordID(),
	})
	if err != nil {
		handler.Logger.Error("error updating Warning UGC: " + err.Error())
		return
	}
}

// Creates an array of UGC record IDs for the database
func (handler *vtecHandler) segmentUGC() error {
	segment := handler.segment

	ids := []models.RecordID{}
	// For each state...
	for _, state := range segment.UGC.States {
		// ...and for each area...
		for _, area := range state.Areas {
			ugcType := state.Type
			// Fire weather (FW) events have different zones
			if handler.vtec.Phenomena == "FW" {
				ugcType = "F"
			}
			if area == "000" || area == "ALL" {
				// Find all UGC codes from the state
				key := state.ID + ugcType
				fmt.Println(key)
				res, err := surrealdb.Query[[]struct {
					ID models.RecordID `json:"id"`
				}](handler.DB, fmt.Sprintf(`SELECT id FROM ugc WHERE string::slice(record::id(id), 0, 3) == "%s"`, key), map[string]interface{}{})
				if err != nil {
					return fmt.Errorf("error retrieving ALL ugc: %s", err.Error())
				}

				result := (*res)[0].Result

				if len(result) == 0 {
					return fmt.Errorf("got 0 UGC for %s. Expected ALL", key)
				}

				for _, i := range result {
					ids = append(ids, i.ID)
				}
			} else {
				// Get the needed UGCs
				id := state.ID + ugcType + area
				ids = append(ids, models.NewRecordID("ugc", id))
			}
		}
	}

	res, err := surrealdb.Select[[]db.UGC, []models.RecordID](handler.DB, ids)
	if err != nil {
		return fmt.Errorf("error retrieving ugcs: %s", err.Error())
	}

	ugcs := map[string]db.UGC{}
	for _, i := range *res {
		ugcs[i.ID.String()] = i
	}

	handler.ugc = ugcs

	return nil
}
