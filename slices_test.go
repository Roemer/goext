package goext

import "testing"

func TestSliceAppendIf(t *testing.T) {
	testSlice := []int{1, 2, 3}

	testSlice = SliceAppendIf(testSlice, true, 4, 5)
	if len(testSlice) != 5 || testSlice[3] != 4 || testSlice[4] != 5 {
		t.Errorf("Expected slice to be [1, 2, 3, 4, 5], got %v", testSlice)
	}

	testSlice = SliceAppendIf(testSlice, false, 5, 6)
	if len(testSlice) != 5 || testSlice[3] != 4 || testSlice[4] != 5 {
		t.Errorf("Expected slice to remain [1, 2, 3, 4, 5], got %v", testSlice)
	}
}

func TestSliceAppendIfFunc(t *testing.T) {
	f1 := func() []int { return []int{4, 5} }
	f2 := func() []int { return []int{6, 7} }
	testSlice := []int{1, 2, 3}

	testSlice = SliceAppendIfFunc(testSlice, true, f1)
	if len(testSlice) != 5 || testSlice[3] != 4 || testSlice[4] != 5 {
		t.Errorf("Expected slice to be [1, 2, 3, 4, 5], got %v", testSlice)
	}

	testSlice = SliceAppendIfFunc(testSlice, false, f2)
	if len(testSlice) != 5 || testSlice[3] != 4 || testSlice[4] != 5 {
		t.Errorf("Expected slice to remain [1, 2, 3, 4, 5], got %v", testSlice)
	}
}

func TestSlicePrepend(t *testing.T) {
	testSlice := []int{1, 2, 3}

	testSlice = SlicePrepend(testSlice, 4, 5)
	if len(testSlice) != 5 || testSlice[0] != 4 || testSlice[1] != 5 || testSlice[2] != 1 || testSlice[3] != 2 || testSlice[4] != 3 {
		t.Errorf("Expected slice to be [4, 5, 1, 2, 3], got %v", testSlice)
	}
}

func TestSlicePrependIf(t *testing.T) {
	testSlice := []int{1, 2, 3}

	testSlice = SlicePrependIf(testSlice, true, 4, 5)
	if len(testSlice) != 5 || testSlice[0] != 4 || testSlice[1] != 5 || testSlice[2] != 1 || testSlice[3] != 2 || testSlice[4] != 3 {
		t.Errorf("Expected slice to be [4, 5, 1, 2, 3], got %v", testSlice)
	}

	testSlice = SlicePrependIf(testSlice, false, 5, 6)
	if len(testSlice) != 5 || testSlice[0] != 4 || testSlice[1] != 5 || testSlice[2] != 1 || testSlice[3] != 2 || testSlice[4] != 3 {
		t.Errorf("Expected slice to remain [4, 5, 1, 2, 3], got %v", testSlice)
	}
}

func TestSlicePrependIfFunc(t *testing.T) {
	f1 := func() []int { return []int{4, 5} }
	f2 := func() []int { return []int{6, 7} }
	testSlice := []int{1, 2, 3}

	testSlice = SlicePrependIfFunc(testSlice, true, f1)
	if len(testSlice) != 5 || testSlice[0] != 4 || testSlice[1] != 5 || testSlice[2] != 1 || testSlice[3] != 2 || testSlice[4] != 3 {
		t.Errorf("Expected slice to be [4, 5, 1, 2, 3], got %v", testSlice)
	}

	testSlice = SlicePrependIfFunc(testSlice, false, f2)
	if len(testSlice) != 5 || testSlice[0] != 4 || testSlice[1] != 5 || testSlice[2] != 1 || testSlice[3] != 2 || testSlice[4] != 3 {
		t.Errorf("Expected slice to remain [4, 5, 1, 2, 3], got %v", testSlice)
	}
}
