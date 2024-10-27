package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

var floatType = reflect.TypeOf(float64(0))

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		PrintHelp()
		return
	}

	if args[0] == "-d" {
		err := ParseAll(args[1])
		if err != nil {
			slog.Error(err.Error())
		}
	} else if args[0] == "-f" {
		err := Parse(args[1])
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func ParseAll(directory string) error {

	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	fileRegexp := regexp.MustCompile("([a-z_]{2}[0-9]{2}){2}(.zip)")

	for _, d := range dir {
		if d.IsDir() {
			continue
		}

		if fileRegexp.MatchString(d.Name()) {
			err = Parse(d.Name())
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}

	return nil
}

func Parse(filename string) error {

	if len(filename) != 12 {
		return fmt.Errorf("filename has invalid length. Expected 12 characters (e.g. c_05mr24.zip) but found %d", len(filename))
	}

	// The type of the shapes in the shapefile
	thing := filename[:2]

	// The day the shapefile was generated
	day, err := strconv.Atoi(filename[2:4])
	if err != nil {
		return err
	}

	// The month the shapefile was generated
	monthString := filename[4:6]
	var month time.Month
	// NOTE: This is kind of just added as new months are found.
	switch monthString {
	case "mr":
		month = time.March
	case "my":
		month = time.May
	case "jn":
		month = time.June
	case "au":
		month = time.August
	case "se":
		month = time.September
	case "oc":
		month = time.October
	default:
		return fmt.Errorf("could not find month for %s in %s", monthString, filename)
	}

	year, err := strconv.Atoi(filename[6:8])
	if err != nil {
		return err
	}

	t := time.Date(2000+year, month, day, 0, 0, 0, 0, time.UTC)

	scanner, err := CreateZipScanner(filename)
	if err != nil {
		return err
	}

	switch thing {
	case "c_":
		err = ParseCounties(scanner, t)
	case "z_":
		err = ParseZones(scanner, t)
	case "mz":
		err = ParseMarineZones(scanner, t)
	case "fz":
		err = ParseFire(scanner, t)
	case "w_":
		err = ParseCWA(scanner, t)
	}

	return err
}

func ParseWFO(filename string) error {
	return nil
}

func PrintHelp() {
	fmt.Println("\nHow to use this thing:")
	fmt.Println("-d <path>		Attempt to parse supported shapefiles from this directory")
	fmt.Println("-f <filename>	Attempt to parse the specified shapefile")

}
