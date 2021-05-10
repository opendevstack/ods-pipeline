package exchange

import "time"

type PipelineInfo struct {
	Time time.Time
	URL  string
}

type ItemInfo struct {
	Name        string
	Description string
}

type TestReport struct {
	Name        string
	Description string
	Executor    string
	Time        time.Time
	Suites      []TestSuite
	Properties  map[string]string // extension mechanims
	// extra prop for linking to e.q. requirement(s)?
}

type TestSuite struct {
	Name        string
	Description string
	Cases       []TestCase
	Duration    time.Duration
	Properties  map[string]string // extension mechanims
	// extra prop for linking to e.q. requirement(s)?
}

type TestCase struct {
	Name        string
	Description string
	Result      string // should be enum: Passed / Failed / Skipped
	Message     string
	Duration    time.Duration
	Properties  map[string]string // extension mechanims
	// extra prop for linking to e.q. requirement(s)?
}
