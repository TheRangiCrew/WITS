package models

import (
	"time"

	"github.com/twpayne/go-geos"
)

type VTECEvent struct {
	ID           int       `json:"id,omitempty"`
	EventID      string    `json:"event_id"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	Issued       time.Time `json:"issued"`
	Starts       time.Time `json:"starts,omitempty"`
	Expires      time.Time `json:"expires"`
	Ends         time.Time `json:"ends,omitempty"`
	EndInitial   time.Time `json:"endInitial,omitempty"`
	Class        string    `json:"class"`
	Phenomena    string    `json:"phenomena"`
	WFO          string    `json:"wfo"`
	Significance string    `json:"significance"`
	EventNumber  int       `json:"eventNumber"`
	Year         int       `json:"year"`
	Title        string    `json:"title"`
	IsEmergency  bool      `json:"isEmergency"`
	IsPDS        bool      `json:"isPDS"`
	PolygonStart *geos.Geom
}

func GetVTEC() {

}
