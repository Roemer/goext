package goext

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type CmdRunner struct {
	Executable            string
	Arguments             []string
	WorkingDirectory      string
	OutputToConsole       bool
	SkipPostProcessOutput bool
	UseCurrentEnv         bool
	AdditionalEnv         map[string]string
}

// Creates a new CmdRunner with the given executable and arguments.
func NewCmdRunner(executable string, arguments ...string) *CmdRunner {
	return &CmdRunner{
		Executable:    executable,
		Arguments:     arguments,
		AdditionalEnv: make(map[string]string),
	}
}

// Runs the command with the given options.
func (r *CmdRunner) Run() error {
	cmd := r.asCmd()
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
func (r *CmdRunner) RunGetOutput() (string, string, error) {
	cmd := r.asCmd()
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
func (r *CmdRunner) RunGetCombinedOutput() (string, error) {
	cmd := r.asCmd()
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
	r.WorkingDirectory = workingDirectory
	return r
}

// Sets to output to console.
func (r *CmdRunner) WithConsoleOutput() *CmdRunner {
	r.OutputToConsole = true
	return r
}

// Sets to skip post-processing of output (trimming newlines).
func (r *CmdRunner) WithSkipPostProcessOutput() *CmdRunner {
	r.SkipPostProcessOutput = true
	return r
}

// Adds the current environment variables to the command.
func (r *CmdRunner) WithCurrentEnvironment() *CmdRunner {
	r.UseCurrentEnv = true
	return r
}

// Adds an environment variable to the command.
func (r *CmdRunner) WithEnv(key, value string) *CmdRunner {
	if r.AdditionalEnv == nil {
		r.AdditionalEnv = make(map[string]string)
	}
	r.AdditionalEnv[key] = value
	return r
}

////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////

func (r *CmdRunner) asCmd() *exec.Cmd {
	cmd := exec.Command(r.Executable, r.Arguments...)
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
