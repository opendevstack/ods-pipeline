package testoutput

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// TestEvent taken from `go doc test2json`, unfortunately, this struct is not
// exposed from any core go packages so we have to define it ourselves
type TestEvent struct {
	Time    time.Time `json:",omitempty"`
	Action  string
	Package string  `json:",omitempty"`
	Test    string  `json:",omitempty"`
	Elapsed float64 `json:",omitempty"`
	Output  string  `json:",omitempty"`
}

// DisentangleTestOutputs separates test outputs from a stream of interleaved
// test events into separate files per individual testcase
func DisentangleTestOutputs(testResultsPath string, in *os.File, out *os.File) {
	// holds a mapping from test to the file handles the output is redirected to
	var fileHandles = make(map[string]*os.File)
	decoder := json.NewDecoder(in)

	for {
		var event TestEvent
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// switch over fixed set of actions (c.f. `go doc test2json`)
		switch event.Action {
		case "run":
			filePath := filepath.Join(testResultsPath, event.Test+".log")
			dir := filepath.Dir(filePath)
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				fmt.Printf("Failed to create directory: %v", err)
			}

			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Failed to open output file: %v\n", err)
			}
			fileHandles[event.Test] = f
		case "pause", "cont", "skip": // do nothing
		case "pass", "fail":
			f, ok := fileHandles[event.Test]
			if ok {
				fmt.Println("fileHandles: ", fileHandles)
				err := f.Close()
				if err != nil {
					fmt.Printf("Could not close file handle: %v\n", err)
				}
			}
		case "bench", "output":
			// 1. redirect to file
			f, ok := fileHandles[event.Test]
			if ok {
				if _, err := fmt.Fprint(f, event.Output); err != nil {
					fmt.Printf("Failed to write to output file: %v\n", err)
				}
			}
			// 2. echo the output with the json stripped off to `out`
			if _, err := fmt.Fprint(out, event.Test, event.Output); err != nil {
				fmt.Printf("Failed to write to output stream: %v", err)
			}
		}
	}
}
