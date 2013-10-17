package robot

import (
	"run/data"
	"math"
	"strings"
	"testing"
	"time"
)

func TestGivenKnownLocationsToMoveReturnCorrectSpeed(t *testing.T) {

	n := time.Now()

	nextOne := data.Location{Time: n.Add(time.Duration(1) * time.Second), Latitude: 123, Longitude: 1}

	robot := Robot{
		CurrentLocation: data.Location{Time: n, Latitude: 123, Longitude: 0},
		NextLocations:   []data.Location{nextOne}}

	speed := robot.Move()

	if math.Floor(speed) != 60560 {
		t.Errorf("Speed should be around 60k m/s and it was %f", speed)
	}
}

func TestGivenLocationsReturnNewCurrentLocationAndNewocations(t *testing.T) {
	n := time.Now()

	nextOne := data.Location{Time: n.Add(time.Duration(1) * time.Second), Latitude: 123, Longitude: 1}
	nextTwo := data.Location{Time: n.Add(time.Duration(2) * time.Second)}
	robot := Robot{
		CurrentLocation: data.Location{Time: n, Latitude: 123, Longitude: 0},
		NextLocations:   []data.Location{nextOne, nextTwo}}

	robot.Move()

	if robot.CurrentLocation.Time != nextOne.Time {
		t.Error("Didn't change current location")
	}

	if len(robot.NextLocations) != 1 {
		t.Error("Didn't slice locations")
	}

	if robot.NextLocations[0].Time != nextTwo.Time {
		t.Error("Didn't show the correct new list")
	}
}

func TestGivenCurrentLocationIsNextToStationOutputToDispatcher(t *testing.T) {
	robot := Robot{
		CurrentLocation: data.Location{},
		Stations:        []data.Tube{data.Tube{}},
		Name:            "driver",
		Dispatcher:      data.RobotToDispatcher{Output: make(chan string, 2)}}

	speed := 5.0

	robot.CheckStationsAndReport(speed)

	result := <-robot.Dispatcher.Output
	if !strings.Contains(result, robot.Name) {
		t.Errorf("Output should have name, was %s", robot.Name)
	}
}
