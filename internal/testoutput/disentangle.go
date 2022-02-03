package testoutput

import (
	"encoding/json"
	"fmt"
	"io"
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
func DisentangleTestOutputs(testResultsPath string, in *os.File, out *os.File) error {
	// holds a mapping from test to the file handles the output is redirected to
	var fileHandles = make(map[string]*os.File)
	decoder := json.NewDecoder(in)

	for {
		var event TestEvent
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// switch over fixed set of actions (c.f. `go doc test2json`)
		switch event.Action {
		case "run":
			filePath := filepath.Join(testResultsPath, event.Test+".log")
			dir := filepath.Dir(filePath)
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("failed to open output file: %w", err)
			}
			fileHandles[event.Test] = f
		case "pause", "cont", "skip": // do nothing
		case "pass", "fail":
			f, ok := fileHandles[event.Test]
			if ok {
				err := f.Close()
				if err != nil {
					return fmt.Errorf("could not close file handle: %w", err)
				}
			}
		case "bench", "output":
			// 1. redirect to file
			f, ok := fileHandles[event.Test]
			if ok {
				if _, err := fmt.Fprint(f, event.Output); err != nil {
					return fmt.Errorf("failed to write to output file: %w", err)
				}
			}
			// 2. echo the output with the json stripped off to `out`
			if _, err := fmt.Fprint(out, event.Output); err != nil {
				return fmt.Errorf("failed to write to output stream: %w", err)
			}
		}
	}
	return nil
}
