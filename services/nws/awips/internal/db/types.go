package db

import (
	"fmt"
	"time"

	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type DBConfig struct {
	surrealdb.Auth
	Endpoint string
	AsRoot   bool
}

// State
type State struct {
	ID           *models.RecordID `json:"id"`
	Name         string           `json:"name"`
	Abbreviation string           `json:"abbreviation"`
	FIPS         string           `json:"fips"`
	NS           string           `json:"NS"`
}

// Offshore
type Offshore struct {
	ID           *models.RecordID `json:"id"`
	Name         string           `json:"name"`
	Abbreviation string           `json:"abbreviation"`
	FIPS         string           `json:"fips"`
}

// NWS/NOAA Office
type Office struct {
	ID       *models.RecordID      `json:"id"`
	Code     string                `json:"code"`
	ICAO     string                `json:"icao"`
	Name     string                `json:"name"`
	State    *models.RecordID      `json:"state"`
	Location *models.GeometryPoint `json:"location"`
}

// County Warning Area
type CWA struct {
	ID        *models.RecordID             `json:"id"`
	Code      string                       `json:"code"`
	Name      string                       `json:"name"`
	Centre    *models.GeometryPoint        `json:"centre"`
	Geometry  *models.GeometryMultiPolygon `json:"geometry,omitempty"`
	Area      float32                      `json:"area,omitempty"`
	WFO       *models.RecordID             `json:"wfo"`
	Region    string                       `json:"region"`
	ValidFrom *models.CustomDateTime       `json:"valid_from"`
}

// UGC
type UGC struct {
	ID        *models.RecordID             `json:"id"`
	Name      string                       `json:"name"`
	State     *models.RecordID             `json:"state"`
	Type      string                       `json:"type"`
	Number    string                       `json:"number"`
	Area      float32                      `json:"area,omitempty"`
	Centre    *models.GeometryPoint        `json:"centre"`
	Geometry  *models.GeometryMultiPolygon `json:"geometry"`
	CWA       []*models.RecordID           `json:"cwa"`
	IsMarine  bool                         `json:"is_marine"`
	IsFire    bool                         `json:"is_fire"`
	ValidFrom *models.CustomDateTime       `json:"valid_from"`
	ValidTo   *models.CustomDateTime       `json:"valid_to,omitempty"`
}

// Text Product
type Product struct {
	ID         *models.RecordID       `json:"id"`
	CreatedAt  *models.CustomDateTime `json:"created_at,omitempty"`
	Product    string                 `json:"product"`
	Issuer     *models.RecordID       `json:"issuer"`
	Text       string                 `json:"text"`
	ReceivedAt *models.CustomDateTime `json:"received_at"`
	WMO        WMO                    `json:"wmo"`
}

type ProductID struct {
	ID       string
	Issued   time.Time
	Sequence int
}

type WMO struct {
	Datatype string                `json:"datatype"`
	Issued   models.CustomDateTime `json:"issued"`
	Original string                `json:"original"`
	WFO      string                `json:"wfo"`
	BBB      string                `json:"bbb"`
}

// VTEC Phenomena
type VTECPhenomena struct {
	ID   *models.RecordID `json:"id"`
	Code string           `json:"code"`
	Name string           `json:"name"`
}

// VTEC Actions
type VTECActions struct {
	ID   *models.RecordID `json:"id"`
	Code string           `json:"code"`
	Name string           `json:"name"`
}

// VTEC Significance
type VTECSignificance struct {
	ID   *models.RecordID `json:"id"`
	Code string           `json:"code"`
	Name string           `json:"name"`
}

// VTEC Event
type VTECEvent struct {
	ID           *models.RecordID       `json:"id,omitempty"`
	CreatedAt    *models.CustomDateTime `json:"created_at,omitempty"`
	UpdatedAt    *models.CustomDateTime `json:"updated_at,omitempty"`
	Updates      int                    `json:"updates"`
	Issued       *models.CustomDateTime `json:"issued"`
	Start        *models.CustomDateTime `json:"start,omitempty"`
	Expires      *models.CustomDateTime `json:"expires"`
	End          *models.CustomDateTime `json:"end,omitempty"`
	EndInitial   *models.CustomDateTime `json:"end_initial,omitempty"`
	Phenomena    *models.RecordID       `json:"phenomena"`
	Office       *models.RecordID       `json:"office"`
	Significance *models.RecordID       `json:"significance"`
	EventNumber  int                    `json:"event_number"`
	Title        string                 `json:"title"`
	IsEmergency  bool                   `json:"is_emergency"`
	IsPDS        bool                   `json:"is_pds"`
}

type VTECEventID struct {
	EventNumber  int    `json:"event_number"`
	Phenomena    string `json:"phenomena"`
	Office       string `json:"office"`
	Significance string `json:"significance"`
	Year         int    `json:"year"`
}

func (id *VTECEventID) String() string {
	return fmt.Sprintf(`vtec_event:{event_number: %d, phenomena:'%s', office: '%s', significance: '%s', year: %d}`, id.EventNumber, id.Phenomena, id.Office, id.Significance, id.Year)
}

func (id *VTECEventID) RecordID() models.RecordID {
	return models.NewRecordID("vtec_event", id)
}

type VTECUGC struct {
	ID         *models.RecordID       `json:"id"`
	CreatedAt  *models.CustomDateTime `json:"created_at,omitempty"`
	UpdatedAt  *models.CustomDateTime `json:"updated_at,omitempty"`
	In         *models.RecordID       `json:"in"`
	Out        *models.RecordID       `json:"out"`
	Issued     *models.CustomDateTime `json:"issued"`
	Start      *models.CustomDateTime `json:"start,omitempty"`
	Expires    *models.CustomDateTime `json:"expires"`
	End        *models.CustomDateTime `json:"end,omitempty"`
	EndInitial *models.CustomDateTime `json:"end_initial,omitempty"`
	Action     *models.RecordID       `json:"action"`
	Latest     *models.RecordID       `json:"latest"`
}

type VTECUGCID struct {
	EventNumber  int    `json:"event_number"`
	Phenomena    string `json:"phenomena"`
	Office       string `json:"office"`
	Significance string `json:"significance"`
	Year         int    `json:"year"`
	UGC          string `json:"ugc"`
}

func (id *VTECUGCID) String() string {
	return fmt.Sprintf(`vtec_ugc:{event_number: %d, phenomena:'%s', office: '%s', significance: '%s', year: %d, ugc: '%s'}`, id.EventNumber, id.Phenomena, id.Office, id.Significance, id.Year, id.UGC)
}

// VTEC Event History
type VTECHistory struct {
	ID           *models.RecordID          `json:"id"`
	CreatedAt    *models.CustomDateTime    `json:"created_at,omitempty"`
	Issued       *models.CustomDateTime    `json:"issued"`
	Start        *models.CustomDateTime    `json:"start,omitempty"`
	Expires      *models.CustomDateTime    `json:"expires"`
	End          *models.CustomDateTime    `json:"end,omitempty"`
	Original     string                    `json:"original"`
	Title        string                    `json:"title"`
	Action       *models.RecordID          `json:"action"`
	Phenomena    *models.RecordID          `json:"phenomena"`
	Office       *models.RecordID          `json:"office"`
	Significance *models.RecordID          `json:"significance"`
	EventNumber  int                       `json:"event_number"`
	VTEC         VTEC                      `json:"vtec"`
	HVTEC        HVTEC                     `json:"h_vtec,omitempty"`
	IsEmergency  bool                      `json:"is_emergency"`
	IsPDS        bool                      `json:"is_pds"`
	LatLon       *LatLon                   `json:"lat_lon,omitempty"`
	Polygon      *models.GeometryPolygon   `json:"polygon,omitempty"`
	BBox         models.GeometryMultiPoint `json:"bbox,omitempty"`
	Tags         map[string]string         `json:"tags"`
	TML          *TML                      `json:"tml,omitempty"`
	Product      *models.RecordID          `json:"product"`
	UGC          []*models.RecordID        `json:"ugc"`
}

type VTECHistoryID struct {
	EventNumber  int    `json:"event_number"`
	Phenomena    string `json:"phenomena"`
	Office       string `json:"office"`
	Significance string `json:"significance"`
	Year         int    `json:"year"`
	Sequence     int    `json:"sequence"`
}

func (id *VTECHistoryID) String() string {
	return fmt.Sprintf(`vtec_history:{event_number: %d, phenomena:'%s', office: '%s', significance: '%s', year: %d, sequence: %d}`, id.EventNumber, id.Phenomena, id.Office, id.Significance, id.Year, id.Sequence)
}

func (id *VTECHistoryID) RecordID() *models.RecordID {
	recordID := models.NewRecordID("vtec_history", id)
	return &recordID
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

func (vtec *VTEC) FromAWIPS(v awips.VTEC) {
	vtec.Class = v.Class
	vtec.Action = v.Action
	vtec.WFO = v.WFO
	vtec.Phenomena = v.Phenomena
	vtec.Significance = v.Significance
	vtec.EventNumber = v.EventNumber
	vtec.Start = v.StartString
	vtec.End = v.EndString
}

// TODO
type HVTEC struct{}

type LatLon struct {
	Original string                 `json:"original"`
	Points   models.GeometryPolygon `json:"points"`
}

type TML struct {
	Direction   int                   `json:"direction"`
	Location    models.GeometryPoint  `json:"location"`
	Speed       int                   `json:"speed"`
	SpeedString string                `json:"speedString"`
	Time        models.CustomDateTime `json:"time"`
	Original    string                `json:"original"`
}

// Warning
type Warning struct {
	ID           *models.RecordID          `json:"id,omitempty"`
	CreatedAt    *models.CustomDateTime    `json:"created_at,omitempty"`
	UpdatedAt    *models.CustomDateTime    `json:"updated_at,omitempty"`
	Updates      int                       `json:"updates"`
	Issued       *models.CustomDateTime    `json:"issued"`
	Start        *models.CustomDateTime    `json:"start,omitempty"`
	Expires      *models.CustomDateTime    `json:"expires"`
	End          *models.CustomDateTime    `json:"end,omitempty"`
	Phenomena    *models.RecordID          `json:"phenomena"`
	Office       *models.RecordID          `json:"office"`
	Significance *models.RecordID          `json:"significance"`
	EventNumber  int                       `json:"event_number"`
	Title        string                    `json:"title"`
	IsEmergency  bool                      `json:"is_emergency"`
	IsPDS        bool                      `json:"is_pds"`
	Polygon      *models.GeometryPolygon   `json:"polygon,omitempty"`
	BBox         models.GeometryMultiPoint `json:"bbox,omitempty"`
}

type WarningID struct {
	EventNumber  int    `json:"event_number"`
	Phenomena    string `json:"phenomena"`
	Office       string `json:"office"`
	Significance string `json:"significance"`
	Year         int    `json:"year"`
}

func (id *WarningID) String() string {
	return fmt.Sprintf(`warning:{event_number: %d, phenomena:'%s', office: '%s', significance: '%s', year: %d}`, id.EventNumber, id.Phenomena, id.Office, id.Significance, id.Year)
}

func (id *WarningID) RecordID() models.RecordID {
	return models.RecordID{
		Table: "warning",
		ID:    id,
	}
}

type WarningUGC struct {
	ID         *models.RecordID       `json:"id"`
	CreatedAt  *models.CustomDateTime `json:"created_at,omitempty"`
	UpdatedAt  *models.CustomDateTime `json:"updated_at,omitempty"`
	In         *models.RecordID       `json:"in"`
	Out        *models.RecordID       `json:"out"`
	Issued     *models.CustomDateTime `json:"issued"`
	Start      *models.CustomDateTime `json:"start,omitempty"`
	Expires    *models.CustomDateTime `json:"expires"`
	End        *models.CustomDateTime `json:"end,omitempty"`
	EndInitial *models.CustomDateTime `json:"end_initial,omitempty"`
	Action     *models.RecordID       `json:"action"`
	Latest     *models.RecordID       `json:"latest"`
}

type WarningUGCID struct {
	EventNumber  int    `json:"event_number"`
	Phenomena    string `json:"phenomena"`
	Office       string `json:"office"`
	Significance string `json:"significance"`
	Year         int    `json:"year"`
	UGC          string `json:"ugc"`
}

func (id *WarningUGCID) String() string {
	return fmt.Sprintf(`vtec_ugc:{event_number: %d, phenomena:'%s', office: '%s', significance: '%s', year: %d, ugc: '%s'}`, id.EventNumber, id.Phenomena, id.Office, id.Significance, id.Year, id.UGC)
}

// Warning History
type WarningHistory struct {
	ID           *models.RecordID          `json:"id"`
	CreatedAt    *models.CustomDateTime    `json:"created_at,omitempty"`
	Issued       *models.CustomDateTime    `json:"issued"`
	Start        *models.CustomDateTime    `json:"start,omitempty"`
	Expires      *models.CustomDateTime    `json:"expires"`
	End          *models.CustomDateTime    `json:"end,omitempty"`
	Original     string                    `json:"original"`
	Title        string                    `json:"title"`
	Action       *models.RecordID          `json:"action"`
	Phenomena    *models.RecordID          `json:"phenomena"`
	Office       *models.RecordID          `json:"office"`
	Significance *models.RecordID          `json:"significance"`
	EventNumber  int                       `json:"event_number"`
	IsEmergency  bool                      `json:"is_emergency"`
	IsPDS        bool                      `json:"is_pds"`
	Polygon      *models.GeometryPolygon   `json:"polygon,omitempty"`
	BBox         models.GeometryMultiPoint `json:"bbox,omitempty"`
	Tags         map[string]string         `json:"tags"`
	TML          *TML                      `json:"tml,omitempty"`
	Product      *models.RecordID          `json:"product"`
	UGC          []*models.RecordID        `json:"ugc"`
}

type WarningHistoryID struct {
	EventNumber  int    `json:"event_number"`
	Phenomena    string `json:"phenomena"`
	Office       string `json:"office"`
	Significance string `json:"significance"`
	Year         int    `json:"year"`
	Sequence     int    `json:"sequence"`
}

func (id *WarningHistoryID) String() string {
	return fmt.Sprintf(`warning_history:{event_number: %d, phenomena:'%s', office: '%s', significance: '%s', year: %d, sequence: %d}`, id.EventNumber, id.Phenomena, id.Office, id.Significance, id.Year, id.Sequence)
}

func (id *WarningHistoryID) RecordID() *models.RecordID {
	recordID := models.NewRecordID("warning_history", id)
	return &recordID
}

// Product Errors
type Log struct {
	ID        *models.RecordID       `json:"id,omitempty"`
	CreatedAt *models.CustomDateTime `json:"created_at,omitempty"`
	Time      *models.CustomDateTime `json:"time"`
	Level     string                 `json:"level"`
	Product   *models.RecordID       `json:"product,omitempty"`
	AWIPS     string                 `json:"awips,omitempty"`
	WMO       string                 `json:"wmo,omitempty"`
	Text      string                 `json:"text,omitempty"`
	Message   string                 `json:"message"`
}

type MCDID struct {
	Number int `json:"number"`
	Year   int `json:"year"`
}

type MCD struct {
	ID               *models.RecordID        `json:"id,omitempty"`
	CreatedAt        *models.CustomDateTime  `json:"created_at,omitempty"`
	UpdatedAt        *models.CustomDateTime  `json:"updated_at,omitempty"`
	Original         *models.RecordID        `json:"original"`
	Issued           *models.CustomDateTime  `json:"issued"`
	Expires          *models.CustomDateTime  `json:"expires"`
	Concerning       string                  `json:"concerning"`
	Polygon          *models.GeometryPolygon `json:"polygon"`
	WatchProbability int                     `json:"watch_probability"`
}

type Watch struct {
	ID        *models.RecordID        `json:"id,omitempty"`
	CreatedAt *models.CustomDateTime  `json:"created_at,omitempty"`
	UpdatedAt *models.CustomDateTime  `json:"updated_at,omitempty"`
	Issued    *models.CustomDateTime  `json:"issued"`
	Expires   *models.CustomDateTime  `json:"expires"`
	SAW       *models.RecordID        `json:"saw"`
	SEL       *models.RecordID        `json:"sel"`
	WOU       *models.RecordID        `json:"wou"`
	WWP       WatchWWP                `json:"wwp"`
	IsPDS     bool                    `json:"is_pds"`
	Polygon   *models.GeometryPolygon `json:"polygon"`
}

type WatchID struct {
	Number    int    `json:"number"`
	Phenomena string `json:"phenomena"`
	Year      int    `json:"year"`
}

type WatchWWP struct {
	Text             *models.RecordID `json:"text"`
	Degrees          int              `json:"degrees"`
	Speed            int              `json:"speed"`
	MaxTops          int              `json:"max_tops"`
	MaxHail          float64          `json:"max_hail"`
	OneOrMoreHail    string           `json:"one_or_more_hail"`
	TenOrMoreSevHail string           `json:"ten_or_more_sev_hail"`
	MaxWind          int              `json:"max_wind"`
	OneOrMoreWind    string           `json:"one_or_more_wind"`
	TenOrMoreSevWind string           `json:"ten_or_more_sev_wind"`
	TwoOrMoreTor     string           `json:"two_or_more_tor"`
	StrongTor        string           `json:"strong_tor"`
	SixOrMoreCombo   string           `json:"six_or_more_combo"`
}
