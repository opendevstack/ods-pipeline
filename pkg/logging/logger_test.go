package logging

import (
	"bytes"
	"testing"
	"time"
)

// fixedClock is frozen in time.
type fixedClock struct{}

// Now always returns the same time.
func (fixedClock) Now() time.Time {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	t, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
	return t
}

func TestLogger(t *testing.T) {
	var logger LeveledLoggerInterface
	tests := map[string]struct {
		timestamp bool
		tag       string
		level     Level
		want      string
	}{
		"without time": {
			level: LevelDebug,
			want:  "INFO  | bar",
		},
		"with time": {
			timestamp: true,
			level:     LevelInfo,
			want:      "2013-02-03T19:54:00Z | INFO  | bar",
		},
		"higher level": {
			level: LevelError,
			want:  "",
		},
		"with tag": {
			tag:   "foo",
			level: LevelInfo,
			want:  "INFO  | foo: bar",
		},
		"with time and tag": {
			timestamp: true,
			tag:       "foo",
			level:     LevelInfo,
			want:      "2013-02-03T19:54:00Z | INFO  | foo: bar",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte{})
			logger = &LeveledLogger{
				StdoutOverride: b,
				Tag:            tc.tag,
				Timestamp:      tc.timestamp,
				ClockOverride:  fixedClock{},
				Level:          tc.level,
			}
			want := tc.want
			if tc.want != "" {
				want = want + "\n"
			}
			logger.Infof("bar")
			got := b.String()
			if got != want {
				t.Fatalf("want: %s, got: %s", want, got)
			}
		})
	}
}

func TestWithTag(t *testing.T) {
	var logger LeveledLoggerInterface
	stdout := bytes.NewBuffer([]byte("stdout\n"))
	stderr := bytes.NewBuffer([]byte("stderr\n"))
	logger = &LeveledLogger{
		StdoutOverride: stdout,
		StderrOverride: stderr,
		Tag:            "initial",
		Timestamp:      true,
		ClockOverride:  fixedClock{},
		Level:          LevelDebug,
	}
	taggedLogger := logger.WithTag("tagged")
	taggedLogger.Errorf("oops")
	wantStdout := "stdout\n"
	if stdout.String() != wantStdout {
		t.Fatalf("wamt %s, got: %s", wantStdout, stdout.String())
	}
	wantStderr := "stderr\n2013-02-03T19:54:00Z | ERROR | tagged: oops\n"
	if stderr.String() != wantStderr {
		t.Fatalf("wamt %s, got: %s", wantStderr, stderr.String())
	}
}
