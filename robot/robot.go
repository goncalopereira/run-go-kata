package robot

import (
	"fmt"
	"run/data"
	"run/htime"
	"run/movement"
	"log"
	"math/rand"
	"time"
)

var TRAFFIC_CONDITIONS = []string{"HEAVY", "LIGHT", "MODERATE"}

type Robot struct {
	Stations          []data.Tube
	Name              string
	CurrentLocation   data.Location
	NextLocations     []data.Location
	TimeUntilNextStop time.Duration
	Dispatcher        data.RobotToDispatcher
	Location          chan []data.Location
	Shutdown          chan bool
}

func FindTraffic() string {
	return TRAFFIC_CONDITIONS[rand.Intn(len(TRAFFIC_CONDITIONS))]
}

func (robot *Robot) Move() float64 {
	speed := movement.Move(robot.CurrentLocation, robot.NextLocations[0])
	robot.CurrentLocation = robot.NextLocations[0]
	robot.NextLocations = robot.NextLocations[1:]
	return speed
}

func (robot *Robot) CheckStationsAndReport(speed float64) {
	if speed == 0 {
		return
	}

	if movement.CloseToStation(robot.CurrentLocation, robot.Stations) {
		condition := FindTraffic()
		robot.Dispatcher.Output <- fmt.Sprintf("%s,%s,%fm/s,%s", robot.Name, robot.CurrentLocation.Time.Format(htime.TIME_FORMAT), speed, condition)
	}
}

func (robot *Robot) SyncWithStartTimeToStart(startTime time.Time) {
	startIn := htime.DurationUntilReachingDestination(startTime, robot.CurrentLocation.Time)
	log.Printf("%s wait to start for %s\n", robot.Name, startIn)
	<-time.After(startIn)
}

func (robot *Robot) Log(message string) {
	log.Printf("robot:%s,time:%s,remaining:%d,waited up to:%s=%s\n", robot.Name, robot.CurrentLocation.Time.Format(htime.TIME_FORMAT), len(robot.NextLocations), robot.TimeUntilNextStop, message)
}

func (robot *Robot) GoingToStartMoving() {
	robot.TimeUntilNextStop = htime.DurationUntilReachingDestination(robot.CurrentLocation.Time, robot.NextLocations[0].Time)
}

func (robot *Robot) GoingToWait() {
	robot.TimeUntilNextStop = htime.WAIT_FOR_MORE_LOCATIONS
}

func (robot *Robot) Run(startTime time.Time) {
	rand.Seed(time.Now().UnixNano())

	robot.SyncWithStartTimeToStart(startTime)

	robot.Log("starting up")
	for {
		if len(robot.NextLocations) == 0 {
			robot.Log("asking for locations")
			robot.Dispatcher.AskLocation <- robot.Name
			robot.GoingToWait()
		} else {
			robot.GoingToStartMoving()
		}

		select {
		case robot.NextLocations = <-robot.Location:
			robot.Log("received locations")
		case <-time.After(robot.TimeUntilNextStop):
			speed := robot.Move()
			robot.Log("finished move")
			robot.CheckStationsAndReport(speed)
		case <-robot.Shutdown:
			robot.Log("shutting down")
			robot.Dispatcher.ShutdownConfirmation <- robot.Name
			return
		}
	}

}
