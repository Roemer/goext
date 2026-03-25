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

func TestCmdRunnerWithSystemAndSpecificEnv(t *testing.T) {
	os.Setenv("SYTEM_VAR", "hello")
	runner := NewCmdRunner().WithEnv("HELLO_VAR", "world")
	output, err := runner.RunGetCombinedOutput("cmd", "/C", `echo %SYTEM_VAR% %HELLO_VAR%`)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if output != "hello world" {
		t.Errorf("Expected output to be %q but got %q", "hello world", output)
	}
}

func TestCmdRunnerWithSystemEnvOnly(t *testing.T) {
	os.Setenv("SYTEM_VAR", "hello world")
	runner := NewCmdRunner()
	output, err := runner.RunGetCombinedOutput("cmd", "/C", `echo %SYTEM_VAR%`)
	//output, err := runner.RunGetCombinedOutput("powershell", "-Command", "Write-Host $env:SYTEM_VAR")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if output != "hello world" {
		t.Errorf("Expected output to be %q but got %q", "hello world", output)
	}
}

func TestCmdRunnerWithFileOutput(t *testing.T) {
	logFilePath := "test_output.log"
	defer os.Remove(logFilePath)
	runner := NewCmdRunner().WithLogFile(logFilePath)
	err := runner.Run("cmd", "/C", "echo Hello File Output")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	logContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if strings.TrimSpace(string(logContent)) != strings.TrimSpace("Hello File Output") {
		t.Errorf("Expected log content to be %q but got %q", "Hello File Output", string(logContent))
	}
}

func TestCmdRunnerWithFileAndCombinedOutput(t *testing.T) {
	logFilePath := "test_output.log"
	defer os.Remove(logFilePath)
	runner := NewCmdRunner().WithLogFile(logFilePath)
	output, err := runner.RunGetCombinedOutput("cmd", "/C", "echo Hello File Output")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if output != "Hello File Output" {
		t.Errorf("Expected output to be %q but got %q", "Hello File Output", output)
	}
	logContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if strings.TrimSpace(string(logContent)) != strings.TrimSpace("Hello File Output") {
		t.Errorf("Expected log content to be %q but got %q", "Hello File Output", string(logContent))
	}
}

func TestCmdRunnerWithFileAndSeparateOutput(t *testing.T) {
	logFilePath := "test_output.log"
	defer os.Remove(logFilePath)
	runner := NewCmdRunner().WithLogFile(logFilePath)
	stdout, stderr, err := runner.RunGetOutput("cmd", "/C", "echo Hello File Output")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if stdout != "Hello File Output" {
		t.Errorf("Expected stdout to be %q but got %q", "Hello File Output", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr but got %q", stderr)
	}
	logContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if strings.TrimSpace(string(logContent)) != strings.TrimSpace("Hello File Output") {
		t.Errorf("Expected log content to be %q but got %q", "Hello File Output", string(logContent))
	}
}
