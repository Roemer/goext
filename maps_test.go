package goext

import (
	"fmt"
	"testing"
)

func TestMapSortedByKey(t *testing.T) {
	myMap := map[int]string{
		1: "b",
		0: "a",
		3: "d",
		2: "c",
	}

	index := 0
	for k := range MapSortedByKey(myMap) {
		if k != index {
			t.Errorf("Expected key %d, got %d", index, k)
		}
		index++
	}

	myMap2 := map[int]string{1: "b", 0: "a", 3: "d", 2: "c"}
	for key, value := range MapSortedByKey(myMap2) {
		fmt.Printf("%d->%s\n", key, value)
	}
}
