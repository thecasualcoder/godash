package godash_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/godash"
)

type Person struct {
	name string
	age  int
}

func TestGroupBy(t *testing.T) {

	t.Run(fmt.Sprintf("GroupBy should return err if groupFn is not a function"), func(t *testing.T) {
		in := []Person{}
		var output map[int][]Person

		err := godash.GroupBy(in, &output, "not a func")

		assert.EqualError(t, err, "groupFn should to be a function")
	})

	t.Run(fmt.Sprintf("GroupBy should return err if predicate function do not take exactly one argument"), func(t *testing.T) {
		in := []Person{}
		var out map[int][]Person

		{
			err := godash.GroupBy(in, &out, func() {})

			assert.EqualError(t, err, "group function has to take only one argument")
		}

		{
			err := godash.GroupBy(in, &out, func(int, int) {})

			assert.EqualError(t, err, "group function has to take only one argument")
		}

	})

	t.Run(fmt.Sprintf("GroupBy should return err if group function do not return exactly one value"), func(t *testing.T) {
		in := []Person{}
		var out map[int][]Person

		{
			err := godash.GroupBy(in, &out, func(int) {})

			assert.EqualError(t, err, "group function should return only one return value")
		}
		{
			err := godash.GroupBy(in, &out, func(int) (bool, bool) { return true, true })

			assert.EqualError(t, err, "group function should return only one return value")

		}
	})

	t.Run(fmt.Sprintf("GroupBy should return err if output is not the pointer to the map"), func(t *testing.T) {
		in := []Person{}
		var out []Person

		{
			err := godash.GroupBy(in, &out, func(int) bool { return true })

			assert.EqualError(t, err, "output has to be a map")
		}
	})

	t.Run(fmt.Sprintf("GroupBy should return err if the type's return of group function is not the same type of the key of output"), func(t *testing.T) {
		in := []Person{}
		var out map[int][]Person

		{
			err := godash.GroupBy(in, &out, func(int) bool { return true })

			assert.EqualError(t, err, "group function should return the type of key's output")
		}
	})

	t.Run(fmt.Sprintf("GroupBy should return err if the type of value's output is not a slice"), func(t *testing.T) {
		in := []Person{}
		var out map[int]Person

		{
			err := godash.GroupBy(in, &out, func(person Person) int {
				return person.age
			})

			assert.EqualError(t, err, "The type of value's output should be a slice")
		}
	})

	t.Run(fmt.Sprintf("GroupBy should return err if the type of element of value's slice is not a same the type of input"), func(t *testing.T) {
		in := []Person{}
		var out map[int][]int

		{
			err := godash.GroupBy(in, &out, func(person Person) int {
				return person.age
			})

			assert.EqualError(t, err, "The type of element of value's slice has to be a same the type of input")
		}
	})

	t.Run("GroupBy should return err if the type of element of input is not a same the type of groupFn input", func(t *testing.T) {
		in := []Person{}
		var out map[Person][]Person

		{
			err := godash.GroupBy(in, &out, func(i int) Person {
				return Person{}
			})

			assert.EqualError(t, err, "group function's argument (int) has to be (godash_test.Person)")
		}
	})
}

func ExampleGroupBy() {

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
