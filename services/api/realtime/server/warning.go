package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/TheRangiCrew/WITS/services/api/realtime/util"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/twpayne/go-geos"
)

type WarningDB struct {
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

type Warning struct {
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

func (warning *Warning) StringID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%d", warning.WFO, warning.Phenomena, warning.Significance, util.PadZero(strconv.Itoa(warning.EventNumber), 4), warning.Year)
}

func (server *Server) InitialiseWarningData() error {
	rows, err := server.DB.Query(context.Background(), "SELECT * FROM vtec.warnings WHERE ends > CURRENT_TIMESTAMP AND (action != 'CAN' OR action != 'EXP') ORDER BY created_at ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	now := time.Now()

	for rows.Next() {
		var warningDB WarningDB
		err := rows.Scan(
			&warningDB.ID,
			&warningDB.CreatedAt,
			&warningDB.UpdatedAt,
			&warningDB.Issued,
			&warningDB.Starts,
			&warningDB.Expires,
			&warningDB.Ends,
			&warningDB.EndInitial,
			&warningDB.Text,
			&warningDB.WFO,
			&warningDB.Action,
			&warningDB.Class,
			&warningDB.Phenomena,
			&warningDB.Significance,
			&warningDB.EventNumber,
			&warningDB.Year,
			&warningDB.Title,
			&warningDB.IsEmergency,
			&warningDB.IsPDS,
			&warningDB.Polygon,
			&warningDB.Direction,
			&warningDB.Location,
			&warningDB.Speed,
			&warningDB.SpeedText,
			&warningDB.TMLTime,
			&warningDB.UGC,
			&warningDB.Tornado,
			&warningDB.Damage,
			&warningDB.HailThreat,
			&warningDB.HailTag,
			&warningDB.WindThreat,
			&warningDB.WindTag,
			&warningDB.FlashFlood,
			&warningDB.RainfallTag,
			&warningDB.FloodTagDam,
			&warningDB.SpoutTag,
			&warningDB.SnowSquall,
			&warningDB.SnowSquallTag)
		if err != nil {
			slog.Error("Failed to scan warning data", "error", err)
			continue
		}

		var polygon json.RawMessage = nil
		if warningDB.Polygon != nil {
			polygonGeoJSON := warningDB.Polygon.ToGeoJSON(1)
			polygon = json.RawMessage(polygonGeoJSON)
		}
		var location json.RawMessage = nil
		if warningDB.Location != nil {
			locationGeoJSON := warningDB.Location.ToGeoJSON(1)
			location = json.RawMessage(locationGeoJSON)
		}

		warning := Warning{
			ID:            warningDB.ID,
			CreatedAt:     warningDB.CreatedAt,
			UpdatedAt:     warningDB.UpdatedAt,
			Issued:        warningDB.Issued,
			Starts:        warningDB.Starts,
			Expires:       warningDB.Expires,
			Ends:          warningDB.Ends,
			EndInitial:    warningDB.EndInitial,
			Text:          warningDB.Text,
			WFO:           warningDB.WFO,
			Action:        warningDB.Action,
			Class:         warningDB.Class,
			Phenomena:     warningDB.Phenomena,
			Significance:  warningDB.Significance,
			EventNumber:   warningDB.EventNumber,
			Year:          warningDB.Year,
			Title:         warningDB.Title,
			IsEmergency:   warningDB.IsEmergency,
			IsPDS:         warningDB.IsPDS,
			Polygon:       polygon,
			Direction:     warningDB.Direction,
			Locations:     location,
			Speed:         warningDB.Speed,
			SpeedText:     warningDB.SpeedText,
			TMLTime:       warningDB.TMLTime,
			UGC:           warningDB.UGC,
			Tornado:       warningDB.Tornado,
			Damage:        warningDB.Damage,
			HailThreat:    warningDB.HailThreat,
			HailTag:       warningDB.HailTag,
			WindThreat:    warningDB.WindThreat,
			WindTag:       warningDB.WindTag,
			FlashFlood:    warningDB.FlashFlood,
			RainfallTag:   warningDB.RainfallTag,
			FloodTagDam:   warningDB.FloodTagDam,
			SpoutTag:      warningDB.SpoutTag,
			SnowSquall:    warningDB.SnowSquall,
			SnowSquallTag: warningDB.SnowSquallTag,
		}

		bytes, err := json.Marshal(warning)
		if err != nil {
			slog.Error("Failed to marshal warning data to JSON", "error", err)
			continue
		}

		stat := server.RDB.Set(context.Background(), "warning:"+warning.StringID(), bytes, warning.Ends.Sub(now))
		if stat.Err() != nil {
			slog.Error("Failed to set warning data in Redis", "error", stat.Err())
			continue
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// Setup API routes
	server.Router.GET("/warnings/current", func(c *gin.Context) {
		mcds, err := server.getAllWarnings()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve warning data: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, mcds)
	})

	// Setup Websocket
	hub := newHub()
	server.Hubs["warnings"] = hub
	server.Router.GET("/warnings/live", server.warningSocketHandler)

	return nil
}

func (server *Server) getAllWarnings() ([]Warning, error) {
	ctx := context.Background()
	iter := server.RDB.Scan(ctx, 0, "warning:*", 0).Iterator()
	warnings := []Warning{}
	for iter.Next(ctx) {
		res := server.RDB.Get(ctx, iter.Val())
		var warning Warning
		err := json.Unmarshal([]byte(res.Val()), &warning)
		if err != nil {
			slog.Error("Failed to unmarshal Warning data", "error", err, "key", iter.Val())
			return nil, err
		}
		warnings = append(warnings, warning)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return warnings, nil
}

func (server *Server) warningSocketHandler(c *gin.Context) {
	hub, ok := server.Hubs["warnings"]
	if !ok {
		slog.Error("Failed to find warnings hub")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade warnings socket connection", "error", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	go client.writePump()

	// Read messages only to detect disconnect
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
	hub.unregister <- client
}

func (server *Server) handleWarning(msg amqp.Delivery) error {

	var warning Warning
	err := json.Unmarshal(msg.Body, &warning)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal warning message: %s", err.Error())
	}

	if warning.Action == "CAN" || warning.Action == "EXP" || warning.Action == "UPG" {
		stat := server.RDB.Del(context.Background(), "warning:"+warning.StringID())
		if stat.Err() != nil {
			return stat.Err()
		}
	} else {
		stat := server.RDB.Set(context.Background(), "warning:"+warning.StringID(), msg.Body, warning.Ends.Sub(time.Now()))
		if stat.Err() != nil {
			return stat.Err()
		}
	}

	fmt.Printf("%s %s %s %s\n", warning.WFO, warning.Action, warning.Title, util.PadZero(strconv.Itoa(warning.EventNumber), 4))

	hub, ok := server.Hubs["warnings"]
	if !ok {
		return fmt.Errorf("Could not get warnings socket.")
	}

	hub.broadcast <- msg.Body

	return nil
}
