package utils

import (
	"fmt"
	"strings"
	"time"
)

// UnixTimeNowMS returns the current unix time in milliseconds.
func UnixTimeNowMS() int64 {
	return time.Now().UnixNano() / 1000000
}

// UnixTimeNowMicroSeconds returns the current unix time in microseconds.
func UnixTimeNowMicroSeconds() int64 {
	return time.Now().UnixNano() / 1000
}

// UnixTimeMS converts time to  unix time in milliseconds.
func UnixTimeMS(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// enables an action once in n sec
type OnceInNSec struct {
	tLast    int64
	interval int64
}

// Init inits in seconds
func (o *OnceInNSec) Init(intervalSec int64) *OnceInNSec {
	o.interval = intervalSec * 1000 // in msec
	return o
}

// Reset will trigger action imm
func (o *OnceInNSec) Reset() {
	o.tLast = 0
}

// Expire in n seconds from now
func (o *OnceInNSec) ExpireIn(expireInSeconds int64) {
	o.tLast = UnixTimeNowMS() - expireInSeconds*1000
}

// StartOver resets waiting period
func (o *OnceInNSec) StartOver() {
	o.tLast = UnixTimeNowMS()
}

// returns true if interval seconds elapsed since last action of if action never happen or after reset
func (o *OnceInNSec) CanDo() bool {
	tnow := UnixTimeNowMS()
	if tnow > o.tLast+o.interval {
		o.tLast = tnow
		return true
	}
	return false
}

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

func PrintDuration(d time.Duration) string {
	if d < day {
		return d.String()
	}

	var b strings.Builder
	if d >= year {
		years := d / year
		fmt.Fprintf(&b, "%dy", years)
		d -= years * year
	}

	days := d / day
	d -= days * day
	fmt.Fprintf(&b, "%dd%s", days, d)

	return b.String()
}
