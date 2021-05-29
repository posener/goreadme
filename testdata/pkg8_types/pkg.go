// Package pkg1 is a testing package.
package pkg1

// ExampleType is a type
type ExampleType struct {
	val              int
	ExampleInterface interface{}
}

// ExampleTypeFactory is a factory function for ExampleType.
func ExampleTypeFactory() ExampleType {
	return ExampleType{1, "test"}
}

// ExampleMethod is a method on ExampleType.
func (e ExampleType) ExampleMethod() {
}

// ExampleType2 is a type with an array
type ExampleType2 struct {
	val              []int
	ExampleInterface interface{}
}

// ExampleTypeInt is a one-liner type
type ExampleTypeInt struct{ val int }
