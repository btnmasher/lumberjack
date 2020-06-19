package lumberjack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// Test data to use.
var testobj logbuffer = logbuffer{
	Entries: []LogEntry{
		{
			Level:   ERROR,
			Caller:  "main.main()",
			Path:    "/somewhere",
			File:    "main.go",
			Line:    10,
			Message: "Test Error",
		},
		{
			Level:   INFO,
			Caller:  "main.main()",
			Path:    "/somewhere/else",
			File:    "main.go",
			Line:    11,
			Message: "Test Info",
		},
	},
}

func TestHttpBackendPost(t *testing.T) {

	// Test server that always responds with 200 code.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		// Instantiate an empty struct to use to read the POST data.
		b := logbuffer{}

		// Get the POST bytes.
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		// Attempt to marshal the data sent to us back into a struct.
		if err := json.Unmarshal(bytes, &b); err != nil {
			t.Error(err)
		}

		// There should be 2 entries
		expect(t, len(b.Entries), 2)

		// Unmarshalled entries should be equivilant to the original testobj entries
		expect(t, b.Entries[0], testobj.Entries[0])
		expect(t, b.Entries[1], testobj.Entries[1])
	}))

	defer server.Close()

	// Test the method!
	err := doSend(server.URL, testobj)
	if err != nil {
		t.Error(err)
	}
}
