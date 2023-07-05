package logger

import (
	"bytes"
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"log"
)

// TODO (@hrharder) -- multiple goroutine tests (concurrent logging)?

func TestAsyncLogger(t *testing.T) {
	// stub the clock
	currentTime, _ := time.Parse(time.RFC3339Nano, "2020-09-29T19:05:07.123456Z")
	currentClock = testClock(currentTime)
	defer resetClock()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var capturedOutput bytes.Buffer
	var outputs []logOutput
	testLogger := &AsyncLogger{
		logs: make(chan logMessage),
		outputs: append(outputs, logOutput{
			dst:       &capturedOutput,
			minLevel:  LogLevelDebug,
			filter:    FilterMatchAll,
			formatter: NewSimpleFormatter(false, true),
		}),
	}
	go testLogger.handleLogs(ctx)

	testLogger.Debugf("test-logger", "debugf with %v%v", "args", ".")
	testLogger.Infof("test-logger", "infof with %v%v", "args", ".")
	testLogger.Warnf("test-logger", "warnf with %v%v", "args", ".")
	assert.PanicsWithValue(t, "errorf with args.", func() {
		testLogger.Errorf("test-logger", "errorf with %v%v", "args", ".")
	})

	// ensure all the logs are written before comparing
	time.Sleep(1 * time.Millisecond)

	expectedOutput := "2020-09-29 19:05:07.123 (test-logger) [DEBUG]: debugf with args.\n" +
		"2020-09-29 19:05:07.123 (test-logger) [INFO]: infof with args.\n" +
		"2020-09-29 19:05:07.123 (test-logger) [WARN]: warnf with args.\n" +
		"2020-09-29 19:05:07.123 (test-logger) [ERROR]: errorf with args.\n"

	assert.Equal(t, expectedOutput, capturedOutput.String())
}

func TestAsyncLoggerLogLevels(t *testing.T) {
	// stub the clock
	currentTime, _ := time.Parse(time.RFC3339Nano, "2020-09-29T19:05:07.123456Z")
	currentClock = testClock(currentTime)
	defer resetClock()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var capturedOutput bytes.Buffer
	var outputs []logOutput
	testLogger := &AsyncLogger{
		logs: make(chan logMessage),
		outputs: append(outputs, logOutput{
			dst:       &capturedOutput,
			minLevel:  LogLevelWarn,
			filter:    FilterMatchAll,
			formatter: NewSimpleFormatter(false, true),
		}),
	}
	go testLogger.handleLogs(ctx)

	// if level is warn, we should only see the warnings and up in the output
	testLogger.Debugf("test-logger", "this log should not show")
	testLogger.Infof("test-logger", "this log should not show")
	testLogger.Warnf("test-logger", "this log should show")

	// ensure all the logs are written before comparing
	time.Sleep(1 * time.Millisecond)

	expectedOutput := "2020-09-29 19:05:07.123 (test-logger) [WARN]: this log should show\n"
	assert.Equal(t, expectedOutput, capturedOutput.String())
}

func TestSpecialPrefix(t *testing.T) {
	// stub the clock
	currentTime, _ := time.Parse(time.RFC3339Nano, "2020-09-29T19:05:07.123456Z")
	currentClock = testClock(currentTime)
	defer resetClock()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// two test outputs: one that is the "screen", one that is the "file"
	var screenBuffer, fileBuffer bytes.Buffer
	outputs := []logOutput{
		// this should get all except where the namespace is or starts with "_"
		{dst: &screenBuffer, minLevel: LogLevelDebug, filter: FilterUnderscore, formatter: NewSimpleFormatter(false, true)},

		// this should get all logs regardless of namespace
		{dst: &fileBuffer, minLevel: LogLevelDebug, filter: FilterMatchAll, formatter: NewSimpleFormatter(false, true)},
	}

	testLogger := &AsyncLogger{
		logs: make(chan logMessage),

		outputs: outputs,
	}
	go testLogger.handleLogs(ctx)

	testLogger.Infof("ns", "should show everywhere")
	testLogger.Infof("ns_with_space", "should show everywhere")
	testLogger.Infof("_", "should not show in screen")
	testLogger.Infof("_ns", "should not show in screen")
	testLogger.Infof("__", "should not show in screen")

	// ensure all the logs are written before comparing
	time.Sleep(1 * time.Millisecond)

	expectedDumpFromFile := "" +
		"2020-09-29 19:05:07.123 (ns) [INFO]: should show everywhere\n" +
		"2020-09-29 19:05:07.123 (ns_with_space) [INFO]: should show everywhere\n" +
		"2020-09-29 19:05:07.123 (_) [INFO]: should not show in screen\n" +
		"2020-09-29 19:05:07.123 (_ns) [INFO]: should not show in screen\n" +
		"2020-09-29 19:05:07.123 (__) [INFO]: should not show in screen\n"

	expectedDumpFromScreen := "" +
		"2020-09-29 19:05:07.123 (ns) [INFO]: should show everywhere\n" +
		"2020-09-29 19:05:07.123 (ns_with_space) [INFO]: should show everywhere\n"

	assert.Equal(t, expectedDumpFromScreen, screenBuffer.String())
	assert.Equal(t, expectedDumpFromFile, fileBuffer.String())
}

var testLogger Logger
var stopFunc context.CancelFunc

func resetTestLogger() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	testLogger = NewAsyncLogger(ctx, LogLevelDebug, FilterMatchAll)
	stopFunc = cancelFunc
}

func BenchmarkAsyncLogger(b *testing.B) {
	b.StopTimer()
	resetTestLogger()
	b.N = math.MaxInt32
	b.StartTimer()

	testLogger.Warnf("TST", "here is a warning, %v %v %v %v", "with", "many", "different", "args")
}

func BenchmarkStandardLogger(b *testing.B) {
	b.N = math.MaxInt32

	log.Printf("here is a warning, %v %v %v %v", "with", "many", "different", "args")
}
