package main

import (
	"bytes"
	"strconv"
	"testing"
)

func TestAquaScan(t *testing.T) {
	tests := map[string]struct {
		cmdExitCode int
		wantSuccess bool
		wantErr     bool
	}{
		"scan exits with license validation failure exit code": {
			cmdExitCode: scanLicenseValidationFailureExitCode,
			wantSuccess: false,
			wantErr:     true,
		},
		"scan exits with compliance failure exit code": {
			cmdExitCode: scanComplianceFailureExitCode,
			wantSuccess: false,
			wantErr:     false,
		},
		"scan passes": {
			cmdExitCode: 0,
			wantSuccess: true,
			wantErr:     false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			success, err := aquaScan(
				"../../test/scripts/exit-with-code.sh",
				[]string{"", "", strconv.Itoa(tc.cmdExitCode)},
				&stdout, &stderr,
			)
			if tc.wantErr && err == nil {
				t.Fatal("want err, got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("want no err, got %s", err)
			}
			if tc.wantSuccess != success {
				t.Fatalf("want success=%v, got success=%v", tc.wantSuccess, success)
			}
		})
	}
}