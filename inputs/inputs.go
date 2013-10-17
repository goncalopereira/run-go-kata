package inputs

import (
	"encoding/csv"
	"run/data"
	"os"
	"strconv"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

func ReadFiles(names []string, stationsFilename string) (locations map[string][]data.Location, stations []data.Tube, err error) {
	stations, err = ReadStationsFile(stationsFilename)
	locations, err = ReadAllRobotFiles(names)
	return
}

func ReadStationsFile(filename string) (stations []data.Tube, err error) {
	file, err := os.Open("./files/" + filename + ".csv")

	if err != nil {
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	if err != nil {
		return
	}

	stations = make([]data.Tube, len(lines))

	for i := 0; i < len(lines); i++ {
		lat, err := strconv.ParseFloat(lines[i][1], 64)
		long, err := strconv.ParseFloat(lines[i][2], 64)

		if err == nil {
			stations[i] = data.Tube{Latitude: lat, Longitude: long}
		}
	}

	return
}

func ReadAllRobotFiles(names []string) (locations map[string][]data.Location, err error) {
	locations = make(map[string][]data.Location)

	for _, name := range names {
		locations[name], err = ReadRobotFile(name)

		if err != nil {
			return
		}
	}

	return
}

func ReadRobotFile(name string) (locations []data.Location, err error) {
	file, err := os.Open("./files/" + name + ".csv")

	if err != nil {
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	if err != nil {
		return
	}

	locations = make([]data.Location, len(lines))

	for i := 0; i < len(lines); i++ {
		lat, err := strconv.ParseFloat(lines[i][1], 64)
		long, err := strconv.ParseFloat(lines[i][2], 64)

		t, err := time.Parse(TIME_FORMAT, lines[i][3])

		if err == nil {
			locations[i] = data.Location{Latitude: lat, Longitude: long, Time: t}
		}

	}

	return
}
