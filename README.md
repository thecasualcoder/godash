# Godash

Inspired from [Lodash](https://github.com/lodash/lodash) for golang

# Functions

## Map

Map applies a mapper function on each element of an input and sets it in output. 

_Primitive types_

```go
func main() {
	input := []int{1, 2, 3, 4, 5}
	output := make([]int, 0)

	_ = godash.Map(input, &output, func(el int) int {
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
	output := make([]string, 0)

	_ = godash.Map(input, &output, func(person Person) string {
		return person.Name
	})

	fmt.Println(output) // prints John Doe
}
```
