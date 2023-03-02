package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

func RunBuffered(executable string, args []string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}

func RunBufferedInDir(executable string, args []string, wsDir string) (outBytes, errBytes []byte, err error) {
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

// Run invokes exe with given args and env. Stdout and stderr
// are streamed to outWriter and errWriter, respectively.
func Run(exe string, args []string, env []string, outWriter, errWriter io.Writer) error {
	return RunInDir(exe, args, env, "", outWriter, errWriter)
}

// Run invokes exe with given args and env. Stdout and stderr
// are streamed to outWriter and errWriter, respectively.
// If dir is non-empty, the workdir of exe will be set to it.
func RunInDir(exe string, args []string, env []string, dir string, outWriter, errWriter io.Writer) error {
	cmd := exec.Command(exe, args...)
	cmd.Env = append(os.Environ(), env...)
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("connect stderr pipe: %w", err)
	}
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("connect stdout pipe: %w", err)
	}
	if dir != "" {
		cmd.Dir = dir
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("start cmd: %w", err)
	}

	err = collectOutput(cmdStdout, cmdStderr, outWriter, errWriter)
	if err != nil {
		return fmt.Errorf("collect output: %w", err)
	}

	return cmd.Wait()
}

// RunWithSpecialFailureCode invokes exe with given args and env. Stdout and stderr
// are streamed to outWriter and errWriter, respectively. If exe errors with an
// exit code equal to failureExitCode, no error is returned to the caller,
// but success is false. If exe does not error, success is true.
func RunWithSpecialFailureCode(exe string, args []string, env []string, outWriter, errWriter io.Writer, failureExitCode int) (success bool, err error) {
	err = Run(exe, args, env, outWriter, errWriter)
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() == failureExitCode {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func collectOutput(rcStdout, rcStderr io.ReadCloser, wStdout, wStderr io.Writer) error {
	var stdoutErr, stderrErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdoutErr = scan(rcStdout, wStdout)
		wg.Done()
	}()
	stderrErr = scan(rcStderr, wStderr)
	wg.Wait()
	if stdoutErr != nil || stderrErr != nil {
		return fmt.Errorf("scan stdout = %s, scan stderr = %s", stdoutErr, stderrErr)
	}
	return nil
}

func scan(rc io.ReadCloser, w io.Writer) error {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
	}
	return scanner.Err()
}
