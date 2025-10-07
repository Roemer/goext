package goext

import "slices"

// Appends the given values to a slice if the condition is fulfilled.
func SliceAppendIf[T any](slice []T, cond bool, values ...T) []T {
	if cond {
		return append(slice, values...)
	}
	return slice
}

// Appends the given value if the condition is fulfilled. The value is lazily evaluated.
func SliceAppendIfFunc[T any](slice []T, cond bool, f func() []T) []T {
	if cond {
		values := f()
		return append(slice, values...)
	}
	return slice
}

// Appends the given elements if it is missing in the slice.
func SliceAppendIfMissing[T comparable](slice []T, values ...T) []T {
	for _, value := range values {
		if !slices.Contains(slice, value) {
			slice = append(slice, value)
		}
	}
	return slice
}

// Appends the given elements if it is missing in the slice. The value is lazily evaluated.
func SliceAppendIfMissingFunc[T comparable](slice []T, f func() []T) []T {
	values := f()
	for _, value := range values {
		if !slices.Contains(slice, value) {
			slice = append(slice, value)
		}
	}
	return slice
}

// Prepends the given elements to the given array.
func SlicePrepend[T any](slice []T, elems ...T) []T {
	return append(elems, slice...)
}

// Prepends the given values to a slice if the condition is fulfilled.
func SlicePrependIf[T any](slice []T, cond bool, values ...T) []T {
	if cond {
		return SlicePrepend(slice, values...)
	}
	return slice
}

// Prepends the given value if the condition is fulfilled. The value is lazily evaluated.
func SlicePrependIfFunc[T any](slice []T, cond bool, f func() []T) []T {
	if cond {
		values := f()
		return SlicePrepend(slice, values...)
	}
	return slice
}
