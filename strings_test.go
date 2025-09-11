package goext

import (
	"reflect"
	"testing"
)

func TestStringContainsAny(t *testing.T) {
	a := "Hello"
	b := "World"

	if !StringContainsAny(a, a, b) {
		t.Errorf("Expected %q to contain any of %v", a, []string{a, b})
	}

	if !StringContainsAny(b, a, b) {
		t.Errorf("Expected %q to contain any of %v", b, []string{a, b})
	}

	if StringContainsAny(a, "xxx", "yyy") {
		t.Errorf("Expected false but was true")
	}
}

func TestStringTrimAllPrefix(t *testing.T) {
	expected := "yyyxxx"
	result := StringTrimAllPrefix("xxxyyyxxx", "x")
	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestStringTrimAllSuffix(t *testing.T) {
	expected := "xxxyyy"
	result := StringTrimAllSuffix("xxxyyyxxx", "x")
	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestStringTrimNewlineSuffix(t *testing.T) {
	expected := "my-text"
	result := StringTrimNewlineSuffix("my-text\n\r\n")
	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestStringSplitByNewLine(t *testing.T) {
	expected := []string{"line1", "line2", "line3"}
	result := StringSplitByNewLine("line1\r\nline2\nline3")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}
