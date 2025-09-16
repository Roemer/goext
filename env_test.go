package goext

import (
	"os"
	"testing"
)

func TestEnvExists(t *testing.T) {
	if Env.Exists("HELLO_VAR") {
		t.Errorf("Expected env var HELLO_VAR to not exist")
	}
	os.Setenv("HELLO_VAR", "world")
	if !Env.Exists("HELLO_VAR") {
		t.Errorf("Expected env var HELLO_VAR to exist")
	}
	os.Unsetenv("HELLO_VAR")
}

func TestGetEnvOrDefault(t *testing.T) {
	if value, ok := Env.ValueOrDefault("HELLO_VAR", "default"); ok {
		t.Errorf("Expected env var HELLO_VAR to not exist, got %q", value)
	} else if value != "default" {
		t.Errorf("Expected default value to be 'default', got %q", value)
	}

	os.Setenv("HELLO_VAR", "world")
	if value, ok := Env.ValueOrDefault("HELLO_VAR", "default"); !ok {
		t.Errorf("Expected env var HELLO_VAR to exist")
	} else if value != "world" {
		t.Errorf("Expected value to be 'world', got %q", value)
	}

	os.Unsetenv("HELLO_VAR")
}
