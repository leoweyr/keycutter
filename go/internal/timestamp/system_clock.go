package timestamp

import "time"

// SystemClock reads the current instant from the operating-system wall clock.
type SystemClock struct{}

var sharedSystemClock *SystemClock = &SystemClock{}

// SharedSystemClock returns the shared SystemClock singleton.
func SharedSystemClock() *SystemClock {
	return sharedSystemClock
}

// NowUnixSeconds returns the current Unix epoch second count from the system wall clock.
func (systemClock *SystemClock) NowUnixSeconds() uint64 {
	return uint64(time.Now().Unix())
}
