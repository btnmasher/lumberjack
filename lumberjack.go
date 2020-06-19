//Package lumberjack is a simple leveled logger that provides the means to send
//structured log messages on a number of various backends. New backends can be
//added to the package by implementing the Backend interface.
package lumberjack

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

//Logger holds the configuration for LogLevel state and references to in-use backends.
type Logger struct {
	logLevels map[LogLevel]struct{}
	backends  map[string]Backend
	sync.Mutex
}

//defaultLevels contains sensible defaults for most regular logging needs.
var defaultLevels map[LogLevel]struct{} = map[LogLevel]struct{}{
	INFO:     {},
	WARN:     {},
	ERROR:    {},
	CRITICAL: {},
	FATAL:    {},
}

//NewLogger returns an empty instance of Logger.
func NewLogger() *Logger {
	logger := Logger{}
	logger.logLevels = map[LogLevel]struct{}{}
	logger.backends = map[string]Backend{}
	return &logger
}

//NewLoggerWithDefaults returns an instance of Logger with sensible defaults and a print backend.
func NewLoggerWithDefaults() *Logger {
	logger := Logger{}

	//Start withdefault log levels (all minus DEBUG)
	logger.logLevels = copyMap(defaultLevels)

	//Start with default print logger
	logger.backends = map[string]Backend{"print": &PrintBackend{Verbosity: ERROR}}

	return &logger
}

//copyMap makes a non-reference copy of a map of LogLevel keys.
func copyMap(original map[LogLevel]struct{}) map[LogLevel]struct{} {
	newmap := map[LogLevel]struct{}{}
	for k, v := range original {
		newmap[k] = v
	}
	return newmap
}

//AddLevel adds a LogLevel to the current Logger.
func (l *Logger) AddLevel(level LogLevel) error {
	if !l.levelSet(level) && validLevel(level) {
		l.Lock()
		l.logLevels[level] = struct{}{}
		l.Unlock()
	} else {
		return fmt.Errorf("LogLevel already set: %s", level)
	}
	return nil
}

//RemoveLevel removes a LogLevel from the current Logger.
func (l *Logger) RemoveLevel(level LogLevel) error {
	if l.levelSet(level) {
		l.Lock()
		delete(l.logLevels, level)
		l.Unlock()
	} else {
		return fmt.Errorf("LogLevel not set: %s", level)
	}
	return nil
}

//AddBackend adds an object implementing the Backend interface to the current Logger.
//A name must be specified to add the Backend to the collection as to differentiate
//it from other Backends. This allows multiple instances of the same Backend object
//to be added to the collection with different configurations.
func (l *Logger) AddBackend(name string, backend Backend) error {
	if !l.backendAdded(name) {
		l.Lock()
		l.backends[name] = backend
		l.Unlock()
	} else {
		return fmt.Errorf("Backend with that name already exists: %s", name)
	}
	return nil
}

//GetBackend returns a reference to an object implementing the Backend
//interface associated with the current Logger. The reference is retrieved
//from a map collection with the name of the Backend as the key that was
//Specified when adding the Backend to the collection.
func (l *Logger) GetBackend(name string) (*Backend, error) {
	l.Lock()
	defer l.Unlock()
	if b, exists := l.backends[name]; exists {
		return &b, nil
	}
	return nil, fmt.Errorf("Backend with that name does not exist: %s", name)

}

//RemoveBackend removes a specified object implementing the Backend interface
//from the current Logger. The name is used to specify which reference to be
//removed from the collection.
func (l *Logger) RemoveBackend(name string, backend Backend) error {
	if l.backendAdded(name) {
		l.Lock()
		delete(l.backends, name)
		l.Unlock()
	} else {
		return fmt.Errorf("Backend with that name does not exist: %s", name)
	}
	return nil
}

//Infof logs a formatted string built from the specified args to all added
//Backend objects aded to the current Logger if the DEBUG LogLevel currently
//added to the Logger.
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.levelSet(INFO) {
		l.log(INFO, fmt.Sprintf(format, args...))
	}
}

//Warnf logs a formatted string built from the specified args to all added
//Backend objects aded to the current Logger if the WARN LogLevel currently
//added to the Logger.
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.levelSet(WARN) {
		l.log(WARN, fmt.Sprintf(format, args...))
	}
}

//Errorf logs a formatted string built from the specified args to all added
//Backend objects aded to the current Logger if the ERROR LogLevel currently
//added to the Logger.
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.levelSet(ERROR) {
		l.log(ERROR, fmt.Sprintf(format, args...))
	}
}

//Criticalf logs a formatted string built from the specified args to all added
//Backend objects aded to the current Logger if the CRITICAL LogLevel currently
//added to the Logger.
func (l *Logger) Criticalf(format string, args ...interface{}) {
	if l.levelSet(CRITICAL) {
		l.log(CRITICAL, fmt.Sprintf(format, args...))
	}
}

//Fatalf logs a formatted string built from the specified args to all added
//Backend objects aded to the current Logger if the FATAL LogLevel currently
//added to the Logger, then it will cause the application to os.Exit with status 1
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FATAL, fmt.Sprintf(format, args...))
	os.Exit(1)
}

//Debugf logs a formatted string built from the specified args to all added
//Backend objects aded to the current Logger if the DEBUG LogLevel currently
//added to the Logger.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.levelSet(DEBUG) {
		l.log(DEBUG, fmt.Sprintf(format, args...))
	}
}

//Info logs a string built from the specified args to all added Backend
//objects aded to the current Logger if the INFO LogLevel currently added
//to the Logger.
func (l *Logger) Info(args ...interface{}) {
	if l.levelSet(INFO) {
		l.log(INFO, fmt.Sprint(args...))
	}
}

//Warn logs a string built from the specified args to all added Backend
//objects aded to the current Logger if the WARN LogLevel currently added
//to the Logger.
func (l *Logger) Warn(args ...interface{}) {
	if l.levelSet(WARN) {
		l.log(WARN, fmt.Sprint(args...))
	}
}

//Error logs a string built from the specified args to all added Backend
//objects aded to the current Logger if the WARN LogLevel currently added
//to the Logger.
func (l *Logger) Error(args ...interface{}) {
	if l.levelSet(ERROR) {
		l.log(ERROR, fmt.Sprint(args...))
	}
}

//Critical logs a string built from the specified args to all added Backend
//objects aded to the current Logger if the CRITICAL LogLevel currently added
//to the Logger.
func (l *Logger) Critical(args ...interface{}) {
	if l.levelSet(CRITICAL) {
		l.log(CRITICAL, fmt.Sprint(args...))
	}
}

//Debug logs a string built from the specified args to all added Backend
//objects aded to the current Logger if the DEBUG LogLevel currently added
//to the Logger.
func (l *Logger) Debug(args ...interface{}) {
	if l.levelSet(DEBUG) {
		l.log(DEBUG, fmt.Sprint(args...))
	}
}

//Fatal logs a string built from the specified args to all added Backend
//objects aded to the current Logger if the FATAL LogLevel currently added
//to the Logger, then it will cause the application to os.Exit with status 1.
func (l *Logger) Fatal(args ...interface{}) {
	l.log(FATAL, fmt.Sprint(args...))
	os.Exit(1)
}

//validLevel checks the specified LogLevel if it is a valid LogLevel constant
//and returns true or false based on that check.
func validLevel(level LogLevel) bool {
	return level >= INFO && level <= DEBUG
}

//levelSet checks the specified LogLevel if it is added to the current Logger
//and returns true or false based on that check.
func (l *Logger) levelSet(level LogLevel) bool {
	l.Lock()
	defer l.Unlock()
	_, exists := l.logLevels[level]
	return exists
}

//backendAdded checks the specified name string of a Backend if it is added
//to the current Logger and returns true or false based on that check.
func (l *Logger) backendAdded(name string) bool {
	l.Lock()
	defer l.Unlock()
	_, exists := l.backends[name]
	return exists
}

//log will accept the specified LogLevel and message, build a LogEntry
//from that information, then send it to all backends added to the
//current Logger.
func (l *Logger) log(level LogLevel, message string) {
	entry := buildLogEntry(level, message)
	l.sendToBackends(entry)
}

//buildLogEntry accepts a specified LogLevel and message string, uses
//the Go runtime to determine where the original call to log originated
//with the name of the source file, line number, and function block it was
//called from.
//
//It then returns a pointer to a LogEntry with this
//information contianed within the fields for consumption by the
//various objects implementing the Backend interface.
func buildLogEntry(level LogLevel, message string) *LogEntry {
	pc, file, line, ok := runtime.Caller(3)
	fname := "???"
	path := ""
	if !ok {
		file = "???"
		line = 0
	} else {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
		p, f := filepath.Split(file)
		file = f
		path = p
	}
	return &LogEntry{
		Level:   level,
		Path:    path,
		File:    file,
		Line:    line,
		Caller:  fname,
		Message: message,
	}
}

//sendToBackends accepts a specified LogEntry, then calls the Log
//function on all backends added to the current Logger.
func (l *Logger) sendToBackends(entry *LogEntry) {
	l.Lock()
	defer l.Unlock()
	for _, backend := range l.backends {
		backend.Log(entry)
	}
}

//logInternalf is a function that is used to log errors that occur
//within the scope of the lumberjack package itself. It accepts a
//LogLevel, a string format, and arbitrary args to format a string
//and send it using a PrintBackend.
func logInteralf(level LogLevel, format string, args ...interface{}) {
	sendToInternal(level, fmt.Sprintf(format, args...))
}

//logInternalf is a function that is used to log errors that occur
//within the scope of the lumberjack package itself. It accepts a
//LogLevel, and arbitrary args to build a string and send it using
//a PrintBackend.
func logInternal(level LogLevel, args ...interface{}) {
	sendToInternal(level, fmt.Sprint(args...))
}

//sendToInternal is a function that will accept a LogLevel and a
//message string, build a LogEntry, then send it on a PrintBackend.
//This function is used for internal logging of errors that occur
//within the scope of the lumberjack package itself.
func sendToInternal(level LogLevel, message string) {
	entry := buildLogEntry(level, message)
	printLog(level, entry)
}
