package main

import (
	"run/data"
	"run/dispatcher"
	"run/htime"
	"run/inputs"
	"run/robot"
	"log"
	"time"
)

func BuildDispatcherToRobots(locations map[string][]data.Location) (r map[string]data.DispatcherToRobot) {
	r = make(map[string]data.DispatcherToRobot)

	for name, nameLocations := range locations {
		r[name] = data.DispatcherToRobot{
			Locations: nameLocations,
			Shutdown:  make(chan bool),
			Location:  make(chan []data.Location, data.ROBOT_LOCATION_STORAGE)}
	}
	return
}

func LaunchRobots(robotsData map[string]data.DispatcherToRobot, startTime time.Time, robotToDispatcher data.RobotToDispatcher, stations []data.Tube) {
	for name, data := range robotsData {
		robot := robot.Robot{
			Name:            name,
			CurrentLocation: data.Locations[0],
			Shutdown:        data.Shutdown,
			Location:        data.Location,
			Dispatcher:      robotToDispatcher,
			Stations:        stations}

		go robot.Run(startTime)
	}
}

func ShutdownCountdown(shutdownAll chan bool, startTime time.Time, shutdownTime time.Time) {

	d := htime.DurationUntilReachingDestination(startTime, shutdownTime)
	log.Printf("Start time: %s, Shutdown time: %s, Shutdown in... %s", startTime, shutdownTime, d)

	go func() {
		time.Sleep(d)
		shutdownAll <- true
	}()
}

func main() {
	names := []string{"5937", "6043"}

	shutdownTime := time.Date(2011, time.March, 22, 8, 12, 0, 0, time.UTC)
	shutdownAll := make(chan bool)

	locations, stations, err := inputs.ReadFiles(names, "tube")

	if err != nil {
		log.Fatal(err)
	}

	dispatcherToRobots := BuildDispatcherToRobots(locations)

	startTime := htime.FindStartTime(dispatcherToRobots)

	dispatcherData := data.RobotToDispatcher{AskLocation: make(chan string), ShutdownConfirmation: make(chan string), Output: make(chan string)}

	LaunchRobots(dispatcherToRobots, startTime, dispatcherData, stations)

	dispatcher := dispatcher.Dispatcher{
		ShutdownConfirmation: dispatcherData.ShutdownConfirmation,
		Output:               dispatcherData.Output,
		AskLocation:          dispatcherData.AskLocation,
		Robots:               dispatcherToRobots}

	ShutdownCountdown(shutdownAll, startTime, shutdownTime)

	dispatcher.Run(startTime, shutdownAll)
}
