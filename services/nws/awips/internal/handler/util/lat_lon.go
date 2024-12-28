package util

import (
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/paulmach/orb"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func LatLonFromAwips(src awips.LatLon) db.LatLon {
	output := db.LatLon{
		Original: src.Original,
		Points:   models.GeometryPolygon{},
	}
	line := models.GeometryLine{}
	for _, i := range src.Points {
		p := models.GeometryPoint{
			Latitude:  i[0],
			Longitude: i[1],
		}
		line = append(line, p)
	}
	output.Points = append(output.Points, line)

	return output
}

func BBoxFromAwips(src awips.LatLon) models.GeometryMultiPoint {
	ring := orb.Ring{}
	for _, i := range src.Points {
		var p orb.Point = [2]float64{i[0], i[1]}

		ring = append(ring, p)
	}
	var polygon orb.Polygon = []orb.Ring{ring}
	bound := polygon.Bound()
	min := bound.Min
	max := bound.Max
	return []models.GeometryPoint{models.NewGeometryPoint(min.Lon(), min.Lat()), models.NewGeometryPoint(max.Lon(), max.Lat())}
}

func BBoxFromMultiPolygon(src []models.GeometryMultiPolygon) models.GeometryMultiPoint {
	c := orb.Collection{}
	for _, i := range src {
		mp := orb.MultiPolygon{}
		for _, j := range i {
			p := orb.Polygon{}
			for _, k := range j {
				r := orb.Ring{}
				for _, l := range k {
					r = append(r, l.GetCoordinates())
				}
				p = append(p, r)
			}
			mp = append(mp, p)
		}
		c = append(c, mp)
	}
	bound := c.Bound()
	min := bound.Min
	max := bound.Max
	return []models.GeometryPoint{models.NewGeometryPoint(min.Lon(), min.Lat()), models.NewGeometryPoint(max.Lon(), max.Lat())}
}

func PolygonFromAwips(src awips.PolygonFeature) models.GeometryPolygon {
	line := models.GeometryLine{}
	for _, i := range src.Coordinates[0] {
		p := models.GeometryPoint{
			Latitude:  i[0],
			Longitude: i[1],
		}
		line = append(line, p)
	}

	return []models.GeometryLine{line}
}
