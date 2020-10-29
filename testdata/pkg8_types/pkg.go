// Package pkg1 is a testing package.
package pkg1

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
