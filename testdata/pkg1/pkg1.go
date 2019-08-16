// Package pkg1 is a testing package.
//
// Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt
// ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
// laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
// voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat
// cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
//
// Section Header
//
// A local link should just start with period and slash: ./internal, another local is ./internal/file.go.
// A web page link should just be written as is: https://goreadme.herokuapp.com, and with path: https://goreadme.herokuapp.com/projects.
// A url can also have a (title) http://example.org.
// A local path can also have a (title) ./pkg.go.
// A local path in inline code `go test ./`.
// Go path ellipsis (also inline ./...) should not be converted to link ./...
//
// Another Section Header
//
// Inline code can be defined with backticks: `prinlnt("hello world")`, or with indentation:
//
// 	func main() {
//		println("hello world")
// 	}
//
// Diff code block:
//
//	 func main() {
//	-	println("hello world")
//	+	fmt.Println("hello, world")
//	 }
//
// You could also use lists:
//
// 1. List item number 1.
// 1. List item number 2.
// 1. List item number 3.
//
// An image:
//
// (image/gopher) https://golang.org/doc/gopher/frontpage.png
package pkg1

import "fmt"

func Func() {
	fmt.Println("hello")
}
