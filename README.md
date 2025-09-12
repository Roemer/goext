# goext

Various small go extensions without dependencies.

Feel free to include this library or just copy the needed files or parts of them to your projects.

CommandRunner:
- [Run](#run)
- [RunGetOutput](#rungetoutput)
- [RunGetCombinedOutput](#rungetcombinedoutput)

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
- [StringTrimAllPrefix](#stringtrimallprefix)
- [StringTrimAllSuffix](#stringtrimallsuffix)
- [StringTrimNewlineSuffix](#stringtrimnewlinesuffix)
- [StringSplitByNewLine](#stringsplitbynewline)

TablePrinter:
- [TablePrinter](#tableprinter)

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

## CommandRunner
Allows running executable with arguments and various options.

The following options are available:
- WorkingDirectory: Runs the command in the given working directory
- OutputToConsole: Outputs stdout to the console
- SkipPostProcessOutput: Does not post-process the output (remove newlines)
- UseCurrentEnv
- AdditionalEnv

Usage is as follows:
```go
err := goext.NewCmdRunner("cmd", "arg1", "arg2").With<...>.Run()
```

### Run
Runs the command with the given options.
```go
err := goext.NewCmdRunner("myapp", "arg1", "arg2").Run()
```

### RunGetOutput
Runs the command and returns the separate output from stdout and stderr.
```go
stdout, stderr, err := goext.NewCmdRunner("myapp").RunGetOutput()
```

### RunGetCombinedOutput
Runs the command and returns the output from stdout and stderr combined.
```go
output, err := goext.NewCmdRunner("myapp").RunGetCombinedOutput()
```

## Slices

### SliceAppendIf
Appends the given values to a slice if the condition is fulfilled.
```go
mySlice := []int{1, 2, 3}
mySlice = goext.SliceAppendIf(mySlice, myCondition, 4, 5)
```

### SliceAppendIfFunc
Appends the given value if the condition is fulfilled. The value is lazily evaluated.
```go
valueFunc := func() []int { return []int{4, 5} }
mySlice := []int{1, 2, 3}
mySlice = goext.SliceAppendIfFunc(mySlice, myCondition, valueFunc)
```

### SlicePrepend
Prepends the given elements to the given array.
```go
mySlice := []int{1, 2, 3}
mySlice = goext.SlicePrepend(testSlice, 4, 5)
// [4, 5, 1, 2, 3]
```

### SlicePrependIf
Prepends the given values to a slice if the condition is fulfilled.
```go
mySlice := []int{1, 2, 3}
mySlice = goext.SlicePrependIf(mySlice, myCondition, 4, 5)
```

### SlicePrependIfFunc
Prepends the given value if the condition is fulfilled. The value is lazily evaluated.
```go
valueFunc := func() []int { return []int{4, 5} }
mySlice := []int{1, 2, 3}
mySlice = goext.SlicePrependIfFunc(mySlice, myCondition, valueFunc)
```

## Strings

### StringContainsAny
Checks if the string contains at least one of the substrings.
```go
goext.StringContainsAny("Hello", "Hello", "World")
goext.StringContainsAny("World", []string{"Hello", "World"}...)
// Both return true
```

### StringTrimAllPrefix
Trims all occurrences of the given prefix from the start of the string.
```go
value := goext.StringTrimAllPrefix("xxxyyyxxx", "x")
```

### StringTrimAllSuffix
Trims all occurrences of the given suffix from the end of the string.
```go
value := goext.StringTrimAllSuffix("xxxyyyxxx", "x")
```

### StringTrimNewlineSuffix
Trims all occurrences of newline characters from the end of the string.
```go
value := goext.StringTrimNewlineSuffix("my-text\n\r\n")
```

### StringSplitByNewLine
Splits the string by new line characters, supporting both "\n" and "\r\n".
```go
lines := goext.StringSplitByNewLine("line1\r\nline2\nline3")
```

## TablePrinter

### Examples
```go
tablePrinter := goext.NewTablePrinter(nil)
tablePrinter.SetHeaders("Id", "Name", "Age", "City")
tablePrinter.Columns[0].ValueAlignment = goext.TABLE_PRINTER_ALIGNMENT_RIGHT
tablePrinter.AddRows(
    []any{1, "Alice", "30", "New York"},
    []any{2, "Bob", "25", "Los Angeles"},
    []any{3, "Charlie", "35", "Chicago"},
)
tablePrinter.PrintStdout()
/* Prints:
┌──────┬───────────┬───────┬───────────────┐
│  Id  │  Name     │  Age  │  City         │
├──────┼───────────┼───────┼───────────────┤
│   1  │  Alice    │  30   │  New York     │
│   2  │  Bob      │  25   │  Los Angeles  │
│   3  │  Charlie  │  35   │  Chicago      │
└──────┴───────────┴───────┴───────────────┘
*/
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
