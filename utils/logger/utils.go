package logger

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// FilterUnderscore will match all namespaces except those prefixed with "_"
var FilterUnderscore = Filter(regexp.MustCompile(`^([^_]+(.*)?)?$`))

// ParseLogLevel tries to parse raw into a log level, if it cant, returns defaultLevel
func ParseLogLevel(raw string, defaultLevel LogLevel) LogLevel {
	for level, levelString := range logLevels {
		if strings.ToLower(raw) == strings.ToLower(levelString) {
			return level
		}
	}
	return defaultLevel
}

// 2020-10-15 10:28:21.333
const timeFormat = "%d-%02d-%02d %02d:%02d:%02d.%03d"

func formatTime(unixNanoseconds int64) string {
	t := time.Unix(0, unixNanoseconds).UTC()

	yr, mo, day := t.Year(), t.Month(), t.Day()
	hr, min, sec, nsec := t.Hour(), t.Minute(), t.Second(), t.Nanosecond()
	ms := nsec / 1e6

	return fmt.Sprintf(timeFormat, yr, mo, day, hr, min, sec, ms)
}

func newLogMessage(level LogLevel, namespace, format string, a ...interface{}) logMessage {
	return logMessage{
		level:           level,
		namespace:       namespace,
		message:         fmt.Sprintf(format, a...),
		unixTimestampNS: currentClock.Now().UnixNano(),
	}
}

var logLevels = map[LogLevel]string{
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARN",
	LogLevelError: "ERROR",
	LogLevelFatal: "FATAL",
}
