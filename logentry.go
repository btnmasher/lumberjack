package lumberjack

//LogEntry is the object used to contain the relevant
//information for a particular log event.
type LogEntry struct {
	Level   LogLevel `json:"level"`
	Caller  string   `json:"caller"`
	Path    string   `json:"path"`
	File    string   `json:"file"`
	Line    int      `json:"line"`
	Message string   `json:"message"`
}
