package movement

import (
	"run/data"
	"math"
	"testing"
	"time"
)

func TestGivenSameLocationTwiceReturnZeroDistance(t *testing.T) {
	var lat float64 = 123
	var long float64 = 456

	distance := Distance(lat, long, lat, long)

	if distance != 0 {
		t.Errorf("Same location should have 0 distance but it was %f", distance)
	}
}

func TestGivenSameTimeTwiceReturnZeroSpeed(t *testing.T) {
	n := time.Now()

	speed := Speed(123, n, n)

	if speed != 0 {
		t.Errorf("Same time should have 0 speed but it was %f", speed)
	}
}

func TestGivenTwoKnownLocationsReturnKnownDistance(t *testing.T) {
	var fromLat float64 = 123
	var fromLong float64 = 456
	var toLat float64 = 124
	var toLong float64 = 456
	var knownResult float64 = 111194

	distance := Distance(fromLat, fromLong, toLat, toLong)

	if math.Floor(distance) != knownResult {
		t.Errorf("Known distance in km should be aprox 111194 but it was %f", distance)
	}
}

func BenchmarkDistance(t *testing.B) {
	var fromLat float64 = 123
	var fromLong float64 = 456
	var toLat float64 = 124
	var toLong float64 = 456

	for i := 0; i < t.N; i++ {
		Distance(fromLat, fromLong, toLat, toLong)
	}
}

func TestGivenValidTimeAndDistanceReturnCorrectSpeed(t *testing.T) {
	fromTime := time.Now()

	toTime := fromTime.Add(5 * time.Second)

	meters := 25.0

	speed := Speed(meters, fromTime, toTime)

	if math.Floor(speed) != 5 {

		t.Errorf("Was expected 5m/s speed and got %f", math.Floor(speed))
	}
}

func TestGivenNoStationsCloseReturnsFalse(t *testing.T) {
	stations := make([]data.Tube, 0)
	location := data.Location{}

	result := CloseToStation(location, stations)

	if result != false {
		t.Error("Empty station list should return false")
	}
}

func TestGivenStationLess350MetersReturnTrue(t *testing.T) {

	location := data.Location{}
	stations := []data.Tube{data.Tube{}}

	Distance = func(flat, flong, tlat, tlong float64) float64 { return 300 }

	result := CloseToStation(location, stations)

	if result != true {
		t.Error("Should return true if it's close")
	}
}

func TestGivenStationMore350MetersReturnFalse(t *testing.T) {

	location := data.Location{}
	stations := []data.Tube{data.Tube{}}

	Distance = func(flat, flong, tlat, tlong float64) float64 { return 400 }

	result := CloseToStation(location, stations)

	if result != false {
		t.Error("Should return false if it's far")
	}
}

func BenchmarkCloseToStation(t *testing.B) {
	location := data.Location{}

	stationsForBenchmark := make([]data.Tube, 100)

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		CloseToStation(location, stationsForBenchmark)
	}

}
