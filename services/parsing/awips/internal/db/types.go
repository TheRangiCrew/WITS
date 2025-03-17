package db

import (
	"time"

	"github.com/twpayne/go-geos"
)

// State
type State struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FIPS       string `json:"fips"`
	IsOffshore bool   `json:"is_offshore"`
}

// NWS/NOAA Office
type Office struct {
	ID       string     `json:"id"`
	ICAO     string     `json:"icao"`
	Name     string     `json:"name"`
	State    string     `json:"state"`
	Location *geos.Geom `json:"location"`
}

// County Warning Area
type CWA struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Area      float32    `json:"area,omitempty"`
	Geom      *geos.Geom `json:"geom,omitempty"`
	WFO       string     `json:"wfo"`
	Region    string     `json:"region"`
	ValidFrom *time.Time `json:"valid_from"`
}

// TODO
type HVTEC struct{}

// Product Errors
type Log struct {
	ID        int        `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Time      *time.Time `json:"time"`
	Level     string     `json:"level"`
	Product   string     `json:"product,omitempty"`
	AWIPS     string     `json:"awips,omitempty"`
	WMO       string     `json:"wmo,omitempty"`
	Text      string     `json:"text,omitempty"`
	Message   string     `json:"message"`
}

type MCD struct {
	ID               int        `json:"id,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	Product          int        `json:"product"`
	Issued           *time.Time `json:"issued"`
	Expires          *time.Time `json:"expires"`
	Year             int        `json:"year"`
	Concerning       string     `json:"concerning"`
	Geom             *geos.Geom `json:"geom"`
	WatchProbability int        `json:"watch_probability"`
}
