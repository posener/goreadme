# pkg1

Package pkg1 is a testing package.

## Constants

```golang
const (
    // ConstVal1 is a const in a const block.
    ConstVal1 int = 1
)
```

ConstVal2 is a const outside a const block.

```golang
const ConstVal2 string = "2"
```

## Variables

```golang
var (
    // VarVal1 is a var in a var block.
    VarVal1 int = 3
)
```

VarVal2 is a var outside a var block.

```golang
var VarVal2 string = "4"
```

## Types

### type [ExampleType](/pkg.go#L21)

`type ExampleType struct { ... }`

ExampleType is a type

### Assignment

ExampleExampleType tests using the type ExampleType

```golang

example := new(ExampleType)
example.val = 1

```

#### func [ExampleFactoryFunction](/pkg.go#L36)

`func ExampleFactoryFunction() ExampleType`

ExampleFactoryFunction is a function that returns an ExampleType by value.

#### func [ExampleFactoryFunction2](/pkg.go#L44)

`func ExampleFactoryFunction2() (*ExampleType, error)`

ExampleFactoryFunction2 is a function that returns an ExampleType by pointer, and an error.

#### func (ExampleType) [ExampleMethod](/pkg.go#L52)

`func (et ExampleType) ExampleMethod() string`

ExampleMethod is a method on an ExampleType that takes the receiver by value.

#### func (*ExampleType) [ExampleMethod2](/pkg.go#L57)

`func (et *ExampleType) ExampleMethod2() string`

ExampleMethod2 is a method on an ExampleType that takes the receiver by pointer.

### type [ExampleType2](/pkg.go#L27)

`type ExampleType2 struct { ... }`

ExampleType2 is a type with an array

### type [ExampleTypeInt](/pkg.go#L33)

`type ExampleTypeInt struct { ... }`

ExampleTypeInt is a one-liner type
