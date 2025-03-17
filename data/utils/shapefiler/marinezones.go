package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/everystreet/go-shapefile"
	orbjson "github.com/paulmach/orb/geojson"
)

func ParseMarineZones(scanner *shapefile.ZipScanner, t time.Time) error {

	slog.Info("Parsing marine zones...")

	// Start the scanner
	err := scanner.Scan()
	if err != nil {
		return err
	}

	info, err := scanner.Info()
	if err != nil {
		return err
	}

	ugcRecords := make([]UGC, info.NumRecords)
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

		idAttr, _ := record.Attributes.Field("ID")
		id := fmt.Sprintf("%v", idAttr.Value())
		state := id[0:2]
		number := id[3:]

		zonename, _ := record.Attributes.Field("NAME")
		name := fmt.Sprintf("%v", zonename.Value())

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

		centre := [2]float64{lon, lat}

		cwaAttr, _ := record.Attributes.Field("WFO")
		cwa := fmt.Sprintf("%v", cwaAttr.Value())
		cwaArr := make([]string, len(cwa)/3)

		for i := 0; i < len(cwa); i += 3 {
			cwaArr[i/3] = cwa[i : i+3]
		}

		ugc := UGC{
			ID:        id,
			Name:      name,
			State:     state,
			Number:    number,
			Type:      "Z",
			Area:      0.0, // Will be calculated in DB
			Centre:    centre,
			Geometry:  *mpolygon,
			CWA:       cwaArr,
			IsMarine:  true,
			IsFire:    false,
			ValidFrom: t,
			ValidTo:   nil,
		}

		ugcRecords[count] = ugc

		count++

	}

	// Err() returns the first error encountered during calls to Record()
	err = scanner.Err()
	if err != nil {
		return err
	}

	out, err := ToSQL(ugcRecords)
	if err != nil {
		return err
	}

	err = WriteToFile("marinezones.sql", []byte(out))
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Wrote %d records to marinezones.sql\n", len(ugcRecords)))

	collection := orbjson.NewFeatureCollection()

	for _, ugc := range ugcRecords {
		feature := orbjson.NewFeature(ugc.Geometry)
		feature.Properties = map[string]interface{}{
			"id":        ugc.ID,
			"name":      ugc.Name,
			"state":     ugc.State,
			"type":      ugc.Type,
			"number":    ugc.Number,
			"is_marine": ugc.IsMarine,
			"is_fire":   ugc.IsFire,
			"cwa":       ugc.CWA,
		}
		collection.Append(feature)
	}

	data, err := collection.MarshalJSON()
	if err != nil {
		return err
	}

	WriteToFile("marinezones.geojson", data)

	slog.Info(fmt.Sprintf("Wrote %d records to marinezones.geojson\n", len(ugcRecords)))

	return err
}
