# Godash

[![Build Status](https://travis-ci.org/thecasualcoder/godash.svg?branch=master)](https://travis-ci.org/thecasualcoder/godash)
[![Go Doc](https://godoc.org/github.com/thecasualcoder/godash?status.svg)](https://godoc.org/github.com/thecasualcoder/godash)

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
For more [docs](https://godoc.org/github.com/thecasualcoder/godash#Map).

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

Filter out elements that fail the predicate.
For more [docs](https://godoc.org/github.com/thecasualcoder/godash#Filter).

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

Reduce can accept a reducer and apply the reducer on each element of the input slice while providing an accumulator to save the reduce output
For more [docs](https://godoc.org/github.com/thecasualcoder/godash#Reduce).

```go
func main() {
    input := []string{"count", "words", "and", "print", "words", "count"}
	accumulator := map[string]int{}

	_ = godash.Reduce(input, &accumulator, func(acc map[string]int, element string) map[string]int {
		if _, present := acc[element]; present {
			acc[element] = acc[element] + 1
		} else {
			acc[element] = 1
		}
		return acc
	})

	bytes, _ := json.MarshalIndent(accumulator, "", "  ")
	fmt.Println(string(bytes))

	// Output:
	//{
	//   "and": 1,
	//   "count": 2,
	//   "print": 1,
	//   "words": 2
	//}

}
```

```go
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
