package lumberjack

import (
	"encoding/json"
	"fmt"
)

type LogLevel byte

//Constants used to define the various LogLevels
const (
	INFO LogLevel = iota
	WARN
	ERROR
	CRITICAL
	FATAL
	DEBUG
)

//logLevelNameToValue is a map that will allow for conversion
//of a string to a LogLevel.
var logLevelNameToValue = map[string]LogLevel{
	"INFO":     INFO,
	"WARN":     WARN,
	"ERROR":    ERROR,
	"CRITICAL": CRITICAL,
	"FATAL":    FATAL,
	"DEBUG":    DEBUG,
}

//logLevelNameToValue is a map that will allow for conversion
//of a LogLevel to a string.
var logLevelValueToName = map[LogLevel]string{
	INFO:     "INFO",
	WARN:     "WARN",
	ERROR:    "ERROR",
	CRITICAL: "CRITICAL",
	FATAL:    "FATAL",
	DEBUG:    "DEBUG",
}

//String satisfies fmt.Stringer interface fo use in Marshalling
//a LogLevel to JSON or for console printing.
func (l LogLevel) String() string {
	if v, exists := logLevelValueToName[l]; exists {
		return v
	}
	return ""
}

// MarshalJSON satisfies json.Marshaler.
func (l LogLevel) MarshalJSON() ([]byte, error) {
	s, ok := logLevelValueToName[l]
	if !ok {
		return nil, fmt.Errorf("invalid LogLevel: %d", l)
	}
	return json.Marshal(s)
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (l *LogLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("LogLevel should be a string, got %s", data)
	}

	v, ok := logLevelNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid LogLevel %q", s)
	}

	*l = v
	return nil
}
