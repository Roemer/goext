package goext

import (
	"errors"
	"testing"
)

func TestTernary(t *testing.T) {
	a := "Hello"
	b := "World"

	if got := Ternary(true, a, b); got != a {
		t.Errorf("Want %q, got %q", a, got)
	}
	if got := Ternary(false, a, b); got != b {
		t.Errorf("Want %q, got %q", b, got)
	}
}

func TestTernaryFunc(t *testing.T) {
	var isACalled, isBCalled bool
	aValue := "Hello"
	bValue := "World"
	a := func() string { isACalled = true; return aValue }
	b := func() string { isBCalled = true; return bValue }

	isACalled = false
	isBCalled = false
	if got := TernaryFunc(true, a, b); got != aValue {
		t.Errorf("Want %q, got %q", aValue, got)
	} else if isBCalled {
		t.Error("b should not be called when condition is true")
	}

	isACalled = false
	isBCalled = false
	if got := TernaryFunc(false, a, b); got != bValue {
		t.Errorf("Want %q, got %q", bValue, got)
	} else if isACalled {
		t.Error("a should not be called when condition is false")
	}
}

func TestTernaryFuncErr(t *testing.T) {
	var isACalled, isBCalled bool
	aValue := "Hello"
	bValue := "World"
	const aErrMsg = "a err"
	const bErrMsg = "b err"

	a := func() (string, error) { isACalled = true; return aValue, nil }
	b := func() (string, error) { isBCalled = true; return bValue, nil }
	aErr := func() (string, error) { isACalled = true; return aValue, errors.New(aErrMsg) }
	bErr := func() (string, error) { isBCalled = true; return bValue, errors.New(bErrMsg) }

	isACalled = false
	isBCalled = false
	if got, err := TernaryFuncErr(true, a, b); got != aValue {
		t.Errorf("Want %q, got %q", aValue, got)
	} else if isBCalled {
		t.Error("b should not be called when condition is true")
	} else if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	isACalled = false
	isBCalled = false
	if got, err := TernaryFuncErr(false, a, b); got != bValue {
		t.Errorf("Want %q, got %q", bValue, got)
	} else if isACalled {
		t.Error("a should not be called when condition is false")
	} else if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	isACalled = false
	isBCalled = false
	if _, err := TernaryFuncErr(true, aErr, bErr); err == nil {
		t.Error("Expected error but got nil")
	} else if isBCalled {
		t.Error("b should not be called when condition is true")
	} else if err.Error() != aErrMsg {
		t.Errorf("Expected error %q, got %q", aErrMsg, err.Error())
	}

	isACalled = false
	isBCalled = false
	if _, err := TernaryFuncErr(false, aErr, bErr); err == nil {
		t.Error("Expected error but got nil")
	} else if isACalled {
		t.Error("b should not be called when condition is true")
	} else if err.Error() != bErrMsg {
		t.Errorf("Expected error %q, got %q", bErrMsg, err.Error())
	}
}
