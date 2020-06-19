package lumberjack

//Backend is an interface that must be implemented
//in order to be utilized by an instance of Logger.
type Backend interface {
	Log(*LogEntry)
}
