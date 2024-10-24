package main

import (
	"fmt"
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

func (ugc *UGC) SurQL() (string, error) {
	centre, err := orbjson.NewGeometry(ugc.Centre).MarshalJSON()
	if err != nil {
		return "", err
	}

	geometry, err := orbjson.NewGeometry(ugc.Geometry).MarshalJSON()
	if err != nil {
		return "", err
	}

	cwa := "["
	for i, c := range ugc.CWA {
		cwa += fmt.Sprintf("cwa:%s", c)
		if i < len(ugc.CWA)-1 {
			cwa += ","
		}
	}
	cwa += "]"

	return fmt.Sprintf("INSERT INTO ugc {id: ugc:%s, name: \"%s\", state: %s, type: \"%s\", number: \"%s\", centre: %s, geometry: %s, cwa: %s, is_marine: %v, is_fire: %v, valid_from: %s, valid_to: %s	};\n",
		ugc.ID, ugc.Name, ugc.State, ugc.Type, ugc.Number, string(centre), string(geometry), cwa, ugc.IsMarine, ugc.IsFire, DateToString(&ugc.ValidFrom), DateToString(ugc.ValidTo)), nil
}
