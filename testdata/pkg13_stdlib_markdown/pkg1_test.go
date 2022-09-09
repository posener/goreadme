package pkg13

import "fmt"

// Example_hello prints hello
func Example_hello() {
	fmt.Println("hello")
	// Output: hello
}

func Example_noDoc() {
	fmt.Println("hello")
	// Output: hello
}

// ExampleFunc tests func
func ExampleFunc() {
	Func()
	// Output: hello
}

// ExampleFunc_withName tests func with a name
func ExampleFunc_withName() {
	Func()
	// Output: hello
}

// ExampleExampleType tests using the type ExampleType
func ExampleExampleType_assignment() {
	example := new(ExampleType)
	example.val = 1
}
