package goext

import "testing"

func TestStringContainsAny(t *testing.T) {
	a := "Hello"
	b := "World"

	if !StringContainsAny(a, []string{a, b}) {
		t.Errorf("Expected %q to contain any of %v", a, []string{a, b})
	}

	if !StringContainsAny(b, []string{a, b}) {
		t.Errorf("Expected %q to contain any of %v", b, []string{a, b})
	}

	if StringContainsAny(a, []string{"xxx", "yyy"}) {
		t.Errorf("Expected false but was true")
	}
}
