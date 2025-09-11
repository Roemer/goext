package goext

import "testing"

func TestRun(t *testing.T) {
	err := Run("where", RunWithArgs("where"), RunWithConsoleOutput(true))
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestRunGetOutput(t *testing.T) {
	// TODO: For linux, use whereis: /usr/bin/whereis
	stdout, stderr, err := RunGetOutput("where", RunWithArgs("where"))
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
