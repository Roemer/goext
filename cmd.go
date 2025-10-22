package goext

import (
	"errors"
	"os/exec"
	"regexp"
)

var argumentsRegex = regexp.MustCompile(`[^\s"]+|"((\\"|[^"])*)"`)

type cmdNamespace int

// Contains methods regarding commands.
var Cmd cmdNamespace = 0

// Splits command line arguments into a slice of strings, respecting quoted substrings.
func (cmdNamespace) SplitArgs(arguments ...string) []string {
	finalArgs := []string{}
	for _, args := range arguments {
		finalArgs = append(finalArgs, argumentsRegex.FindAllString(args, -1)...)
	}
	return finalArgs
}

// Gets the exit code from a command error.
func (cmdNamespace) ErrorExitCode(err error) int {
	if err == nil {
		// No error
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		// The commands exit code
		return exitErr.ExitCode()
	}
	// Some other error (e.g. command not found)
	return -1
}
