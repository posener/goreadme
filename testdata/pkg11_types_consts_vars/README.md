# pkg11

Package pkg11 is a testing package.

## Constants

```go
const (
    // ConstVal1 is a const in a const block.
    ConstVal1 int = 1
)
```

ConstVal2 is a const outside a const block.

```go
const ConstVal2 string = "2"
```

## Variables

```go
var (
    // VarVal1 is a var in a var block.
    VarVal1 int = 3
)
```

VarVal2 is a var outside a var block.

```go
var VarVal2 string = "4"
```

## Types

### type [ExampleType](/pkg.go#L21)

```go
type ExampleType struct { ... }
```

ExampleType is a type

#### func [ExampleFactoryFunction](/pkg.go#L42)

```go
func ExampleFactoryFunction() ExampleType
```

ExampleFactoryFunction is a function that returns an ExampleType by value.

#### func [ExampleFactoryFunction2](/pkg.go#L50)

```go
func ExampleFactoryFunction2() (*ExampleType, error)
```

ExampleFactoryFunction2 is a function that returns an ExampleType by pointer, and an error.

#### func (ExampleType) [ExampleMethod](/pkg.go#L58)

```go
func (et ExampleType) ExampleMethod() string
```

ExampleMethod is a method on an ExampleType that takes the receiver by value.

#### func (*ExampleType) [ExampleMethod2](/pkg.go#L63)

```go
func (et *ExampleType) ExampleMethod2() string
```

ExampleMethod2 is a method on an ExampleType that takes the receiver by pointer.

### type [ExampleType2](/pkg.go#L27)

```go
type ExampleType2 struct { ... }
```

ExampleType2 is a type with an array

### type [ExampleTypeInt](/pkg.go#L33)

```go
type ExampleTypeInt struct { ... }
```

ExampleTypeInt is a one-liner type

#### Constants

ConstType1 is a constant of type ExampleTypeInt.

```go
const ConstType1 ExampleTypeInt = ExampleTypeInt{5}
```

#### Variables

VarType1 is a variable of type ExampleTypeInt.

```go
var VarType1 ExampleTypeInt = ExampleTypeInt{6}
```
