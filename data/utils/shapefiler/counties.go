package main

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/everystreet/go-shapefile"
	orbjson "github.com/paulmach/orb/geojson"
)

func ParseCounties(scanner *shapefile.ZipScanner, t time.Time) error {

	slog.Info("Parsing counties...")

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

	// Call Record() to get each record in turn, until either the end of the file, or an error occurs
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

		fips, _ := record.Attributes.Field("FIPS")
		id := fmt.Sprintf("%v", fips.Value())[2:]
		if len(id) != 3 {
			return errors.New("fips did not have the correct length")
		}

		stateAttr, _ := record.Attributes.Field("STATE")
		state := fmt.Sprintf("%v", stateAttr.Value())

		ugcID := state + "C" + id

		countyname, _ := record.Attributes.Field("COUNTYNAME")
		name := fmt.Sprintf("%v", countyname.Value())

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
		cwaarr := make([]string, len(cwa)/3)

		for i := 0; i < len(cwa); i += 3 {
			cwaarr[i/3] = cwa[i : i+3]
		}

		ugc := UGC{
			ID:        ugcID,
			Name:      name,
			State:     state,
			Number:    id,
			Type:      "C",
			Area:      0.0,
			Centre:    centre,
			Geometry:  *mpolygon,
			CWA:       cwaarr,
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

	err = WriteToFile("counties.sql", []byte(out))
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Wrote %d records to counties.sql\n", len(ugcRecords)))

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

	WriteToFile("counties.geojson", data)

	slog.Info(fmt.Sprintf("Wrote %d records to counties.geojson\n", len(ugcRecords)))

	return err
}
