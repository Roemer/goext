package goext

import "os"

type envNamespace int

// Contains methods regarding environment variables.
var Env envNamespace = 0

// Checks if the given environment variable exists or not.
func (envNamespace) Exists(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

// Returns the value if the environment variable exists or the default otherwise.
func (envNamespace) ValueOrDefault(key string, defaultValue string) (string, bool) {
	if value, ok := os.LookupEnv(key); ok {
		return value, true
	}
	return defaultValue, false
}
