# Godash

Inspired from [Lodash](https://github.com/lodash/lodash) for golang

## Why?

I did not like most map/reduce implementations that returned an `interface{}` which had to be typecasted. This library follows the concept of how `json.Marshal` works. Create an output variable **outside** the functions and pass a **pointer reference** to it, so it can be **set**.
This library heavily makes use of `reflect` package and hence will have an **impact on performance**. Use it with care. All functions have **validations** on how mapper function/predicate functions should be written. So even if we lose out on compile time validation, the library still **does not panic** if it does not know how to handle an argument passed to it.

## Available Functions

1. [Map](#Map)
2. [Filter](#Filter)
3. [Reduce](#Reduce)

## Usages

### Map

Map applies a mapper function on each element of an input and sets it in output. 

_Primitive types_

```go
func main() {
	input := []int{1, 2, 3, 4, 5}
	var output []int

	godash.Map(input, &output, func(el int) int {
		return el * el
	})

	fmt.Println(output) // prints 1 4 9 16 25
}
```

_Struct type_

```go
type Person struct {
	Name string
	Age Int
}

func main() {
	input := []Person{
		{Name: "John", Age: 22},
		{Name: "Doe", Age: 23},
	}
	var output []string

	godash.Map(input, &output, func(person Person) string {
		return person.Name
	})

	fmt.Println(output) // prints John Doe
}
```

### Filter

Filter out elements that fail the predicate

```go
func main() {
	input := []int{1, 2, 3, 4, 5}
	var output []int

	godash.Filter(input, &output, func(element int) bool {
		return element % 2 == 0
	})

	fmt.Println(output) // prints 2 4
}
```

```go
func main() {
	input := []Person{
		{Name: "John", Age: 20},
		{Name: "Doe", Age: 30},
	}
	var output []string

	godash.Filter(input, &output, func(person Person) string {
		return person.Age > 25
	})

	fmt.Println(output) // prints {Doe 30}
}
```

### Reduce

Reduce reduces the given collection using given reduce function 

_Primitive types_

```go
func main() {
	input := []int{1, 2, 3, 4, 5}
	var output int

	godash.Reduce(input, &output, func(sum, element int) int {
		return sum + element
	})

	fmt.Println(output) // prints 15
}
```

_Struct type_

```go
type Person struct {
	Name string
	Age Int
}

func main() {
	input := []Person{
		{Name: "John", Age: 22},
		{Name: "Doe", Age: 23},
	}
	var output int

	godash.Reduce(input, &output, func(sum int, person Person) int {
		return sum + person.Age
	})

	fmt.Println(output) // prints 45
}
```
