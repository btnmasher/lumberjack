package lumberjack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/http"
)

//HttpClientBackend is an object that holds the configuration data
//to be used for implementing an instance of an HTTP POST logging
//backend that sends LogEntry messages via JSON to a specified URL.
//
//The exported Stop channel should be used during cleanup code
//to close down the internal Goroutine of the HttpClientBackend.
type HttpClientBackend struct {
	logchan chan LogEntry
	Stop    chan struct{}
	timer   *time.Ticker
	//TODO: Add more options like cookie, client certificate, basic auth, etc.
}

//logbuffer is a structure to be used to encapsulate LogEntry objects
//in an slice so that when Marshalled into JSON will be contained
//within a JSON array.
type logbuffer struct {
	Entries []LogEntry `json:"logentries"`
}

//NewHttpClientBackend is a function that accepts the url string, LogEntry buffer size
//and interval time.Duration to be used to configure and start a new instnace of
//HttpClientBackend with the given arguments. It will start a Goroutine that will
//act as a method of decoupling the blocking nature of HTTP requests from any
//application that sends a LogEntry to this Backend.
//
//It also implements a method of buffering LogEntry messages to pipeline in single
//HTTP POST requests. The buffer will empty into an HTTP POST request when it is
//full, or when the specified time.Duration interval passes on a time.Ticker. This
//will keep slowly moving logs from sitting too long in the buffer. If no interval
//is specified, a default of 1 second will be chosen. If no bufsize is specified,
//each LogEntry will be sent via HTTP POST individually.
func NewHttpClientBackend(url string, bufsize int, interval time.Duration) *HttpClientBackend {
	if interval == 0 {
		interval = time.Second * 1 //Default to 1 second interval incase they decided to be a poo-head and not set it.
	}

	h := HttpClientBackend{
		logchan: make(chan LogEntry, 50),  //Some breathing room to keep from blocking
		Stop:    make(chan struct{}),      //So we can kill our goroutine cleanly, implementer must close(h.Stop)
		timer:   time.NewTicker(interval), //how often we want to clear the buffer if not full.
	}

	go startClient(url, bufsize, h.logchan, h.timer, h.Stop)

	return &h
}

//startClient is an internal function used by the NewHttpClientBackend function to start up
//the Goroutine that will be ultimately handling the buffered LogEntry messages and
//sending via HTTP POST as JSON.
func startClient(url string, bufsize int, logchan chan LogEntry, timer *time.Ticker, stop chan struct{}) {
	var buffer logbuffer

	defer timer.Stop()

	for {
		select {
		case entry := <-logchan:
			buffer.Entries = append(buffer.Entries, entry)

			if bufsize > 0 { //Are we even trying to buffer requests?
				if len(buffer.Entries) < bufsize { //Have we filled the buffer yet?
					continue //Nope, just skip back to the start of the select loop
				}
			}

			err := doSend(url, buffer) //Send that buffer!
			if err != nil {
				logInternal(ERROR, err)
			}
			buffer.Entries = buffer.Entries[:0] //Clear that buffer!

		case <-timer.C:
			if len(buffer.Entries) > 0 {
				err := doSend(url, buffer)
				if err != nil {
					logInternal(ERROR, err)
				}
				buffer.Entries = buffer.Entries[:0] //Time's up, send what we have!
			}

		case <-stop:
			break
		}
	}
}

//doSend is an internal function that accepts a url and a logbuffer object that
//contains LogEntry objects to be Marshalled to JSON then sent via HTTP POST
//to the specified url. It returns an error if the http reqeust fails.
func doSend(url string, buffer logbuffer) error {
	data, err := json.Marshal(buffer)
	if err != nil {
		return fmt.Errorf("HTTP Backend: unable to Marshal JSON from logbuffer struct: %s", err)
	}

	b := bytes.NewBuffer(data)

	err = http.Post(url, b)
	if err != nil {
		return fmt.Errorf("HTTP Backend: unable to POST to specified URL, library returned error: %s", err)
	}
	return nil
}

//Log implements the Backend interface's requirements and will send LogEntry
//object references to the channel on the current HttpClientBackend to be
//buffered then sent via HTTP POST as JSON.
func (h *HttpClientBackend) Log(entry *LogEntry) {
	h.logchan <- *entry
}
