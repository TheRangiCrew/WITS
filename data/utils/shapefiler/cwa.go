package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/everystreet/go-shapefile"
	"github.com/paulmach/orb"
	orbjson "github.com/paulmach/orb/geojson"
)

type CWA struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Centre    orb.Point    `json:"centre"`
	Geometry  orb.Geometry `json:"geometry"`
	Area      float64      `json:"area"`
	WFO       string       `json:"wfo"`
	Region    string       `json:"region"`
	ValidFrom time.Time    `json:"valid_from"`
}

func ParseCWA(scanner *shapefile.ZipScanner, t time.Time) error {

	// Start the scanner
	err := scanner.Scan()
	if err != nil {
		return err
	}

	info, err := scanner.Info()
	if err != nil {
		return err
	}

	cwaRecords := make([]CWA, info.NumRecords)
	count := 0

	for {
		record := scanner.Record()
		if record == nil {
			break
		}

		shape := record.Shape.GeoJSONFeature()
		mpolygon, err := GetShape(shape)
		if err != nil {
			return err
		}

		cwaAttr, _ := record.Attributes.Field("CWA")
		id := fmt.Sprintf("%v", cwaAttr.Value())

		cwaName, _ := record.Attributes.Field("CITY")
		name := fmt.Sprintf("%v", cwaName.Value())

		lonAttr, _ := record.Attributes.Field("LON")
		lon, err := getFloat(lonAttr.Value())
		if err != nil {
			return err
		}

		latAttr, _ := record.Attributes.Field("LAT")
		lat, err := getFloat(latAttr.Value())
		if err != nil {
			return err
		}

		location := [2]float64{lon, lat}

		regionAttr, _ := record.Attributes.Field("REGION")
		region := fmt.Sprintf("%v", regionAttr.Value())

		cwa := CWA{
			ID:        id,
			Name:      name,
			Centre:    location,
			Geometry:  *mpolygon,
			Area:      0.0,
			WFO:       id,
			Region:    region,
			ValidFrom: t,
		}

		cwaRecords[count] = cwa

		count++

	}

	// Err() returns the first error encountered during calls to Record()
	err = scanner.Err()
	if err != nil {
		return err
	}

	result := "INSERT INTO postgis.cwas(id, name, area, geom, wfo, region, valid_from) VALUES\n"

	for _, cwa := range cwaRecords {
		geometry, err := orbjson.NewGeometry(cwa.Geometry).MarshalJSON()
		if err != nil {
			return err
		}

		result += fmt.Sprintf("('%s', '%s', ST_Area(ST_GeomFromGeoJSON('%s')), ST_GeomFromGeoJSON('%s'), '%s', '%s', %v),\n",
			cwa.ID, cwa.Name, geometry, geometry, cwa.WFO, cwa.Region, DateToString(&cwa.ValidFrom))
	}

	result = result[:len(result)-2]

	err = WriteToFile("cwa.sql", []byte(result))
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Wrote %d records to cwa.sql\n", len(cwaRecords)))

	collection := orbjson.NewFeatureCollection()

	for _, cwa := range cwaRecords {
		feature := orbjson.NewFeature(cwa.Geometry)
		feature.Properties = map[string]interface{}{
			"id":     cwa.ID,
			"name":   cwa.Name,
			"wfo":    cwa.WFO,
			"region": cwa.Region,
		}
		collection.Append(feature)
	}

	data, err := collection.MarshalJSON()
	if err != nil {
		return err
	}

	WriteToFile("cwa.geojson", data)

	slog.Info(fmt.Sprintf("Wrote %d records to cwa.geojson\n", len(cwaRecords)))

	return err
}
