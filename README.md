# Godash

[![Build Status](https://travis-ci.org/thecasualcoder/godash.svg?branch=master)](https://travis-ci.org/thecasualcoder/godash)
[![Go Doc](https://godoc.org/github.com/thecasualcoder/godash?status.svg)](https://godoc.org/github.com/thecasualcoder/godash)

[![DeepSource](https://static.deepsource.io/deepsource-badge-light.svg)](https://deepsource.io/gh/thecasualcoder/godash/?ref=repository-badge)

Inspired from [Lodash](https://github.com/lodash/lodash) for golang

## Why?

- I did not like most map/reduce implementations that returned an `interface{}` which had to be typecasted. This library follows the concept of how `json.Marshal` works. Create an output variable **outside** the functions and pass a **pointer reference** to it, so it can be **set**.
- This library heavily makes use of `reflect` package and hence will have an **impact on performance**. **DO NOT USE THIS IN PRODUCTION**. This repository is more of a way to learn the reflect package and measure its performance impact.
- All functions have **validations** on how mapper function/predicate functions should be written. So even if we lose out on compile time validation, the library still **does not panic** if it does not know how to handle an argument passed to it.

## Available Functions

1. [Map](#Map)
2. [Filter](#Filter)
3. [Reduce](#Reduce)
4. [Any](#Any-or-Some) or [Some](#Any-or-Some)
5. [Find](#Find)
6. [All](#All-or-Every) or [Every](#All-or-Every)
7. [GroupBy](#GroupBy)

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

```go
func main() {
	input := map[string]int{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		}
	var output []int

	godash.Map(input, &output, func(el int) int {
		return el * el
	})

	fmt.Println(output) // prints 1 4 9
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

	godash.Filter(input, &output, func(person Person) bool {
		return person.Age > 25
	})

	fmt.Println(output) // prints {Doe 30}
}
```

### Reduce

Reduce can accept a reducer and apply the reducer on each element of the input slice while providing an accumulator to save the reduce output.
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

### Any or Some

Any or Some checks if predicate returns truthy for any element of collection. Iteration is stopped once predicate returns truthy.
For more [docs](https://godoc.org/github.com/thecasualcoder/godash#Any).


```go
func main() {
	input := []int{1, 2, 3, 4, 5}
	var output []int
	output, _ := godash.Any(input, func(num int) bool {
		return num % 7 == 0
	})
	fmt.Println(output) // prints false
}
```

```go
func main() {
	input := []Person{
		{Name: "John", Age: 25},
		{Name: "Doe", Age: 15},
	}
	var output int
	output, _ := godash.Some(input, func(person Person) bool {
		return person.Age < 18
	})
	fmt.Println(output) // prints true
}
```


### Find

Returns the first element which passes the predicate.
For more [docs](https://godoc.org/github.com/thecasualcoder/godash#Find).
```go
func main() {
	input := []string{"john","wick","will"}
	var output string

	godash.Find(input, &output, func(element string) bool {
    	return strings.HasPrefix(element, "w") // starts with
	}
	// output is "wick"
	fmt.Println(output)
}
```

### All or Every 

All or Every checks if predicate returns truthy for all element of collection. Iteration is stopped once predicate returns falsely. 
For more [docs](https://godoc.org/github.com/thecasualcoder/godash#All). 

```go 
func main() { 
	input := []int{1, 2, 3, 4, 5} 
	var output bool 
	output, _ := godash.All(input, func(num int) bool { 
		return num >= 1 
	}) 
	fmt.Println(output) // prints true 
} 
``` 

```go 
func main() { 
	input := []Person{ 
		{Name: "John", Age: 25}, 
		{Name: "Doe", Age: 15}, 
	} 
	var output bool 
	output, _ := godash.Every(input, func(person Person) bool {
		return person.Age < 18 
	}) 
	fmt.Println(output) // prints false 
}
```

### GroupBy

GroupBy creates an object composed of keys generated from the results of running each element of slice throught iteration. The order of grouped values is determined by the order they occur in slice. The corresponding value of each key is an array of elements responsible for generating the key.
For more [docs](https://godoc.org/github.com/thecasualcoder/godash#GroupBy).

```go
func main() {
	input := []Person{
		Person{name: "John", age: 25},
		Person{name: "Doe", age: 30},
		Person{name: "Wick", age: 25},
	}

	var output map[int][]Person

	godash.GroupBy(input, &output, func(person Person) int {
		return person.age
	})
	
	fmt.Printf("Groups Count: %d", len(output))
	for k, v := range output {
		fmt.Printf("\nGroup %d: ", k)
		for _, elem := range v {
			fmt.Printf(" %s,", elem.name)
		}
	}

	// Output:
	// Groups Count: 2
	// Group 25:  John, Wick,
	// Group 30:  Doe,
}
```
