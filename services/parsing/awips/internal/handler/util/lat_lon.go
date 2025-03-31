package util

import (
	"sync"

	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/twpayne/go-geos"
)

var geosMutex sync.Mutex

func PolygonFromAwips(src awips.PolygonFeature) *geos.Geom {
	geosMutex.Lock()
	defer geosMutex.Unlock()

	geom := geos.NewPolygon(src.Coordinates)

	geom.SetSRID(4326)
	return geom
}
