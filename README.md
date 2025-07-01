# goext

Various small go extensions without dependencies.

Feel free to includ this library or just copy the needed files or parts of them to your projects.

Maps:
- [MapSortedByKey](#mapsortedbykey)

Slices:
- [SliceAppendIf](#sliceappendif)
- [SliceAppendIfFunc](#sliceappendiffunc)
- [SlicePrepend](#sliceprepend)
- [SlicePrependIf](#sliceprependif)
- [SlicePrependIfFunc](#sliceprependiffunc)

Strings:
- [StringContainsAny](#stringcontainsany)

Ternary:
- [Ternary](#ternary)
- [TernaryFunc](#ternaryfunc)
- [TernaryFuncErr](#ternaryfuncerr)

## Maps

### MapSortedByKey
Returns an iterator for the given map that yields the key-value pairs in sorted order.
```go
myMap := map[int]string{1: "b", 0: "a", 3: "d", 2: "c"}
for key, value := range goext.MapSortedByKey(myMap) {
    fmt.Printf("%d->%s\n", key, value)
}
// Always prints
// 0->a
// 1->b
// 2->c
// 3->d
```

## Slices

### SliceAppendIf
Appends the given values to a slice if the condition is fulfilled.
```go
mySlice := []int{1, 2, 3}
mySlice = SliceAppendIf(mySlice, myCondition, 4, 5)
```

### SliceAppendIfFunc
Appends the given value if the condition is fulfilled. The value is lazily evaluated.
```go
valueFunc := func() []int { return []int{4, 5} }
mySlice := []int{1, 2, 3}
mySlice = SliceAppendIfFunc(mySlice, myCondition, valueFunc)
```

### SlicePrepend
Prepends the given elements to the given array.
```go
mySlice := []int{1, 2, 3}
mySlice = SlicePrepend(testSlice, 4, 5)
// [4, 5, 1, 2, 3]
```

### SlicePrependIf
Prepends the given values to a slice if the condition is fulfilled.
```go
mySlice := []int{1, 2, 3}
mySlice = SlicePrependIf(mySlice, myCondition, 4, 5)
```

### SlicePrependIfFunc
Prepends the given value if the condition is fulfilled. The value is lazily evaluated.
```go
valueFunc := func() []int { return []int{4, 5} }
mySlice := []int{1, 2, 3}
mySlice = SlicePrependIfFunc(mySlice, myCondition, valueFunc)
```

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
value := goext.Ternary(myCondition, "Value A", "Value B")
```

### TernaryFunc
Like Ternary but uses functions to lazily evaluate the values.
```go
aFunc := func() string { return "Value A" }
bFunc := func() string { return "Value B"}
value := goext.TernaryFunc(myCondition, aFunc, bFunc)
```

### TernaryFuncErr
Like TernaryFunc but returns an error as well.
```go
aFunc := func() (string,error) { return "Value A", nil }
bFunc := func() (string,error) { return "Value B", nil }
value, err := goext.TernaryFuncErr(myCondition, aFunc, bFunc)
```
