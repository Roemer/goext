package goext

import (
	"testing"
)

func TestCmdSplitArgs(t *testing.T) {
	args := Cmd.SplitArgs(`arg1 arg2 "arg 3" "arg \"4\""`)
	expected := []string{"arg1", "arg2", `"arg 3"`, `"arg \"4\""`}
	if len(args) != len(expected) {
		t.Fatalf("Expected %d args, got %d", len(expected), len(args))
	}

	for i, arg := range args {
		if arg != expected[i] {
			t.Errorf("Expected arg %d to be %q, got %q", i, expected[i], arg)
		}
	}
}

func TestCmdSplitArgsMultiple(t *testing.T) {
	args := Cmd.SplitArgs(`arg1 arg2 "arg 3" "arg \"4\""`, "arg77 arg88", "arg99")
	expected := []string{"arg1", "arg2", `"arg 3"`, `"arg \"4\""`, "arg77", "arg88", "arg99"}
	if len(args) != len(expected) {
		t.Fatalf("Expected %d args, got %d", len(expected), len(args))
	}

	for i, arg := range args {
		if arg != expected[i] {
			t.Errorf("Expected arg %d to be %q, got %q", i, expected[i], arg)
		}
	}
}
