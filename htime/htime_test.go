package htime

import (
	"run/data"
	"testing"
	"time"
)

func TestGivenTwoTimesReturnEarliest(t *testing.T) {
	n := time.Now()
	firstLocation := data.Location{Time: n.Add(time.Second)}
	secondLocation := data.Location{Time: n}

	robots := make(map[string]data.DispatcherToRobot)
	firstList := []data.Location{firstLocation}
	secondList := []data.Location{secondLocation}

	robots["one"] = data.DispatcherToRobot{Locations: firstList}
	robots["two"] = data.DispatcherToRobot{Locations: secondList}

	result := FindStartTime(robots)

	if result != n {
		t.Errorf("Should be %s but was %s", n, result)
	}
}

func TestGivenNoTimesReturnDefault(t *testing.T) {
	robots := make(map[string]data.DispatcherToRobot)

	result := FindStartTime(robots)

	if !result.IsZero() {
		t.Errorf("Should be default but was %s", result)
	}
}
