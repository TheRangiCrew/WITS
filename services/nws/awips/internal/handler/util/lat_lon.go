package util

import (
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
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
