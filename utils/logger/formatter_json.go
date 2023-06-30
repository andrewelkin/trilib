package logger

import (
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
}

type jsonLogMessage struct {
	Severity  string `json:"severity"`
	Time      int64  `json:"time"`
	Namespace string `json:"context"`
	Message   string `json:"message"`
}

// Implement the String method for JSONFormatter
func (f *JSONFormatter) String(lm logMessage) string {
	level, hasLevel := logLevels[lm.level]
	if !hasLevel {
		panic(fmt.Errorf("unknown log level: %v", lm.level))
	}
	var msg = jsonLogMessage{
		Severity:  level,
		Time:      lm.unixTimestampNS,
		Namespace: lm.namespace,
		Message:   lm.message,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Errorf("failed to marshal log message to JSON: %v", err))
	}

	return string(data)
}

// Implement the NoDateNextLine method for JSONFormatter
func (f *JSONFormatter) NoDateNextLine() {
}

// Implement the NewLine method for JSONFormatter
func (f *JSONFormatter) NewLine() {
}

func NewJsonFormatter() Formatter {
	return &JSONFormatter{}
}
