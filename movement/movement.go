package movement

import (
	"run/data"
	"math"
	"time"
)

var MINIMUM_DISTANCE float64 = 350

func Speed(distance float64, fTime, tTime time.Time) float64 {
	duration := tTime.Sub(fTime)

	if duration == 0 {
		return 0
	}
	return distance / duration.Seconds()
}

func Move(from, to data.Location) float64 {
	distance := Distance(from.Latitude, from.Longitude, to.Latitude, to.Longitude)
	return Speed(distance, from.Time, to.Time)
}

func CloseToStation(location data.Location, stations []data.Tube) bool {
	for _, tube := range stations {
		if Distance(location.Latitude, location.Longitude, tube.Latitude, tube.Longitude) < MINIMUM_DISTANCE {
			return true
		}
	}

	return false
}

//www.johndcook.com/lat_long_distance.html
var Distance = func(fromLat, fromLong, toLat, toLong float64) float64 {
	var radius float64 = 6371.0 //kms
	var phi_1 float64 = (90.0 - fromLat) * math.Pi / 180.0
	var phi_2 float64 = (90.0 - toLat) * math.Pi / 180.0
	var theta_1 float64 = fromLong * math.Pi / 180.0
	var theta_2 float64 = toLong * math.Pi / 180.0

	d := radius * math.Acos(math.Sin(phi_1)*math.Sin(phi_2)*math.Cos(theta_1-theta_2)+math.Cos(phi_1)*math.Cos(phi_2))

	return d * 1000 //meters
}
