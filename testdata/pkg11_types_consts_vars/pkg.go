// Package pkg11 is a testing package.
package pkg11

const (
	// ConstVal1 is a const in a const block.
	ConstVal1 int = 1
)

// ConstVal2 is a const outside a const block.
const ConstVal2 string = "2"

var (
	// VarVal1 is a var in a var block.
	VarVal1 int = 3
)

// VarVal2 is a var outside a var block.
var VarVal2 string = "4"

// ExampleType is a type
type ExampleType struct {
	val              int
	ExampleInterface interface{}
}

// ExampleType2 is a type with an array
type ExampleType2 struct {
	val              []int
	ExampleInterface interface{}
}

// ExampleTypeInt is a one-liner type
type ExampleTypeInt struct{ val int }

// ConstType1 is a constant of type ExampleTypeInt.
const ConstType1 ExampleTypeInt = ExampleTypeInt{5}

// VarType1 is a variable of type ExampleTypeInt.
var VarType1 ExampleTypeInt = ExampleTypeInt{6}

// ExampleFactoryFunction is a function that returns an ExampleType by value.
func ExampleFactoryFunction() ExampleType {
	return ExampleType{
		val:              1,
		ExampleInterface: "something",
	}
}

// ExampleFactoryFunction2 is a function that returns an ExampleType by pointer, and an error.
func ExampleFactoryFunction2() (*ExampleType, error) {
	return &ExampleType{
		val:              1,
		ExampleInterface: "something",
	}, nil
}

// ExampleMethod is a method on an ExampleType that takes the receiver by value.
func (et ExampleType) ExampleMethod() string {
	return "something"
}

// ExampleMethod2 is a method on an ExampleType that takes the receiver by pointer.
func (et *ExampleType) ExampleMethod2() string {
	return "*something"
}
