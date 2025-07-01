package goext

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

// Returns an iterator for the given map that yields the key-value pairs in sorted order.
func MapSortedByKey[Map ~map[K]V, K cmp.Ordered, V any](m Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range slices.Sorted(maps.Keys(m)) {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
