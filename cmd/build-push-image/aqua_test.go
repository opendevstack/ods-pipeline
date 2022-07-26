package main

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/command"
)

// TestAquaScanError covers the cases when the aqua scanner can be invoked, but may succeed or fail.
func TestAquaScan(t *testing.T) {
	tests := map[string]struct {
		cmdStdout   string
		cmdStderr   string
		cmdExitCode int
		wantSuccess bool
		wantOut     string
		wantErr     bool
	}{
		"scan exits with license validation failure exit code": {
			cmdStdout:   "summary",
			cmdStderr:   "log output",
			cmdExitCode: scanLicenseValidationFailureExitCode,
			wantSuccess: false,
			wantOut:     "log output\n\nsummary\n",
			wantErr:     true,
		},
		"scan exits with compliance failure exit code": {
			cmdStdout:   "summary",
			cmdStderr:   "log output",
			cmdExitCode: scanComplianceFailureExitCode,
			wantSuccess: false,
			wantOut:     "log output\n\nsummary\n",
			wantErr:     false,
		},
		"scan passes": {
			cmdStdout:   "summary",
			cmdStderr:   "log output",
			cmdExitCode: 0,
			wantSuccess: true,
			wantOut:     "log output\n\nsummary\n",
			wantErr:     false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeAquaScan := aquaScanRunnerFunc(func(opts options, image, htmlReportFile, jsonReportFile string) ([]byte, []byte, error) {
				return command.Run("../../test/scripts/exit-with-code.sh", []string{tc.cmdStdout, tc.cmdStderr, strconv.Itoa(tc.cmdExitCode)})
			})
			success, out, err := aquaScan(fakeAquaScan, options{}, "image", "html", "json")
			if tc.wantErr && err == nil {
				t.Fatal("want err, got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("want no err, got %s", err)
			}
			if diff := cmp.Diff(tc.wantOut, out); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
			if tc.wantSuccess != success {
				t.Fatalf("want success=%v, got success=%v", tc.wantSuccess, success)
			}
		})
	}
}

// TestAquaScanError covers the case when the aqua scanner cannot be invoked at all.
func TestAquaScanError(t *testing.T) {
	fakeAquaScan := aquaScanRunnerFunc(func(opts options, image, htmlReportFile, jsonReportFile string) ([]byte, []byte, error) {
		return command.Run("bogus.sh", []string{})
	})
	success, out, err := aquaScan(fakeAquaScan, options{}, "image", "html", "json")
	if err == nil || err.Error() != "scan error: fork/exec bogus.sh: no such file or directory" {
		t.Fatal("want err, got none")
	}
	if success {
		t.Fatal("scan should not be successful")
	}
	if diff := cmp.Diff("\n", out); diff != "" {
		t.Fatalf("output mismatch (-want +got):\n%s", diff)
	}
}
