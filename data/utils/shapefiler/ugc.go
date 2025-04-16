package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/paulmach/orb"
	orbjson "github.com/paulmach/orb/geojson"
)

type UGC struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	State     string       `json:"state"`
	Number    string       `json:"number"`
	Type      string       `json:"type"`
	Area      float64      `json:"area"`
	Centre    orb.Point    `json:"centre"`
	Geometry  orb.Geometry `json:"geometry"`
	CWA       []string     `json:"cwa"`
	IsMarine  bool         `json:"is_marine"`
	IsFire    bool         `json:"is_fire"`
	ValidFrom time.Time    `json:"valid_from"`
	ValidTo   *time.Time   `json:"valid_to"`
}

func ToSQL(ugcs []UGC) (string, error) {
	result := "INSERT INTO postgis.ugcs(ugc, name, state, type, number, area, geom, cwa, is_marine, is_fire, valid_from) VALUES \n"

	for _, ugc := range ugcs {

		geometry, err := orbjson.NewGeometry(ugc.Geometry).MarshalJSON()
		if err != nil {
			return "", err
		}

		ugc.Name = strings.ReplaceAll(ugc.Name, "'", "''")

		cwa := "'{"
		for i, c := range ugc.CWA {
			cwa += fmt.Sprintf("\"%s\"", c)
			if i < len(ugc.CWA)-1 {
				cwa += ","
			}
		}
		cwa += "}'"

		result += fmt.Sprintf("('%s', '%s', '%s', '%s', %s, ST_Area(ST_GeomFromGeoJSON('%s')), ST_GeomFromGeoJSON('%s'), %s, %v, %v, %s),\n",
			ugc.ID, ugc.Name, ugc.State, ugc.Type, ugc.Number, string(geometry), string(geometry), cwa, ugc.IsMarine, ugc.IsFire, DateToString(&ugc.ValidFrom))
	}

	result = result[:len(result)-2]

	return result, nil
}
