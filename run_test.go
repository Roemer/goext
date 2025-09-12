package goext

import (
	"os"
	"testing"
)

func TestRunWithEnv(t *testing.T) {
	os.Setenv("PRE_VAR_1", "baz")
	os.Setenv("PRE_VAR_2", "hello")
	assertEnvEquals(t, "PRE_VAR_1", "baz")
	assertEnvEquals(t, "PRE_VAR_2", "hello")
	assertEnvIsUnset(t, "VAR_1")
	assertEnvIsUnset(t, "VAR_2")

	err := RunWithOptions(func() error {
		assertEnvEquals(t, "PRE_VAR_1", "baz")
		assertEnvEquals(t, "PRE_VAR_2", "world")
		assertEnvEquals(t, "VAR_1", "foo")
		assertEnvEquals(t, "VAR_2", "bar")
		return nil
	}, RunOptionWithEnvs(map[string]string{
		"VAR_1":     "foo",
		"VAR_2":     "bar",
		"PRE_VAR_2": "world",
	}))

	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	assertEnvEquals(t, "PRE_VAR_1", "baz")
	assertEnvEquals(t, "PRE_VAR_2", "hello")
	assertEnvIsUnset(t, "VAR_1")
	assertEnvIsUnset(t, "VAR_2")
}

func assertEnvEquals(t *testing.T, name string, expected string) {
	actual, ok := os.LookupEnv(name)
	if !ok {
		t.Errorf("Expected env var %q to be set to %q but it is not set", name, expected)
	}
	if actual != expected {
		t.Errorf("Expected env var %q to be %q but got %q", name, expected, actual)
	}
}

func assertEnvIsUnset(t *testing.T, name string) {
	_, ok := os.LookupEnv(name)
	if ok {
		t.Errorf("Expected env var %q to be unset but it is set", name)
	}
}
