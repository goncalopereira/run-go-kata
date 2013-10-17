package htime

import (
	"run/data"
	"time"
)

const WAIT_FOR_MORE_LOCATIONS = 5 * time.Minute

const SECONDS_TO_10MS = 100

const TIME_FORMAT = "15:04:05"

func FindStartTime(robots map[string]data.DispatcherToRobot) (startTime time.Time) {
	for _, r := range robots {
		if startTime.IsZero() || r.Locations[0].Time.Before(startTime) {
			startTime = r.Locations[0].Time
		}
	}

	return
}

func DurationUntilReachingDestination(currentTime time.Time, nextLocationTime time.Time) time.Duration {
	return nextLocationTime.Sub(currentTime) / SECONDS_TO_10MS
}
