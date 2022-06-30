package main

import (
	"bytes"
	"regexp"
)

const helmDiffDetectedMarker = `Error: identified at least one change, exiting with non-zero exit code (detailed-exitcode parameter enabled)`
const desiredDiffMessage = `plugin "diff" identified at least one change`

// cleanHelmDiffOutput removes error messages from the given Helm output.
// Those error messages are confusing, because they do not come from  an actual
// error, but from detecting drift between desired and current Helm state.
func cleanHelmDiffOutput(out []byte) []byte {
	if !bytes.Contains(out, []byte(helmDiffDetectedMarker)) {
		return out
	}
	cleanedOut := bytes.Replace(
		out, []byte(helmDiffDetectedMarker), []byte(desiredDiffMessage), -1,
	)
	r := regexp.MustCompile(`Error: plugin "(diff|secrets)" exited with error[\n]?`)
	cleanedOut = r.ReplaceAll(cleanedOut, []byte{})
	r = regexp.MustCompile(`helm.go:81: \[debug\] plugin "(diff|secrets)" exited with error[\n]?`)
	cleanedOut = r.ReplaceAll(cleanedOut, []byte{})
	return cleanedOut
}
