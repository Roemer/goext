package goext

import (
	"bytes"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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

	stdoutWriter, stderrWriter, cleanup, err := r.prepareWriters(nil, nil)
	if err != nil {
		return err
	}
	defer cleanup()
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// Runs the command and returns the separate output from stdout and stderr.
func (r *CmdRunner) RunGetOutput(executable string, arguments ...string) (string, string, error) {
	cmd := r.asCmd(executable, arguments...)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutWriter, stderrWriter, cleanup, err := r.prepareWriters(&stdoutBuf, &stderrBuf)
	if err != nil {
		return "", "", err
	}
	defer cleanup()
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

	err = cmd.Run()
	if r.SkipPostProcessOutput {
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	return r.processOutputString(stdoutBuf.String()), r.processOutputString(stderrBuf.String()), err
}

// Runs the command and returns the output from stdout and stderr combined.
func (r *CmdRunner) RunGetCombinedOutput(executable string, arguments ...string) (string, error) {
	cmd := r.asCmd(executable, arguments...)

	var outBuf bytes.Buffer
	stdoutWriter, stderrWriter, cleanup, err := r.prepareWriters(&outBuf, &outBuf)
	if err != nil {
		return "", err
	}
	defer cleanup()
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

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

func (r *CmdRunner) asCmd(executable string, arguments ...string) *exec.Cmd {
	// Remove empty arguments that might cause issues on some platforms (e.g. Windows)
	arguments = slices.DeleteFunc(arguments, func(arg string) bool {
		return arg == ""
	})
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

func (r *CmdRunner) prepareWriters(stdoutBuf, stderrBuf *bytes.Buffer) (stdoutWriter, stderrWriter io.Writer, cleanup func(), err error) {
	cleanup = func() {}
	// Prepare the slices of the writers
	var stdoutWriters, stderrWriters []io.Writer
	// Add the console writers if needed
	if r.OutputToConsole {
		stdoutWriters = append(stdoutWriters, os.Stdout)
		stderrWriters = append(stderrWriters, os.Stderr)
	}
	// Add the log file writer if needed
	if r.LogFilePath != "" {
		if err := os.MkdirAll(filepath.Dir(r.LogFilePath), os.ModePerm); err != nil {
			return nil, nil, nil, err
		}
		logFile, err := os.OpenFile(r.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, nil, err
		}
		cleanup = func() {
			logFile.Close()
		}
		stdoutWriters = append(stdoutWriters, logFile)
		stderrWriters = append(stderrWriters, logFile)
	}
	// Add the buffer writers if needed
	if stdoutBuf != nil {
		stdoutWriters = append(stdoutWriters, stdoutBuf)
	}
	if stderrBuf != nil {
		stderrWriters = append(stderrWriters, stderrBuf)
	}
	// If no writers were added, add a dummy one to avoid nil writers
	if len(stdoutWriters) == 0 {
		stdoutWriters = append(stdoutWriters, io.Discard)
	}
	if len(stderrWriters) == 0 {
		stderrWriters = append(stderrWriters, io.Discard)
	}
	return io.MultiWriter(stdoutWriters...), io.MultiWriter(stderrWriters...), cleanup, nil
}
