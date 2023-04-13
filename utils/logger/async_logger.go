package logger

import (
	"context"
	"fmt"
	"github.com/mgutz/ansi"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

// LogDefaultFilter is the regular expression that is used in new loggers as the filter applied to the default
// output (stdout). By default, it is set to filter out logs where the namespace is prefixed with an underscore.
var LogDefaultFilter *regexp.Regexp = FilterUnderscore

// LogBufferSize is the size of the pending log buffer (number of logs in the write queue that can be present until log
// requests start blocking in the calling go-routine).
var LogBufferSize = 256

// LogCacheSize is the size of the log cache: logs that are maintained before/after writing until they are dumped to a
// file (or other output) at the request of the user. Sets the number of recent logs retained (FIFO order).
var LogCacheSize = 1000

var defaultScreenDst io.Writer = os.Stdout

var DefaultScreenOutputFunc = func(dst io.Writer, txt string) {
	if dst == nil {
		dst = defaultScreenDst
	}
	io.WriteString(dst, txt)
}

// AsyncLogger implements Logger and handles logs from writing go-routines in an asynchronous manner.
type AsyncLogger struct {
	mux       sync.RWMutex
	wg        sync.WaitGroup
	blockLogs bool

	logs    chan logMessage
	outputs []logOutput

	skipDate  int32  // start next screen output without date.service/level mark
	needCRLF  int32  // start next screen output with new line
	ansiReset string // code for resetting ansi
	//	ansiBlack string // code for black -- fixing a bug in ansi lib
}

func SetDefaultScreenIO(dst io.Writer) {
	defaultScreenDst = dst
}

func SetDefaultScreenIOfunc(f func(io.Writer, string)) {
	DefaultScreenOutputFunc = f
}

// NewAsyncLogger returns a new async logger where writes are handled in an independent goroutine
// Any logs sent after the context is cancelled or expired may be lost.
//
// If stdoutFilter is specified, logs written to stdout must have a namspace that matches. If no
// stdout filter is specified, the regular expression defined by LogDefaultFilter is used.
func NewAsyncLogger(ctx context.Context, level LogLevel, stdoutFilter *regexp.Regexp) *AsyncLogger {
	logger := &AsyncLogger{
		logs:      make(chan logMessage, LogBufferSize),
		ansiReset: ansi.ColorCode("reset"),
		//ansiBlack: ansi.ColorCode("black"),
	}

	if stdoutFilter == nil {
		stdoutFilter = LogDefaultFilter
	}
	logger.outputs = append(logger.outputs, logOutput{filter: stdoutFilter, dst: defaultScreenDst, minLevel: level, ansi: true, trailCR: true})

	go logger.handleLogs(ctx)
	return logger
}

// Debugf implements Logger
func (lgr *AsyncLogger) Debugf(namespace, format string, a ...interface{}) {
	lgr.withLogsNotBlocked(func() {
		lgr.logs <- newLogMessage(LogLevelDebug, namespace, format, a...)
	})
}

// Infof implements Logger
func (lgr *AsyncLogger) Infof(namespace, format string, a ...interface{}) {
	lgr.withLogsNotBlocked(func() {
		lgr.logs <- newLogMessage(LogLevelInfo, namespace, format, a...)
	})
}

// Warnf implements Logger
func (lgr *AsyncLogger) Warnf(namespace, format string, a ...interface{}) {
	lgr.withLogsNotBlocked(func() {
		lgr.logs <- newLogMessage(LogLevelWarn, namespace, format, a...)
	})
}

// Errorf implements Logger
// Will trigger panic in the calling goroutine
func (lgr *AsyncLogger) Errorf(namespace, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	lgr.withLogsNotBlocked(func() {
		lgr.logs <- newLogMessage(LogLevelError, namespace, message)
	})
	panic(message)
}

// Fatalf implements Logger
// Will trigger immediate process termination from the log writer goroutine
func (lgr *AsyncLogger) Fatalf(namespace, format string, a ...interface{}) {
	lgr.withLogsNotBlocked(func() {
		lgr.logs <- newLogMessage(LogLevelFatal, namespace, format, a...)
	})
}

// AddOutput implements Logger
func (lgr *AsyncLogger) AddOutput(filter *regexp.Regexp, output io.Writer, minLevel LogLevel, ansi bool, trailCR bool) {
	lgr.outputs = append(lgr.outputs, logOutput{filter: filter, minLevel: minLevel, dst: output, ansi: ansi, trailCR: trailCR})
}

// NewLine inserts \n before next output
func (lgr *AsyncLogger) NewLine() {
	atomic.StoreInt32(&lgr.needCRLF, 1)

}

// NoDateNextLine starts next line without date/debug/servie label
func (lgr *AsyncLogger) NoDateNextLine() {
	atomic.StoreInt32(&lgr.skipDate, 1)
}

// Flush implements Logger
func (lgr *AsyncLogger) Flush() {
	for {
		select {
		case log := <-lgr.logs:
			lgr.writeLogAccordingToLevel(log)
			break
		default:
			return
		}
	}
}

func (lgr *AsyncLogger) handleLogs(ctx context.Context) {
	for {
		select {
		case log := <-lgr.logs:
			lgr.writeLogAccordingToLevel(log)
		case <-ctx.Done():
			lgr.mux.Lock()
			defer lgr.mux.Unlock()

			// flush any pending logs, then stop accepting new onces
			lgr.Flush()
			lgr.blockLogs = true
			close(lgr.logs)
			return
		}
	}
}

func (lgr *AsyncLogger) withLogsNotBlocked(f func()) {
	lgr.mux.RLock()
	defer lgr.mux.RUnlock()

	if lgr.blockLogs {
		return
	}

	f()
}

func expandOrStripAnsi(t string, expand bool) (string, bool) {
	modified := false
	for {
		n := strings.IndexByte(t, '{')
		if n < 0 {
			break
		}
		n1 := strings.IndexByte(t[n:], '}')
		if n1 < 0 {
			break
		}
		if n1 > 15 {
			break
		}

		n1 += n
		cc := ansi.ColorCode(t[n+1 : n1])
		if len(cc) != 0 {
			if !expand {
				cc = ""
			}
			t = t[:n] + cc + t[n1+1:]
			modified = true
		} else {
			break
		}
	}
	return t, modified
}

func expandAnsi(t string) (string, bool) {
	return expandOrStripAnsi(t, true)
}

func stripAnsi(t string) (string, bool) {
	return expandOrStripAnsi(t, false)
}

func (lgr *AsyncLogger) writeLogAccordingToLevel(msg logMessage) {
	var wg sync.WaitGroup
	for _, outputConfig := range lgr.outputs {
		if !outputConfig.filter.MatchString(msg.namespace) {
			continue
		}

		if uint(outputConfig.minLevel) > uint(msg.level) {
			return
		}

		wg.Add(1)
		go func(dst io.Writer, txt string, ansi bool) {
			defer wg.Done()

			if ansi {
				if 1 == atomic.SwapInt32(&lgr.needCRLF, 0) {
					txt = "\n" + txt
				}
				var mod bool
				if txt, mod = expandAnsi(txt); mod {
					txt += lgr.ansiReset
				}

			} else {
				txt, _ = stripAnsi(txt)
			}

			if ansi {
				DefaultScreenOutputFunc(dst, txt)
			} else {
				io.WriteString(dst, txt)
			}
		}(outputConfig.dst, msg.String(0 == atomic.SwapInt32(&lgr.skipDate, 0), outputConfig.trailCR), outputConfig.ansi)
	}

	wg.Wait()
	if msg.level == LogLevelFatal {
		lgr.Flush()
		os.Exit(1)
	}
}

type logMessage struct {
	level           LogLevel
	namespace       string
	message         string
	unixTimestampNS int64
}

// String implements stringer
func (lm logMessage) String(needPrefxes bool, needCR bool) string {
	level, hasLevel := logLevels[lm.level]
	if !hasLevel {
		panic(fmt.Errorf("unknown log level: %v", lm.level))
	}
	var cr string
	if needCR {
		cr = "\n"
	}
	if needPrefxes {
		return fmt.Sprintf("%s (%s) [%s]: %s%s", formatTime(lm.unixTimestampNS), lm.namespace, level, lm.message, cr)
	} else {
		return fmt.Sprintf("%s\n", lm.message)
	}
}

type logOutput struct {
	filter   *regexp.Regexp
	ansi     bool
	trailCR  bool
	dst      io.Writer
	minLevel LogLevel
}
