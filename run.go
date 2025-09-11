package goext

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// Creates a new exec.Cmd with the given executable and arguments.
func NewCmd(executable string, arguments ...string) *exec.Cmd {
	return exec.Command(executable, arguments...)
}

// Runs an executable with the given options.
func Run(executable string, options ...runOption) error {
	cmd := NewCmd(executable)
	return RunCommand(cmd, options...)
}

// Runs a command with the given options.
func RunCommand(cmd *exec.Cmd, options ...runOption) error {
	runOptions := prepare(cmd, options...)
	if runOptions.outputToConsole {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

// Runs an executable with the given options and returns the separate output from stdout and stderr.
func RunGetOutput(executable string, options ...runOption) (string, string, error) {
	cmd := NewCmd(executable)
	return RunCommandGetOutput(cmd, options...)
}

// Runs a command with the given options and returns the separate output from stdout and stderr.
func RunCommandGetOutput(cmd *exec.Cmd, options ...runOption) (string, string, error) {
	runOptions := prepare(cmd, options...)
	var stdoutBuf, stderrBuf bytes.Buffer
	if runOptions.outputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	} else {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	}
	err := cmd.Run()
	if runOptions.skipPostProcessOutput {
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	return processOutputString(stdoutBuf.String()), processOutputString(stderrBuf.String()), err
}

// Runs an executable with the given options and returns the output from stdout and stderr combined.
func RunGetCombinedOutput(executable string, options ...runOption) (string, error) {
	cmd := NewCmd(executable)
	return RunCommandGetCombinedOutput(cmd, options...)
}

// Runs a command with the given options and returns the output from stdout and stderr combined.
func RunCommandGetCombinedOutput(cmd *exec.Cmd, options ...runOption) (string, error) {
	runOptions := prepare(cmd, options...)
	var outBuf bytes.Buffer
	if runOptions.outputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &outBuf)
	} else {
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf
	}
	err := cmd.Run()
	if runOptions.skipPostProcessOutput {
		return outBuf.String(), err
	}
	return processOutputString(outBuf.String()), err
}

////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////

// Applies the options to a command.
func prepare(cmd *exec.Cmd, options ...runOption) *runSettings {
	// Build the options
	runOptions := &runSettings{}
	for _, option := range options {
		option(runOptions)
	}
	// Arguments
	if len(runOptions.arguments) > 0 {
		cmd.Args = append(cmd.Args, runOptions.arguments...)
	}
	// Working directory
	if runOptions.workingDirectory != "" {
		cmd.Dir = runOptions.workingDirectory
	}
	return runOptions
}

func processOutputString(value string) string {
	return StringTrimNewlineSuffix(value)
}
