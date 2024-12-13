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
	Product      *awips.TextProduct
	DBProduct    *db.Product
	Segment      awips.TextProductSegment
	ParentRecord *db.VTECEvent
	ParentID     db.VTECEventID
	VTEC         awips.VTEC
}

func (handler *Handler) vtec(product *awips.TextProduct, dbProduct *db.Product) error {

	events := map[db.VTECEventID]*db.VTECEvent{}

	for i, segment := range product.Segments {

		if len(segment.VTEC) == 0 {
			handler.Logger.Info(fmt.Sprintf("Product %s segment %d does not have VTECs. Skipping...", dbProduct.ID, i))
			continue
		}

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

			parent := events[eventID]

			if parent == nil {

				eventRes, err := surrealdb.Query[[]db.VTECEvent](handler.DB, "SELECT * FROM $id", map[string]interface{}{
					"id": models.NewRecordID("vtec_event", eventID),
				})
				if err != nil {
					return fmt.Errorf("error retrieving current VTECs: %s", err)
				}

				currentVTECEvents := (*eventRes)[0].Result

				if len(currentVTECEvents) == 0 {
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

					res, err := surrealdb.Create[db.VTECEvent](handler.DB, models.Table("vtec_event"), db.VTECEvent{
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
						Title:        fmt.Sprintf("%s %s", awips.VTECPhenomena[vtec.Phenomena], awips.VTECSignificance[vtec.Significance]),
						IsEmergency:  segment.IsEmergency(),
						IsPDS:        segment.IsPDS(),
					})
					if err != nil {
						return fmt.Errorf("failed to create vtec_event: " + err.Error())
					}

					// The new record can now be our parent
					parent = res
				} else {
					// We have the parent record
					if len(currentVTECEvents) > 1 {
						// Somehow there are multiple events with the same ID. Should not be possible
						handler.Logger.Warn(fmt.Sprintf("Found %d records for %s when there should be 1 or 0", len(currentVTECEvents), eventID.String()))
					}
					parent = &currentVTECEvents[0]
					parent.UpdatedAt = &models.CustomDateTime{Time: time.Now()}

					if segment.IsEmergency() {
						parent.IsEmergency = true
					}
					if segment.IsPDS() {
						parent.IsPDS = true
					}
				}

				events[eventID] = parent
			}

			vh := vtecHandler{
				Handler:      handler,
				Product:      product,
				DBProduct:    dbProduct,
				Segment:      segment,
				ParentRecord: parent,
				ParentID:     eventID,
				VTEC:         vtec,
			}

			var err error
			switch vtec.Action {
			case "NEW":
				err = vh.handleNew()
			case "CON":
				fallthrough
			case "EXA":
				fallthrough
			case "EXT":
				err = vh.handleContinue()
			case "CAN":
				fallthrough
			case "EXP":
				fallthrough
			case "UPG":
				err = vh.handleCancel()
			}

			if err != nil {
				return err
			}
		}
	}

	for id, event := range events {

		// Check if all the UGCs have expired
		ugcRes, err := surrealdb.Query[[]db.VTECUGC](handler.DB, fmt.Sprintf("SELECT id, end, out FROM %s->vtec_ugc WHERE end > time::now()", id.String()), map[string]interface{}{})
		if err != nil {
			handler.Logger.Error("error getting expired UGCs: " + err.Error())
			continue
		}

		// If all UGCs are expired then the event has ended
		if len((*ugcRes)[0].Result) == 0 {
			event.End = &models.CustomDateTime{Time: product.Issued}
		}

		_, err = surrealdb.Merge[db.VTECEvent](handler.DB, *event.ID, event)
		if err != nil {
			handler.Logger.Error("error merging vtec_event: " + err.Error())
			continue
		}
	}

	return nil
}

// Handles action NEW
func (handler *vtecHandler) handleNew() error {
	parent := handler.ParentRecord

	historyID, err := handler.createVTECHistory()
	if err != nil {
		return err
	}

	historyRecordID := models.NewRecordID("vtec_history", historyID)

	// Relate the event to the historical record
	err = surrealdb.Relate(handler.DB, &surrealdb.Relationship{
		In:       *parent.ID,
		Out:      historyRecordID,
		Relation: models.Table("vtec_event_history"),
	})
	if err != nil {
		return fmt.Errorf("error relating vtec_event_history %s: %s", historyRecordID.String(), err.Error())
	}

	handler.relateUGC(historyID)

	return nil
}

// Handles actions CON, EXA, EXT, EXB
func (handler *vtecHandler) handleContinue() error {
	segment := handler.Segment
	vtec := handler.VTEC
	parent := handler.ParentRecord

	historyID, err := handler.createVTECHistory()
	if err != nil {
		return err
	}

	historyRecordID := models.NewRecordID("vtec_history", historyID)

	// Relate the event to the historical record
	err = surrealdb.Relate(handler.DB, &surrealdb.Relationship{
		In:       *parent.ID,
		Out:      historyRecordID,
		Relation: models.Table("vtec_event_history"),
	})
	if err != nil {
		return fmt.Errorf("error relating vtec_event_history %s: %s", historyRecordID.String(), err.Error())
	}

	if handler.VTEC.Action == "EXT" || handler.VTEC.Action == "EXB" {
		parent.End = &models.CustomDateTime{Time: *vtec.End}
		parent.Expires = &models.CustomDateTime{Time: segment.UGC.Expires}
	}

	if handler.VTEC.Action == "EXA" || handler.VTEC.Action == "EXB" {
		handler.relateUGC(historyID)
	} else {
		handler.updateUGC(historyID)
	}

	return nil
}

// Handle action CAN, EXP, UPG
func (handler *vtecHandler) handleCancel() error {
	parent := handler.ParentRecord

	historyID, err := handler.createVTECHistory()
	if err != nil {
		return err
	}

	historyRecordID := models.NewRecordID("vtec_history", historyID)

	// Relate the event to the historical record
	err = surrealdb.Relate(handler.DB, &surrealdb.Relationship{
		In:       *parent.ID,
		Out:      historyRecordID,
		Relation: models.Table("vtec_event_history"),
	})
	if err != nil {
		return fmt.Errorf("error relating vtec_event_history %s: %s", historyRecordID.String(), err.Error())
	}

	handler.cancelUGC(historyID)

	return nil
}

func (handler *vtecHandler) createVTECHistory() (*db.VTECHistoryID, error) {
	product := handler.Product
	segment := handler.Segment
	parent := handler.ParentRecord
	parentID := handler.ParentID
	vtec := handler.VTEC

	// Find out how many other history records there are for this VTEC
	countRes, err := surrealdb.Query[[]struct {
		Count int `json:"count"`
	}](handler.DB, fmt.Sprintf("SELECT count() FROM %s->vtec_event_history->vtec_history GROUP ALL;", parentID.String()), map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("error getting history sequence for %s: %s", parentID.String(), err.Error())
	}

	sequence := 0
	if len((*countRes)[0].Result) > 0 {
		sequence = (*countRes)[0].Result[0].Count
	}

	// We have our sequence so we can create the historical record ID
	historyID := handler.vtecHistoryID(sequence)

	action := models.NewRecordID("vtec_action", vtec.Action)

	// Get any polygons in the product
	var latlon *db.LatLon
	var polygon *models.GeometryPolygon
	if segment.LatLon != nil {
		output := util.LatLonFromAwips(*segment.LatLon)
		latlon = &output
		polygon = &latlon.Points
	}

	// Generate UGC array
	ugcs := handler.vtecHistoryUGC()

	// Any TML data
	tml, err := handler.vtecTML()

	expires := segment.UGC.Expires
	if vtec.End == nil {
		vtec.End = &handler.Segment.Expires
		handler.Logger.Info("VTEC end time is nil. Defaulting to UGC expiry time.")
	}
	end := *vtec.End

	if vtec.Action == "CAN" || vtec.Action == "UPG" {
		expires = product.Issued
		end = product.Issued
	} else if vtec.Action == "EXP" {
		expires = *vtec.End
	}

	title := parent.Title
	if segment.IsEmergency() {
		title = fmt.Sprintf("%s Emergency", awips.VTECPhenomena[vtec.Phenomena])
	}

	// Create the historical product
	historyRecord, err := surrealdb.Create[db.VTECHistory](handler.DB, models.Table("vtec_history"), db.VTECHistory{
		ID:           historyID.RecordID(),
		Issued:       &models.CustomDateTime{Time: product.Issued},
		Start:        parent.Start,
		Expires:      &models.CustomDateTime{Time: expires},
		End:          &models.CustomDateTime{Time: end},
		Original:     segment.Text,
		Title:        title,
		Action:       &action,
		Phenomena:    parent.Phenomena,
		Office:       parent.Office,
		Significance: parent.Significance,
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
		Tags:        segment.Tags,
		TML:         tml,
		Product:     handler.DBProduct.ID,
		UGC:         ugcs,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating %s: %s", historyID.String(), err.Error())
	}

	group := fmt.Sprintf("%s.%s.%04d.%s.%d", historyID.Phenomena, historyID.Significance, historyID.EventNumber, historyID.Office, historyID.Year)

	currentWarnings, err := surrealdb.Query[[]db.Warning](handler.DB, fmt.Sprintf("SELECT id, valid_to, created_at FROM warning WHERE group == \"%s\" ORDER BY created_at ASC", group), map[string]interface{}{})
	if err != nil {
		return &historyID, fmt.Errorf("error retrieving warning group %s: %s", group, err.Error())
	}

	if len((*currentWarnings)[0].Result) > 0 {
		currentWarning := (*currentWarnings)[0].Result[0]

		_, err := surrealdb.Merge[db.Warning](handler.DB, *currentWarning.ID, map[string]interface{}{
			"valid_to": historyRecord.Issued,
		})
		if err != nil {
			return &historyID, fmt.Errorf("error merging warning group %s: %s", group, err.Error())
		}
	}

	warningID := db.WarningID{
		EventNumber:  historyID.EventNumber,
		Phenomena:    historyID.Phenomena,
		Office:       historyID.Office,
		Significance: historyID.Significance,
		Sequence:     historyID.Sequence,
	}

	warning := &db.Warning{
		ID:                    warningID.RecordID(),
		Group:                 group,
		Start:                 historyRecord.Start,
		Expires:               historyRecord.Expires,
		End:                   historyRecord.End,
		ValidFrom:             &models.CustomDateTime{Time: time.Now().UTC()},
		ValidTo:               historyRecord.End,
		Text:                  historyRecord.Original,
		Title:                 historyRecord.Title,
		Action:                historyRecord.Action,
		Phenomena:             historyRecord.Phenomena,
		Office:                historyRecord.Office,
		Significance:          historyRecord.Significance,
		PhenomenaSignificance: warningID.Phenomena + "." + warningID.Significance,
		EventNumber:           historyRecord.EventNumber,
		VTEC:                  historyRecord.VTEC,
		HVTEC:                 historyRecord.HVTEC,
		IsEmergency:           historyRecord.IsEmergency,
		IsPDS:                 historyRecord.IsPDS,
		Polygon:               historyRecord.Polygon,
		Tags:                  historyRecord.Tags,
		TML:                   historyRecord.TML,
		UGC:                   ugcs,
	}

	_, err = surrealdb.Create[db.Warning, models.RecordID](handler.DB, *warningID.RecordID(), warning)
	if err != nil {
		handler.Logger.Error(fmt.Sprintf("error creating warning %s: %s", warningID.String(), err.Error()))
	}

	return &historyID, nil
}

func (handler *vtecHandler) vtecHistoryID(sequence int) db.VTECHistoryID {
	vtec := handler.VTEC
	return db.VTECHistoryID{
		EventNumber:  vtec.EventNumber,
		Phenomena:    vtec.Phenomena,
		Office:       vtec.WFO,
		Significance: vtec.Significance,
		Year:         handler.ParentID.Year,
		Sequence:     sequence,
	}
}

func (handler *vtecHandler) vtecHistoryUGC() []*models.RecordID {
	segment := handler.Segment

	ugcs := []*models.RecordID{}
	for _, state := range segment.UGC.States {
		for _, area := range state.Areas {
			ugcType := state.Type
			if handler.VTEC.Phenomena == "FW" {
				ugcType = "F"
			}
			if area == "000" || area == "ALL" {
				key := state.ID + ugcType
				for k, ugc := range handler.UGCData {
					if k[0:3] == key {
						ugcs = append(ugcs, ugc.ID)
					}
				}
			} else {
				id := state.ID + ugcType + area
				ugc := handler.UGCData[id]
				if ugc.ID == nil {
					handler.Logger.Warn(fmt.Sprintf("Could not find UGC %s", id))
					continue
				}
				ugcs = append(ugcs, ugc.ID)
			}
		}
	}

	return ugcs
}

func (handler *vtecHandler) vtecTML() (*db.TML, error) {
	product := handler.Product
	segment := handler.Segment

	tmlAwips, err := awips.ParseTML(segment.Text, product.Issued)
	if err != nil {
		return nil, err
	}

	if tmlAwips == nil {
		return nil, nil
	}

	tml := db.TML{
		Direction:   tmlAwips.Direction,
		Location:    models.NewGeometryPoint(tmlAwips.Location[0], tmlAwips.Location[1]),
		Speed:       tmlAwips.Speed,
		SpeedString: tmlAwips.SpeedString,
		Time:        models.CustomDateTime{Time: tmlAwips.Time},
		Original:    tmlAwips.Original,
	}

	return &tml, nil
}

func (handler *vtecHandler) relateUGC(historyID *db.VTECHistoryID) {
	product := handler.Product
	parent := handler.ParentRecord
	segment := handler.Segment
	vtec := handler.VTEC

	records := handler.vtecHistoryUGC()

	action := models.NewRecordID("vtec_action", vtec.Action)

	for _, id := range records {
		err := surrealdb.Relate(handler.DB, &surrealdb.Relationship{
			In:       *parent.ID,
			Out:      *id,
			Relation: models.Table("vtec_ugc"),
			Data: map[string]any{
				"id": models.NewRecordID("vtec_ugc", db.VTECUGCID{
					EventNumber:  vtec.EventNumber,
					Phenomena:    vtec.Phenomena,
					Office:       vtec.WFO,
					Significance: vtec.Significance,
					Year:         handler.ParentID.Year,
					UGC:          fmt.Sprintf("%v", id.ID),
				}),
				"created_at":  &models.CustomDateTime{Time: time.Now().UTC()},
				"issued":      &models.CustomDateTime{Time: product.Issued},
				"start":       parent.Start,
				"expires":     &models.CustomDateTime{Time: segment.UGC.Expires},
				"end":         &models.CustomDateTime{Time: *vtec.End},
				"end_initial": &models.CustomDateTime{Time: *vtec.End},
				"action":      &action,
				"latest":      historyID.RecordID(),
			},
		})
		if err != nil {
			handler.Logger.Error(fmt.Sprintf("error relating %s: %s", id.String(), err.Error()))
		}
	}
}

func (handler *vtecHandler) updateUGC(historyID *db.VTECHistoryID) {
	segment := handler.Segment
	vtec := handler.VTEC

	records := handler.vtecHistoryUGC()

	ugcs := []models.RecordID{}

	for _, ugc := range records {
		id := db.VTECUGCID{
			EventNumber:  vtec.EventNumber,
			Phenomena:    vtec.Phenomena,
			Office:       vtec.WFO,
			Significance: vtec.Significance,
			Year:         handler.ParentID.Year,
			UGC:          fmt.Sprintf("%v", ugc.ID),
		}
		ugcs = append(ugcs, models.NewRecordID("vtec_ugc", id))
	}

	action := models.NewRecordID("vtec_action", vtec.Action)

	expires := models.CustomDateTime{Time: segment.UGC.Expires}
	end := models.CustomDateTime{Time: *vtec.End}

	res, err := surrealdb.Merge[[]db.VTECUGC](handler.DB, ugcs, map[string]interface{}{
		"expires": expires,
		"end":     end,
		"action":  action,
		"latest":  historyID,
	})
	if err != nil {
		handler.Logger.Error("error updating VTEC UGC: " + err.Error())
		return
	}

	if len(*res) == 0 {
		handler.Logger.Info(fmt.Sprintf("Missing UGC relation for %s. Creating now.", &handler.ParentID))
		handler.relateUGC(historyID)
	}
}

func (handler *vtecHandler) cancelUGC(historyID *db.VTECHistoryID) {
	vtec := handler.VTEC

	records := handler.vtecHistoryUGC()

	ugcs := []models.RecordID{}

	for _, ugc := range records {
		id := db.VTECUGCID{
			EventNumber:  vtec.EventNumber,
			Phenomena:    vtec.Phenomena,
			Office:       vtec.WFO,
			Significance: vtec.Significance,
			Year:         handler.ParentID.Year,
			UGC:          fmt.Sprintf("%v", ugc.ID),
		}
		ugcs = append(ugcs, models.NewRecordID("vtec_ugc", id))
	}

	action := models.NewRecordID("vtec_action", vtec.Action)

	expires := models.CustomDateTime{Time: handler.Product.Issued}
	end := models.CustomDateTime{Time: handler.Product.Issued}

	res, err := surrealdb.Merge[[]db.VTECUGC](handler.DB, ugcs, map[string]interface{}{
		"expires": expires,
		"end":     end,
		"action":  action,
		"latest":  historyID,
	})
	if err != nil {
		handler.Logger.Error("error updating VTEC UGC: " + err.Error())
		return
	}

	if len(*res) == 0 {
		handler.Logger.Info(fmt.Sprintf("Missing UGC relation for %s. Creating now.", &handler.ParentID))
		handler.relateUGC(historyID)
	}
}
