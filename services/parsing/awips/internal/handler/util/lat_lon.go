package util

import (
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/twpayne/go-geos"
)

func PolygonFromAwips(src awips.PolygonFeature) geos.Geom {
	return *geos.NewPolygon(src.Coordinates).SetSRID(4326)
}
