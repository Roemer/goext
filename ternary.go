package goext

// A simple ternary function that returns one of two values based on a boolean condition.
func Ternary[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Like Ternary but uses functions to lazily evaluate the values.
func TernaryFunc[T any](cond bool, vtrue, vfalse func() T) T {
	if cond {
		return vtrue()
	}
	return vfalse()
}

// Like TernaryFunc but returns an error as well.
func TernaryFuncErr[T any](cond bool, vtrue, vfalse func() (T, error)) (T, error) {
	if cond {
		return vtrue()
	}
	return vfalse()
}
