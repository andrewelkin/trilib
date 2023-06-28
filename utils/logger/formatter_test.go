package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONFormatter_String(t *testing.T) {
	f := &JSONFormatter{}
	lm := logMessage{
		level:           LogLevelInfo,
		unixTimestampNS: 1612345678901234567,
		namespace:       "app",
		message:         "This is a log message",
	}
	expected := `{"severity":"INFO","time":1612345678901234567,"context":"app","message":"This is a log message"}`
	assert.Equal(t, expected, f.String(lm))
}

func TestSimpleFormatter_String(t *testing.T) {
	f := &SimpleFormatter{}
	lm := logMessage{
		level:           LogLevelInfo,
		unixTimestampNS: 1612345678901234567,
		namespace:       "app",
		message:         "This is a log message",
	}
	expected := "2021-02-03 09:47:58.901 (app) [INFO]: This is a log message"
	assert.Equal(t, expected, f.String(lm))
}

func TestSimpleFormatter_NoDateNextLine(t *testing.T) {
	f := &SimpleFormatter{}
	f.NoDateNextLine()
	lm := logMessage{
		level:           LogLevelInfo,
		unixTimestampNS: 1612345678901234567,
		namespace:       "app",
		message:         "This is a log message",
	}
	expected := "This is a log message\n"
	assert.Equal(t, expected, f.String(lm))
}

func TestSimpleFormatter_NewLine(t *testing.T) {
	f := NewSimpleFormatter(false, true)
	f.NewLine()
	lm := logMessage{
		level:           LogLevelInfo,
		unixTimestampNS: 1612345678901234567,
		namespace:       "app",
		message:         "This is a log message",
	}
	expected := "2021-02-03 09:47:58.901 (app) [INFO]: This is a log message\n"
	assert.Equal(t, expected, f.String(lm))
}

func TestNewSimpleFormatter(t *testing.T) {
	f := NewSimpleFormatter(true, false)
	assert.Equal(t, true, f.(*SimpleFormatter).ansi)
	assert.Equal(t, int32(0), f.(*SimpleFormatter).needCRLF)
}

func TestNewJsonFormatter(t *testing.T) {
	f := NewJsonFormatter()
	assert.IsType(t, &JSONFormatter{}, f)
}
