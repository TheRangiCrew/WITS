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
	Code      string       `json:"code"`
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
			ID:        fmt.Sprintf("cwa:%s", id),
			Code:      id,
			Name:      name,
			Centre:    location,
			Geometry:  *mpolygon,
			Area:      0.0,
			WFO:       fmt.Sprintf("office:%s", id),
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

	out := ""

	for _, cwa := range cwaRecords {
		centre, err := orbjson.NewGeometry(cwa.Centre).MarshalJSON()
		if err != nil {
			return err
		}

		geometry, err := orbjson.NewGeometry(cwa.Geometry).MarshalJSON()
		if err != nil {
			return err
		}

		out += fmt.Sprintf("INSERT INTO cwa { id: %s, code: \"%s\", name: \"%s\", centre: %s, geometry: %s, wfo: %s, region: \"%s\", valid_from: %s};\n", cwa.ID, cwa.Code, cwa.Name, string(centre), string(geometry), cwa.WFO, cwa.Region, DateToString(&cwa.ValidFrom))
	}

	err = WriteToFile("cwa.surql", []byte(out))
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Wrote %d records to cwa.surql\n", len(cwaRecords)))

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
