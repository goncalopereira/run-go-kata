package data

import "time"

const ROBOT_LOCATION_STORAGE = 10

type Tube struct {
	Latitude  float64
	Longitude float64
}

type Location struct {
	Latitude  float64
	Longitude float64
	Time      time.Time
}

type DispatcherToRobot struct {
	Shutdown  chan bool
	Location  chan []Location
	Locations []Location
}

type RobotToDispatcher struct {
	ShutdownConfirmation chan string
	Output               chan string
	AskLocation          chan string
}
