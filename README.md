# goext

Various small go extensions without dependencies.

Feel free to includ this library or just copy the needed files or parts of them to your projects.

Strings:
- [StringContainsAny](#stringcontainsany)

Ternary:
- [StringContainsAny](#stringcontainsany)

## Strings

### StringContainsAny
Checks if the string contains at least one of the substrings.
```go
goext.StringContainsAny("Hello", []string{"Hello", "World"})
goext.StringContainsAny("World", []string{"Hello", "World"})
// Both return true
```

## Ternary

### Ternary
A simple ternary function that returns one of two values based on a boolean condition.
```go
value := Ternary(myCondition, "Value A", "Value B")
```

### TernaryFunc
// Like Ternary but uses functions to lazily evaluate the values.
```go
aFunc := func() string { return "Value A" }
bFunc := func() string { return "Value B"}
value := TernaryFunc(myCondition, aFunc, bFunc)
```

### TernaryFuncErr
// Like TernaryFunc but returns an error as well.
```go
aFunc := func() (string,error) { return "Value A", nil }
bFunc := func() (string,error) { return "Value B", nil }
value, err := TernaryFuncErr(myCondition, aFunc, bFunc)
```
