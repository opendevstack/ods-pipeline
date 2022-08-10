package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Run(executable string, args []string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}

func RunInDir(executable string, args []string, wsDir string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = wsDir
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}

func RunWithExtraEnvs(executable string, args []string, extraEnvs []string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	newEnv := append(os.Environ(), extraEnvs...)
	cmd.Env = newEnv
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}

// RunWithStreamingOutput invokes exe with given args and env. Stdout and stderr
// are streamed to outWriter and errWriter, respectively. If exe errors with an
// exit code equal to failureExitCode, no error is returned to the caller,
// but success is false. If exe does not error, success is true.
func RunWithStreamingOutput(exe string, args []string, env []string, outWriter, errWriter io.Writer, failureExitCode int) (success bool, err error) {
	return runWithStreamingOutput(exe, args, env, outWriter, errWriter, failureExitCode, false)
}

// RunWithStreamingOutput invokes exe with given args and env. Stdout and stderr
// are streamed to outWriter and errWriter, respectively. If exe errors with an
// exit code equal to failureExitCode, no error is returned to the caller,
// but success is false. If exe does not error, success is true.
// RunWithStreamingOutputReversed is identical to RunWithStreamingOutput,
// except for streaming stderr before stdout.
func RunWithStreamingOutputReversed(exe string, args []string, env []string, outWriter, errWriter io.Writer, failureExitCode int) (success bool, err error) {
	return runWithStreamingOutput(exe, args, env, outWriter, errWriter, failureExitCode, true)
}

// runWithStreamingOutput is the private implementation of RunWithStreamingOutput(Reversed).
func runWithStreamingOutput(exe string, args []string, env []string, outWriter, errWriter io.Writer, failureExitCode int, stderrFirst bool) (success bool, err error) {
	cmd := exec.Command(exe, args...)
	cmd.Env = append(os.Environ(), env...)
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return false, fmt.Errorf("connect stderr pipe: %w", err)
	}
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("connect stdout pipe: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return false, fmt.Errorf("start cmd: %w", err)
	}

	if stderrFirst {
		collectOutput(cmdStderr, errWriter)
		collectOutput(cmdStdout, outWriter)
	} else {
		collectOutput(cmdStdout, outWriter)
		collectOutput(cmdStderr, errWriter)
	}

	err = cmd.Wait()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() == failureExitCode {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func collectOutput(rc io.ReadCloser, w io.Writer) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
	}
}
