# pkg13

Package pkg1 is a testing package.

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt
ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat
cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

# Section Header

Links in stdlib comment parser are markdown style bottom reference links.
For example [this is a link] which the url is defined in the bottom of the
comment section. Also links can be to local functions: [Func], or [pkg13.Func], or in
other packages [goreadme.New].

# Another Section Header

You can use code blocks:

```go
func main() {
	println("hello world")
}
```

You could also use numbered lists:

```go
1. List item number 1.
2. List item number 2.
3. List item number 3.
```

Or itemized list:

```diff
- Item 1.
- Item 2.
```

[this is a link]: [https://github.com/posener/goreadme](https://github.com/posener/goreadme)

## Examples

### Hello

Example_hello prints hello

```go
fmt.Println("hello")
```

 Output:

```
hello
```

### NoDoc

```go
fmt.Println("hello")
```

 Output:

```
hello
```

### Func

ExampleFunc tests func

```go
Func()
```

 Output:

```
hello
```

### WithName

ExampleFunc_withName tests func with a name

```go
Func()
```

 Output:

```
hello
```

### Assignment

ExampleExampleType tests using the type ExampleType

```go

example := new(ExampleType)
example.val = 1

```
