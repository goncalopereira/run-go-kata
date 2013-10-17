package dispatcher

import (
	"run/data"
	"testing"
	"time"
)

func TestGivenArrayWithMoreThanMaxAmountSendMaxAmountAndReturnRemaining(t *testing.T) {

	name := "driver"
	maxAmount := data.ROBOT_LOCATION_STORAGE

	robotLocation := make(chan []data.Location, maxAmount)
	remainingSize := maxAmount + 5
	remainingLocations := make([]data.Location, remainingSize)
	n := time.Now()

	for i := 0; i < remainingSize; i++ {
		remainingLocations[i] = data.Location{Time: n.Add(time.Duration(i) * time.Second)}
	}

	nextZeroTime := remainingLocations[maxAmount].Time

	robot := data.DispatcherToRobot{
		Locations: remainingLocations,
		Location:  robotLocation}

	robots := make(map[string]data.DispatcherToRobot)
	robots[name] = robot

	dispatcher := Dispatcher{Robots: robots}

	dispatcher.SliceAndSendBuffer(name)

	if len(dispatcher.Robots[name].Locations) != (remainingSize - maxAmount) {
		t.Errorf("Length is %d should be %d", len(dispatcher.Robots[name].Locations), (remainingSize - maxAmount))
	}

	if dispatcher.Robots[name].Locations[0].Time != nextZeroTime {
		t.Errorf("Time for 0 in Slice is %s, should be %s", dispatcher.Robots[name].Locations[0].Time, nextZeroTime)
	}

	result := <-robotLocation
	if len(result) != maxAmount {
		t.Errorf("Length is %d, should be %d", len(result), maxAmount)
	}

	if result[0].Time != n {
		t.Errorf("Time for 0 in CHan is %s, should be %s", result[0].Time, n)
	}

}

func TestGivenArrayWithLessThanMaxAmountSendAllAndReturnEmpty(t *testing.T) {
	name := "driver"
	maxAmount := data.ROBOT_LOCATION_STORAGE

	robotLocation := make(chan []data.Location, maxAmount)
	remainingSize := maxAmount - 5
	remainingLocations := make([]data.Location, remainingSize)
	n := time.Now()

	for i := 0; i < remainingSize; i++ {
		remainingLocations[i] = data.Location{Time: n.Add(time.Duration(i) * time.Second)}
	}

	robot := data.DispatcherToRobot{
		Locations: remainingLocations,
		Location:  robotLocation}

	robots := make(map[string]data.DispatcherToRobot)
	robots[name] = robot

	dispatcher := Dispatcher{Robots: robots}

	dispatcher.SliceAndSendBuffer(name)

	if len(dispatcher.Robots[name].Locations) != 0 {
		t.Errorf("Length is %d, should be 0", len(dispatcher.Robots[name].Locations))
	}

	result := <-robotLocation

	if len(result) != remainingSize {
		t.Errorf("Length of chan is %d should be %d", len(result), remainingSize)
	}

	if result[0].Time != n {
		t.Errorf("Time for 0 at chan is %s should be %s", result[0].Time, n)
	}
}

func BenchmarkSliceAndSend(t *testing.B) {
	name := "driver"

	robotLocation := make(chan []data.Location, data.ROBOT_LOCATION_STORAGE)

	robot := data.DispatcherToRobot{Location: robotLocation}

	robots := make(map[string]data.DispatcherToRobot)
	robots[name] = robot

	dispatcher := Dispatcher{Robots: robots}

	for i := 0; i < t.N; i++ {
		r := dispatcher.Robots[name]
		r.Locations = make([]data.Location, data.ROBOT_LOCATION_STORAGE)
		dispatcher.Robots[name] = r

		dispatcher.SliceAndSendBuffer(name)
		<-robotLocation
	}
}
