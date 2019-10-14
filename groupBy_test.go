package godash_test

import (
	"fmt"
	"testing"
	"github.com/thecasualcoder/godash"
)

func TestGroupBy(t *testing.T) {

}

func ExampleGroupBy() {
	type Person struct {
		name string
		age  int
	}
	john := Person{name: "John", age: 25}
	doe := Person{name: "Doe", age: 30}
	wick := Person{name: "Wick", age: 25}

	input := []Person{
		john,
		doe,
		wick,
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
