package goext

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"regexp"
)

var argumentsRegex = regexp.MustCompile(`[^\s"]+|"((\\"|[^"])*)"`)

type cmdRunners struct {
	// A default runner with no special options set.
	Default *CmdRunner
	// A runner that outputs to the console.
	Console *CmdRunner
}

// Contains a few pre-configured CmdRunners for easy access.
var CmdRunners cmdRunners = cmdRunners{
	Default: NewCmdRunner(),
	Console: NewCmdRunner().WithConsoleOutput(),
}

// The CmdRunner struct that holds the configuration for running commands.
type CmdRunner struct {
	WorkingDirectory      string
	OutputToConsole       bool
	SkipPostProcessOutput bool
	UseCurrentEnv         bool
	AdditionalEnv         map[string]string
}

// Creates a new CmdRunner with the given options.
func NewCmdRunner() *CmdRunner {
	cmdRunner := &CmdRunner{
		AdditionalEnv: make(map[string]string),
	}
	return cmdRunner
}

// Runs the command with the given options.
func (r *CmdRunner) Run(executable string, arguments ...string) error {
	cmd := r.asCmd(executable, arguments...)
	if r.OutputToConsole {
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

// Runs the command and returns the separate output from stdout and stderr.
func (r *CmdRunner) RunGetOutput(executable string, arguments ...string) (string, string, error) {
	cmd := r.asCmd(executable, arguments...)
	var stdoutBuf, stderrBuf bytes.Buffer
	if r.OutputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	} else {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	}
	err := cmd.Run()
	if r.SkipPostProcessOutput {
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	return r.processOutputString(stdoutBuf.String()), r.processOutputString(stderrBuf.String()), err
}

// Runs the command and returns the output from stdout and stderr combined.
func (r *CmdRunner) RunGetCombinedOutput(executable string, arguments ...string) (string, error) {
	cmd := r.asCmd(executable, arguments...)
	var outBuf bytes.Buffer
	if r.OutputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &outBuf)
	} else {
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf
	}
	err := cmd.Run()
	if r.SkipPostProcessOutput {
		return outBuf.String(), err
	}
	return r.processOutputString(outBuf.String()), err
}

// Sets the working directory for the command.
func (r *CmdRunner) WithWorkingDirectory(workingDirectory string) *CmdRunner {
	clone := r.Clone()
	clone.WorkingDirectory = workingDirectory
	return clone
}

// Sets to output to console.
func (r *CmdRunner) WithConsoleOutput() *CmdRunner {
	clone := r.Clone()
	clone.OutputToConsole = true
	return clone
}

// Sets to skip post-processing of output (trimming newlines).
func (r *CmdRunner) WithSkipPostProcessOutput() *CmdRunner {
	clone := r.Clone()
	clone.SkipPostProcessOutput = true
	return clone
}

// Adds the current environment variables to the command.
func (r *CmdRunner) WithCurrentEnvironment() *CmdRunner {
	clone := r.Clone()
	clone.UseCurrentEnv = true
	return clone
}

// Adds an environment variable to the command.
func (r *CmdRunner) WithEnv(key, value string) *CmdRunner {
	clone := r.Clone()
	if clone.AdditionalEnv == nil {
		clone.AdditionalEnv = make(map[string]string)
	}
	clone.AdditionalEnv[key] = value
	return clone
}

// Clones the CmdRunner with its current configuration.
func (r *CmdRunner) Clone() *CmdRunner {
	clone := NewCmdRunner()
	clone.WorkingDirectory = r.WorkingDirectory
	clone.OutputToConsole = r.OutputToConsole
	clone.SkipPostProcessOutput = r.SkipPostProcessOutput
	clone.UseCurrentEnv = r.UseCurrentEnv
	clone.AdditionalEnv = make(map[string]string)
	for k, v := range r.AdditionalEnv {
		clone.AdditionalEnv[k] = v
	}
	return clone
}

// Splits command line arguments into a slice of strings, respecting quoted substrings.
func SplitCmdArgs(arguments string) []string {
	return argumentsRegex.FindAllString(arguments, -1)
}

////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////

func (r *CmdRunner) asCmd(executable string, arguments ...string) *exec.Cmd {
	cmd := exec.Command(executable, arguments...)
	cmd.Dir = r.WorkingDirectory
	if r.UseCurrentEnv {
		cmd.Env = os.Environ()
	}
	for k, v := range r.AdditionalEnv {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	return cmd
}

func (r *CmdRunner) processOutputString(value string) string {
	return StringTrimNewlineSuffix(value)
}
