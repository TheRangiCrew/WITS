package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/twpayne/go-geos"
)

// Warning
type Warning struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at,omitzero"`
	UpdatedAt     time.Time
	Issued        time.Time `json:"issued"`
	Starts        time.Time `json:"starts,omitzero"`
	Expires       time.Time `json:"expires"`
	Ends          time.Time `json:"ends,omitzero"`
	EndInitial    time.Time
	Text          string     `json:"text"`
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

func (handler *vtecHandler) warning(update VTECUpdate) error {
	// Lets check if the Warning is already in the database
	rows, err := handler.db.Query(handler.db.CTX, `
			SELECT * FROM warning WHERE
			wfo = $1 AND phenomena = $2 AND significance = $3 AND event_number = $4 AND year = $5
			`, update.WFO, update.Phenomena, update.Significance, update.EventNumber, update.Year)
	if err != nil {
		handler.logger.Error("failed to get warning: " + err.Error())
		return err
	}

	var warning *Warning

	if !rows.Next() {
		if update.Action != "CAN" && update.Action != "EXP" && update.Action != "UPG" {
			warning = &Warning{
				Issued:        update.Issued,
				Starts:        update.Starts,
				Expires:       update.Expires,
				Ends:          update.Ends,
				EndInitial:    update.Ends,
				Text:          update.Text,
				WFO:           update.WFO,
				Action:        update.Action,
				Class:         update.Class,
				Phenomena:     update.Phenomena,
				Significance:  update.Significance,
				EventNumber:   update.EventNumber,
				Year:          update.Year,
				Title:         update.Title,
				IsEmergency:   update.IsEmergency,
				IsPDS:         update.IsPDS,
				Polygon:       update.Polygon,
				Direction:     update.Direction,
				Location:      update.Location,
				Speed:         update.Speed,
				SpeedText:     update.SpeedText,
				TMLTime:       update.TMLTime,
				UGC:           update.UGC,
				Tornado:       update.Tornado,
				Damage:        update.Damage,
				HailThreat:    update.HailThreat,
				HailTag:       update.HailTag,
				WindThreat:    update.WindThreat,
				WindTag:       update.WindTag,
				FlashFlood:    update.FlashFlood,
				RainfallTag:   update.RainfallTag,
				FloodTagDam:   update.FloodTagDam,
				SpoutTag:      update.SpoutTag,
				SnowSquall:    update.SnowSquall,
				SnowSquallTag: update.SnowSquallTag,
			}

			_, err := handler.db.CopyFrom(
				context.Background(),
				pgx.Identifier{"warning"},
				[]string{"issued", "starts", "expires", "ends", "end_initial", "text", "wfo", "action", "class", "phenomena", "significance", "event_number", "year", "title", "is_emergency", "is_pds", "polygon", "direction", "location", "speed", "speed_text", "tml_time", "ugc", "tornado", "damage", "hail_threat", "hail_tag", "wind_threat", "wind_tag", "flash_flood", "rainfall_tag", "flood_tag_dam", "spout_tag", "snow_squall", "snow_squall_tag"},
				pgx.CopyFromRows([][]interface{}{
					{
						warning.Issued, warning.Starts, warning.Expires, warning.Ends, warning.EndInitial, warning.Text, warning.WFO, warning.Action, warning.Class, warning.Phenomena, warning.Significance, warning.EventNumber, warning.Year, warning.Title, warning.IsEmergency, warning.IsPDS, warning.Polygon, warning.Direction, warning.Location, warning.Speed, warning.SpeedText, warning.TMLTime, warning.UGC, warning.Tornado, warning.Damage, warning.HailThreat, warning.HailTag, warning.WindThreat, warning.WindTag, warning.FlashFlood, warning.RainfallTag, warning.FloodTagDam, warning.SpoutTag, warning.SnowSquall, warning.SnowSquallTag,
					},
				}),
			)

			if err != nil {
				handler.logger.Error(fmt.Sprintf("failed to create warning: %s", err.Error()))
				return err
			}
		}
	} else {
		warning = &Warning{}
		err = rows.Scan(
			&warning.ID,
			&warning.CreatedAt,
			&warning.UpdatedAt,
			&warning.Issued,
			&warning.Starts,
			&warning.Expires,
			&warning.Ends,
			&warning.EndInitial,
			&warning.Text,
			&warning.WFO,
			&warning.Action,
			&warning.Class,
			&warning.Phenomena,
			&warning.Significance,
			&warning.EventNumber,
			&warning.Year,
			&warning.Title,
			&warning.IsEmergency,
			&warning.IsPDS,
			&warning.Polygon,
			&warning.Direction,
			&warning.Location,
			&warning.Speed,
			&warning.SpeedText,
			&warning.TMLTime,
			&warning.UGC,
			&warning.Tornado,
			&warning.Damage,
			&warning.HailThreat,
			&warning.HailTag,
			&warning.WindThreat,
			&warning.WindTag,
			&warning.FlashFlood,
			&warning.RainfallTag,
			&warning.FloodTagDam,
			&warning.SpoutTag,
			&warning.SnowSquall,
			&warning.SnowSquallTag)

		if err != nil {
			handler.logger.Error("failed to scan warning: " + err.Error())
			return err
		}

		for _, ugc := range update.UGC {
			switch update.Action {
			case "CAN":
				fallthrough
			case "UPG":
				fallthrough
			case "EXP":
				index := -1
				for i, u := range warning.UGC {
					if u == ugc {
						index = i
					}
				}

				if index > -1 {
					ret := make([]string, 0)
					ret = append(ret, warning.UGC[:index]...)
					warning.UGC = append(ret, warning.UGC[index+1:]...)
				}
			default:
				index := -1
				for i, u := range warning.UGC {
					if u == ugc {
						index = i
					}
				}

				if index == -1 {
					warning.UGC = append(warning.UGC, ugc)
				} else {

				}
			}
		}

		warning.Expires = update.Expires
		warning.Ends = update.Ends
		warning.Text = update.Text
		warning.Action = update.Action
		warning.Title = update.Title
		warning.IsEmergency = update.IsEmergency
		warning.IsPDS = update.IsPDS
		warning.Polygon = update.Polygon
		warning.Direction = update.Direction
		warning.Location = update.Location
		warning.Speed = update.Speed
		warning.SpeedText = update.SpeedText
		warning.TMLTime = update.TMLTime
		warning.Tornado = update.Tornado
		warning.Damage = update.Damage
		warning.HailThreat = update.HailThreat
		warning.HailTag = update.HailTag
		warning.WindThreat = update.WindThreat
		warning.WindTag = update.WindTag
		warning.FlashFlood = update.FlashFlood
		warning.RainfallTag = update.RainfallTag
		warning.FloodTagDam = update.FloodTagDam
		warning.SpoutTag = update.SpoutTag
		warning.SnowSquall = update.SnowSquall
		warning.SnowSquallTag = update.SnowSquallTag

		_, err := handler.db.Exec(handler.db.CTX, `
		UPDATE warning SET updated_at = $1, expires = $2, ends = $3, text = $4, action = $5, title = $6, 
    is_emergency = $7, is_pds = $8, polygon = $9, direction = $10, location = $11, speed = $12, 
    speed_text = $13, tml_time = $14, ugc = $15, tornado = $16, damage = $17, hail_threat = $18, 
    hail_tag = $19, wind_threat = $20, wind_tag = $21, flash_flood = $22, rainfall_tag = $23, 
    flood_tag_dam = $24, spout_tag = $25, snow_squall = $26, snow_squall_tag = $27
	WHERE wfo = $28 AND phenomena = $29 AND significance = $30 AND event_number = $31 AND year = $32
		`, time.Now().UTC(), warning.Expires, warning.Ends, warning.Text, warning.Action, warning.Title,
			warning.IsEmergency, warning.IsPDS, warning.Polygon, warning.Direction, warning.Location, warning.Speed,
			warning.SpeedText, warning.TMLTime, warning.UGC, warning.Tornado, warning.Damage, warning.HailThreat, warning.HailTag,
			warning.WindThreat, warning.WindTag, warning.FlashFlood, warning.RainfallTag, warning.FloodTagDam,
			warning.SpoutTag, warning.SnowSquall, warning.SnowSquallTag, warning.WFO, warning.Phenomena,
			warning.Significance, warning.EventNumber, warning.Year)
		if err != nil {
			handler.logger.Error("failed to update warning: " + err.Error())
			return err
		}
	}

	return nil
}
