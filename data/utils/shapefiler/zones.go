package main

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/everystreet/go-shapefile"
	orbjson "github.com/paulmach/orb/geojson"
)

func ParseZones(scanner *shapefile.ZipScanner, t time.Time) error {

	slog.Info("Parsing zones...")

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

		zoneAttr, _ := record.Attributes.Field("ZONE")
		zone := fmt.Sprintf("%v", zoneAttr.Value())
		if len(zone) != 3 {
			return errors.New("fips did not have the correct length")
		}

		stateAttr, _ := record.Attributes.Field("STATE")
		state := fmt.Sprintf("%v", stateAttr.Value())

		zonename, _ := record.Attributes.Field("SHORTNAME")
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

		cwaAttr, _ := record.Attributes.Field("CWA")
		cwa := fmt.Sprintf("%v", cwaAttr.Value())
		cwaArr := make([]string, len(cwa)/3)

		for i := 0; i < len(cwa); i += 3 {
			cwaArr[i/3] = cwa[i : i+3]
		}

		ugc := UGC{
			ID:        state + "Z" + zone,
			Name:      name,
			State:     state,
			Number:    zone,
			Type:      "Z",
			Area:      0.0,
			Centre:    centre,
			Geometry:  *mpolygon,
			CWA:       cwaArr,
			IsMarine:  false,
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

	err = WriteToFile("zones.sql", []byte(out))
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Wrote %d records to zones.sql\n", len(ugcRecords)))

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

	WriteToFile("zones.geojson", data)

	slog.Info(fmt.Sprintf("Wrote %d records to zones.geojson\n", len(ugcRecords)))

	return nil
}
