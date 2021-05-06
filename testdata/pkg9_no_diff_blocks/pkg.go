// Package pkg3_skip_examples is a testing package.
//
// Diff code block:
//
// 	 func main() {
// 	-	println("hello world")
// 	+	fmt.Println("hello, world")
// 	 }
package pkg9_no_diff_blocks

import "fmt"

func Func() {
	fmt.Println("hello")
}
