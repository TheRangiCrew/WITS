package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
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

type WarningMessage struct {
	ID            int             `json:"id"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	Issued        time.Time       `json:"issued"`
	Starts        time.Time       `json:"starts"`
	Expires       time.Time       `json:"expires"`
	Ends          time.Time       `json:"ends"`
	EndInitial    time.Time       `json:"endInitial"`
	Text          string          `json:"text"`
	WFO           string          `json:"wfo"`
	Action        string          `json:"action"`
	Class         string          `json:"class"`
	Phenomena     string          `json:"phenomena"`
	Significance  string          `json:"significance"`
	EventNumber   int             `json:"eventNumber"`
	Year          int             `json:"year"`
	Title         string          `json:"title"`
	IsEmergency   bool            `json:"isEmergency"`
	IsPDS         bool            `json:"isPDS"`
	Polygon       json.RawMessage `json:"polygon,omitempty"`
	Direction     *int            `json:"direction"`
	Locations     json.RawMessage `json:"locations"`
	Speed         *int            `json:"speed"`
	SpeedText     *string         `json:"speed_text"`
	TMLTime       *time.Time      `json:"tmlTime"`
	UGC           []string        `json:"ugc"`
	Tornado       string          `json:"tornado"`
	Damage        string          `json:"damage"`
	HailThreat    string          `json:"hailThreat"`
	HailTag       string          `json:"hailTag"`
	WindThreat    string          `json:"windThreat"`
	WindTag       string          `json:"windTag"`
	FlashFlood    string          `json:"flashFlood"`
	RainfallTag   string          `json:"rainfallTag"`
	FloodTagDam   string          `json:"floodTagDam"`
	SpoutTag      string          `json:"spoutTag"`
	SnowSquall    string          `json:"snowSquall"`
	SnowSquallTag string          `json:"snowSquallTag"`
}

func (handler *vtecHandler) warning(update VTECUpdate) error {
	// If the warning ends more than 36 hours ago, ignore it since it may be old data
	if update.Ends.Before(time.Now().Add(-(time.Hour * 36))) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Lets check if the Warning is already in the database
	rows, err := handler.DB.Query(ctx, `
			SELECT * FROM vtec.warnings WHERE
			wfo = $1 AND phenomena = $2 AND significance = $3 AND event_number = $4 AND year = $5
			`, update.WFO, update.Phenomena, update.Significance, update.EventNumber, update.Year)
	if err != nil {
		handler.logger.Error("failed to get warning: " + err.Error())
		return err
	}
	defer rows.Close()

	var warning *Warning

	if !rows.Next() {
		if rows.Err() != nil {
			handler.logger.Error("failed to get warning: " + rows.Err().Error())
			return err
		}

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

		if update.Action != "CAN" && update.Action != "EXP" && update.Action != "UPG" {

			_, err := handler.DB.CopyFrom(
				context.Background(),
				pgx.Identifier{"vtec", "warnings"},
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

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := handler.DB.Exec(ctx, `
		UPDATE vtec.warnings SET updated_at = $1, expires = $2, ends = $3, text = $4, action = $5, title = $6, 
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

	var polygon json.RawMessage = nil
	if warning.Polygon != nil {
		polygonGeoJSON := warning.Polygon.ToGeoJSON(1)
		polygon = json.RawMessage(polygonGeoJSON)
	}
	var location json.RawMessage = nil
	if warning.Location != nil {
		locationGeoJSON := warning.Location.ToGeoJSON(1)
		location = json.RawMessage(locationGeoJSON)
	}

	warningMessage := WarningMessage{
		ID:            warning.ID,
		CreatedAt:     warning.CreatedAt,
		UpdatedAt:     warning.UpdatedAt,
		Issued:        warning.Issued,
		Starts:        warning.Starts,
		Expires:       warning.Expires,
		Ends:          warning.Ends,
		EndInitial:    warning.EndInitial,
		Text:          warning.Text,
		WFO:           warning.WFO,
		Action:        warning.Action,
		Class:         warning.Class,
		Phenomena:     warning.Phenomena,
		Significance:  warning.Significance,
		EventNumber:   warning.EventNumber,
		Year:          warning.Year,
		Title:         warning.Title,
		IsEmergency:   warning.IsEmergency,
		IsPDS:         warning.IsPDS,
		Polygon:       polygon,
		Direction:     warning.Direction,
		Locations:     location,
		Speed:         warning.Speed,
		SpeedText:     warning.SpeedText,
		TMLTime:       warning.TMLTime,
		UGC:           warning.UGC,
		Tornado:       warning.Tornado,
		Damage:        warning.Damage,
		HailThreat:    warning.HailThreat,
		HailTag:       warning.HailTag,
		WindThreat:    warning.WindThreat,
		WindTag:       warning.WindTag,
		FlashFlood:    warning.FlashFlood,
		RainfallTag:   warning.RainfallTag,
		FloodTagDam:   warning.FloodTagDam,
		SpoutTag:      warning.SpoutTag,
		SnowSquall:    warning.SnowSquall,
		SnowSquallTag: warning.SnowSquallTag,
	}

	message, err := json.Marshal(warningMessage)
	if err != nil {
		handler.logger.Error("failed to marshal warning message: " + err.Error())
		return err
	}

	publisher, err := handler.GetPublisher(RealtimeExchange)
	if err != nil {
		handler.logger.Error("failed to get publisher: "+err.Error(), "publisher", RealtimeExchange)
		return err
	}
	err = publisher.PublishWithContext(context.Background(), RealtimeExchange, "realtime.warnings", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	if err != nil {
		handler.logger.Error("failed to publish warning message: " + err.Error())
	}

	return err
}
