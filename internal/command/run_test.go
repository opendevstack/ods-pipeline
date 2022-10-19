package command

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestRunWithStreamingOutput covers the cases when helm-diff can be invoked, but may detect drift or not.
func TestRunWithStreamingOutput(t *testing.T) {
	tests := map[string]struct {
		cmdExitCode int
		wantSuccess bool
		wantErr     bool
	}{
		"cmd exits with generic exit code": {
			cmdExitCode: 1,
			wantSuccess: false,
			wantErr:     true,
		},
		"cmd exits with special failure exit code": {
			cmdExitCode: 2,
			wantSuccess: false,
			wantErr:     false,
		},
		"cmd finishes without error": {
			cmdExitCode: 0,
			wantSuccess: true,
			wantErr:     false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			driftDetected, err := RunWithSpecialFailureCode(
				"../../test/scripts/exit-with-code.sh",
				[]string{"value of FOO=${FOO}", "log msg", strconv.Itoa(tc.cmdExitCode)},
				[]string{"FOO=bar"},
				&stdout, &stderr,
				2,
			)
			if tc.wantErr && err == nil {
				t.Fatal("want err, got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("want no err, got %s", err)
			}
			wantStdout := "value of FOO=bar"
			wantStderr := "log msg"
			if diff := cmp.Diff(wantStdout+"\n", stdout.String()); diff != "" {
				t.Fatalf("stdout mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(wantStderr+"\n", stderr.String()); diff != "" {
				t.Fatalf("stderr mismatch (-want +got):\n%s", diff)
			}
			if tc.wantSuccess != driftDetected {
				t.Fatalf("want success=%v, got success=%v", tc.wantSuccess, driftDetected)
			}
		})
	}
}

// TestRunWithStreamingOutputError covers the case when the aqua scanner cannot be invoked at all.
func TestRunWithStreamingOutputError(t *testing.T) {
	success, err := RunWithSpecialFailureCode(
		"./bogus.sh", []string{}, []string{}, &bytes.Buffer{}, &bytes.Buffer{}, -1,
	)
	if err == nil || err.Error() != "start cmd: fork/exec ./bogus.sh: no such file or directory" {
		t.Fatalf("want err, got: %+v", err)
	}
	if success {
		t.Fatal("cmd should not be successful")
	}
}

func TestInterleavedStdoutAndStderr(t *testing.T) {
	var out bytes.Buffer
	success, err := RunWithSpecialFailureCode(
		"../../test/scripts/interleaved-output.sh", []string{}, []string{}, &out, &out, -1,
	)
	if err != nil {
		t.Fatal(err)
	}
	wantOut := "some stdout\nsome stderr\nmore stdout\nmore stderr\nstderr after sleep\nstdout after sleep"
	if diff := cmp.Diff(wantOut+"\n", out.String()); diff != "" {
		t.Fatalf("stdout mismatch (-want +got):\n%s", diff)
	}
	if !success {
		t.Fatal("cmd should be successful")
	}

}
