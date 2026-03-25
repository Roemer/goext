package goext

import (
	"bytes"
	"io"
	"maps"
	"os"
	"os/exec"
)

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
	AdditionalEnv         map[string]string
	LogFilePath           string
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

	logFile, cleanup, err := r.openLogFile()
	if err != nil {
		return err
	}
	defer cleanup()

	var stdoutWriters, stderrWriters []io.Writer
	if r.OutputToConsole {
		stdoutWriters = append(stdoutWriters, os.Stdout)
		stderrWriters = append(stderrWriters, os.Stderr)
	}
	if logFile != nil {
		stdoutWriters = append(stdoutWriters, logFile)
		stderrWriters = append(stderrWriters, logFile)
	}
	cmd.Stdout = io.MultiWriter(stdoutWriters...)
	cmd.Stderr = io.MultiWriter(stderrWriters...)

	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// Runs the command and returns the separate output from stdout and stderr.
func (r *CmdRunner) RunGetOutput(executable string, arguments ...string) (string, string, error) {
	cmd := r.asCmd(executable, arguments...)

	logFile, cleanup, err := r.openLogFile()
	if err != nil {
		return "", "", err
	}
	defer cleanup()

	var stdoutBuf, stderrBuf bytes.Buffer
	var stdoutWriters, stderrWriters []io.Writer
	if r.OutputToConsole {
		stdoutWriters = append(stdoutWriters, os.Stdout)
		stderrWriters = append(stderrWriters, os.Stderr)
	}
	if logFile != nil {
		stdoutWriters = append(stdoutWriters, logFile)
		stderrWriters = append(stderrWriters, logFile)
	}
	stdoutWriters = append(stdoutWriters, &stdoutBuf)
	stderrWriters = append(stderrWriters, &stderrBuf)
	cmd.Stdout = io.MultiWriter(stdoutWriters...)
	cmd.Stderr = io.MultiWriter(stderrWriters...)

	err = cmd.Run()
	if r.SkipPostProcessOutput {
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	return r.processOutputString(stdoutBuf.String()), r.processOutputString(stderrBuf.String()), err
}

// Runs the command and returns the output from stdout and stderr combined.
func (r *CmdRunner) RunGetCombinedOutput(executable string, arguments ...string) (string, error) {
	cmd := r.asCmd(executable, arguments...)

	logFile, cleanup, err := r.openLogFile()
	if err != nil {
		return "", err
	}
	defer cleanup()

	var outBuf bytes.Buffer
	var writers []io.Writer
	if r.OutputToConsole {
		writers = append(writers, os.Stdout)
	}
	if logFile != nil {
		writers = append(writers, logFile)
	}
	writers = append(writers, &outBuf)
	combined := io.MultiWriter(writers...)
	cmd.Stdout = combined
	cmd.Stderr = combined

	err = cmd.Run()
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

// Enables output to console.
func (r *CmdRunner) WithConsoleOutput() *CmdRunner {
	return r.SetConsoleOutput(true)
}

// Sets output to console.
func (r *CmdRunner) SetConsoleOutput(outputToConsole bool) *CmdRunner {
	clone := r.Clone()
	clone.OutputToConsole = outputToConsole
	return clone
}

// Enables skipping post-processing of output (trimming newlines).
func (r *CmdRunner) WithSkipPostProcessOutput() *CmdRunner {
	return r.SetSkipPostProcessOutput(true)
}

// Sets skipping post-processing of output (trimming newlines).
func (r *CmdRunner) SetSkipPostProcessOutput(skipPostProcessOutput bool) *CmdRunner {
	clone := r.Clone()
	clone.SkipPostProcessOutput = skipPostProcessOutput
	return clone
}

// Adds an environment variable to the command.
func (r *CmdRunner) WithEnv(key, value string) *CmdRunner {
	return r.WithEnvs(map[string]string{key: value})
}

// Adds multiple environment variable to the command.
func (r *CmdRunner) WithEnvs(envs map[string]string) *CmdRunner {
	clone := r.Clone()
	maps.Copy(clone.AdditionalEnv, envs)
	return clone
}

// Sets a file path to which all command output (stdout + stderr) is written.
func (r *CmdRunner) WithLogFile(filePath string) *CmdRunner {
	clone := r.Clone()
	clone.LogFilePath = filePath
	return clone
}

// Clones the CmdRunner with its current configuration.
func (r *CmdRunner) Clone() *CmdRunner {
	clone := NewCmdRunner()
	clone.WorkingDirectory = r.WorkingDirectory
	clone.OutputToConsole = r.OutputToConsole
	clone.SkipPostProcessOutput = r.SkipPostProcessOutput
	clone.LogFilePath = r.LogFilePath
	clone.AdditionalEnv = make(map[string]string)
	maps.Copy(clone.AdditionalEnv, r.AdditionalEnv)
	return clone
}

////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////

func (r *CmdRunner) openLogFile() (*os.File, func(), error) {
	if r.LogFilePath == "" {
		return nil, func() {}, nil
	}
	logFile, err := os.OpenFile(r.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return nil, func() {}, err
	}
	return logFile, func() { _ = logFile.Close() }, nil
}

func (r *CmdRunner) asCmd(executable string, arguments ...string) *exec.Cmd {
	cmd := exec.Command(executable, arguments...)
	// Set the working directory
	if r.WorkingDirectory != "" {
		cmd.Dir = r.WorkingDirectory
	}
	// Set the additional environment variables
	if len(r.AdditionalEnv) > 0 {
		// Make sure to add the current environment variables first
		cmd.Env = os.Environ()
		// Then add the additional ones
		for k, v := range r.AdditionalEnv {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	return cmd
}

func (r *CmdRunner) processOutputString(value string) string {
	return StringTrimNewlineSuffix(value)
}
