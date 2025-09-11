package goext

// This type holds the options for running an executable or command.
type runSettings struct {
	arguments             []string
	outputToConsole       bool
	workingDirectory      string
	skipPostProcessOutput bool
}

// Defines the type for options.
type runOption func(*runSettings)

// Appends the given arguments.
func RunWithArgs(arguments ...string) runOption {
	return func(options *runSettings) {
		options.arguments = append(options.arguments, arguments...)
	}
}

// Sets the working directory for the command.
func RunWithWorkingDirectory(workingDirectory string) runOption {
	return func(options *runSettings) {
		options.workingDirectory = workingDirectory
	}
}

// Allows enabling or disabling console output.
func RunWithConsoleOutput(outputToConsole bool) runOption {
	return func(options *runSettings) {
		options.outputToConsole = outputToConsole
	}
}

// Allows skipping the post-processing of output.
func RunWithSkipPostProcessOutput(skipPostProcessOutput bool) runOption {
	return func(options *runSettings) {
		options.skipPostProcessOutput = skipPostProcessOutput
	}
}
