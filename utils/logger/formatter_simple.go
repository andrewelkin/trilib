package logger

import (
	"fmt"
	"sync/atomic"

	"github.com/mgutz/ansi"
)

type SimpleFormatter struct {
	skipDate  int32  // start next screen output without date.service/level mark
	needCRLF  int32  // start next screen output with new line
	trailCR   bool   // automatically append \n to the string
	ansiReset string // code for resetting ansi
	ansi      bool
}

func (f *SimpleFormatter) String(lm logMessage) string {
	level, hasLevel := logLevels[lm.level]
	if !hasLevel {
		panic(fmt.Errorf("unknown log level: %v", lm.level))
	}
	var cr string
	if f.trailCR {
		cr = "\n"
	}

	var txt = ""
	if f.skipDate == 0 {
		txt = fmt.Sprintf("%s (%s) [%s]: %s%s", formatTime(lm.unixTimestampNS), lm.namespace, level, lm.message, cr)
	} else {
		txt = fmt.Sprintf("%s\n", lm.message)
	}

	if f.ansi {
		if 1 == atomic.SwapInt32(&f.needCRLF, 0) {
			txt = "\n" + txt
		}
		var mod bool
		if txt, mod = expandAnsi(txt); mod {
			txt += f.ansiReset
		}

	} else {
		txt, _ = stripAnsi(txt)
	}

	return txt
}

// NoDateNextLine starts next line without date/debug/servie label
func (f *SimpleFormatter) NoDateNextLine() {
	atomic.StoreInt32(&f.skipDate, 1)
}

// NewLine inserts \n before next output
func (f *SimpleFormatter) NewLine() {
	atomic.StoreInt32(&f.needCRLF, 1)
}

func NewSimpleFormatter(ansiSupport bool, trailCR bool) Formatter {
	return &SimpleFormatter{
		0,
		0,
		trailCR,
		ansi.ColorCode("reset"),
		ansiSupport,
	}
}
