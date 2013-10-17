package dispatcher

import (
	"run/data"
	"log"
	"time"
)

type Dispatcher struct {
	ShutdownConfirmation chan string
	Output               chan string
	AskLocation          chan string
	Robots               map[string]data.DispatcherToRobot
}

func (dispatcher *Dispatcher) SliceAndSendBuffer(name string) {
	r := dispatcher.Robots[name]
	if len(r.Locations) != 0 {
		size := data.ROBOT_LOCATION_STORAGE
		if data.ROBOT_LOCATION_STORAGE > len(r.Locations) {
			size = len(r.Locations)
		}

		r.Location <- r.Locations[0:size]
		r.Locations = r.Locations[size:]
		dispatcher.Robots[name] = r
	}
}

func (dispatcher *Dispatcher) ConfirmShutdownFrom(name string) {
	delete(dispatcher.Robots, name)
}

func (dispatcher *Dispatcher) ReadyForShutdown() bool {
	return len(dispatcher.Robots) == 0
}

func (dispatcher *Dispatcher) SendShutdownToAll() {
	for name, robot := range dispatcher.Robots {
		log.Printf("Dispatcher still had %d locations for %s", len(robot.Locations), name)
		robot.Shutdown <- true
	}
}

func (dispatcher *Dispatcher) Run(startTime time.Time, shutdownAll chan bool) {

	for {
		select {
		case name := <-dispatcher.AskLocation:
			log.Printf("Dispatcher asked for locations by %s, remaining: %d", name, len(dispatcher.Robots[name].Locations))
			dispatcher.SliceAndSendBuffer(name)
		case output := <-dispatcher.Output:
			log.Println("Dispatcher received message: " + output)
		case name := <-dispatcher.ShutdownConfirmation:
			log.Printf("Dispatcher received shutdown confirmation by %s", name)
			dispatcher.ConfirmShutdownFrom(name)
			if dispatcher.ReadyForShutdown() {
				return
			}
		case <-shutdownAll:
			log.Println("Dispatcher sending shutdown to all")
			dispatcher.SendShutdownToAll()
		}
	}
}
