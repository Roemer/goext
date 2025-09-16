package goext

import (
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	err := NewCmdRunner().WithConsoleOutput().Run("where", "where")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestRunGetOutput(t *testing.T) {
	// TODO: For linux, use whereis: /usr/bin/whereis
	stdout, stderr, err := NewCmdRunner().RunGetOutput("where", "where")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if stdout != `C:\Windows\System32\where.exe` {
		t.Errorf("Expected output to be %q but got %q", `C:\Windows\System32\where.exe`, stdout)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr but got %q", stderr)
	}
}

func TestCmdRunnerWith(t *testing.T) {
	os.MkdirAll("test", os.ModePerm)
	defer os.RemoveAll("test")
	runner := NewCmdRunner().WithWorkingDirectory("test").WithEnv("HELLO_VAR", "world")

	pwd, err := runner.RunGetCombinedOutput("cmd", "/C", "cd")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if !strings.HasSuffix(pwd, `\test`) {
		t.Errorf("Expected working directory to end with %q but got %q", `\test`, pwd)
	}

	output, err := runner.RunGetCombinedOutput("cmd", "/C", "echo %HELLO_VAR%")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if output != "world" {
		t.Errorf("Expected output to be %q but got %q", "world", output)
	}
}
