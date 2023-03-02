package logger

import "time"

var currentClock clock

func init() {
	resetClock()
}

func resetClock() {
	currentClock = realClock{}
}

type clock interface {
	Now() time.Time
}

type realClock struct{}

func (rc realClock) Now() time.Time {
	return time.Now()
}

type testClock time.Time

func (tc testClock) Now() time.Time {
	return time.Time(tc)
}
