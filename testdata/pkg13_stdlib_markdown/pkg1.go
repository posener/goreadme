// Package pkg1 is a testing package.
//
// Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt
// ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
// laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
// voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat
// cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
//
// # Section Header
//
// Links in stdlib comment parser are markdown style bottom reference links.
// For example [this is a link] which the url is defined in the bottom of the
// comment section. Also links can be to local functions: [Func], or [pkg13.Func], or in
// other packages [goreadme.New].
//
// # Another Section Header
//
// You can use code blocks:
//
//	func main() {
//		println("hello world")
//	}
//
// You could also use numbered lists:
//  1. List item number 1.
//  2. List item number 2.
//  3. List item number 3.
//
// Or itemized list:
//   - Item 1.
//   - Item 2.
//
// [this is a link]: https://github.com/posener/goreadme
package pkg13

import "fmt"

// SomeConst is a package-level constant.
const SomeConst int = 5

// SomeVar is a package-level variable.
var SomeVar int = 6

func Func() {
	fmt.Println("hello")
}

type ExampleType struct {
	val              int
	ExampleInterface interface{}
}
