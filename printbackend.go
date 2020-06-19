package lumberjack

import "log"

//PrintBackend implements a console printing Backend that currently
//offers two predefined formats based on the Verbosity specified.
type PrintBackend struct {
	Verbosity LogLevel
}

//Log satisfies the Backend interfaces requirements used for accepting
//LogEntry objects to print out to the console.
func (b *PrintBackend) Log(entry *LogEntry) {
	//TODO: Custom Formatting Templates
	printLog(b.Verbosity, entry)
}

//printLog is an internal function to print the log to the console with
//a predefined format determined by the verbosity LogLevel paramter.
func printLog(verbosity LogLevel, entry *LogEntry) {
	if entry.Level >= verbosity {
		log.Printf("(%s) @ %s() %s:%v: %s", entry.Level, entry.Caller, entry.File, entry.Line, entry.Message)
	} else {
		log.Printf("(%s) @ %s(): %s", entry.Level, entry.Caller, entry.Message)
	}
}
