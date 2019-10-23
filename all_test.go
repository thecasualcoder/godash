package godash_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/godash"
)

func TestAllAndEvery(t *testing.T) {
	var funcs = map[string]func(interface{}, interface{}) (bool, error){
		"All()":   godash.All,
		"Every()": godash.Every,
	}

	for fnName, fn := range funcs {
		t.Run(fmt.Sprintf("%s should return err if predicate is not a function", fnName), func(t *testing.T) {
			in := []int{1, 2, 3}

			_, err := fn(in, "not a func")

			assert.EqualError(t, err, "predicateFn has to be a function")
		})

		t.Run(fmt.Sprintf("%s should return err if predicate function do not take exactly one argument", fnName), func(t *testing.T) {
			in := []int{1, 2, 3}

			{
				_, err := fn(in, func() {})

				assert.EqualError(t, err, "predicate function has to take only one argument")
			}
			{
				_, err := fn(in, func(int, int) {})

				assert.EqualError(t, err, "predicate function has to take only one argument")
			}
		})

		t.Run(fmt.Sprintf("%s should return err if predicate function do not return exactly one value", fnName), func(t *testing.T) {
			in := []int{1, 2, 3}

			{
				_, err := fn(in, func(int) {})

				assert.EqualError(t, err, "predicate function should return only one return value")
			}
			{
				_, err := fn(in, func(int) (bool, bool) { return true, true })

				assert.EqualError(t, err, "predicate function should return only one return value")

			}
		})

		t.Run(fmt.Sprintf("%s should return err if predicate function's return value is not a boolean", fnName), func(t *testing.T) {
			in := []int{1, 2, 3}

			_, err := fn(in, func(int) int { return 0 })

			assert.EqualError(t, err, "predicate function should return a boolean value")
		})

		t.Run(fmt.Sprintf("%s should return err if input is not a slice", fnName), func(t *testing.T) {
			in := 1

			_, err := fn(in, func(int) bool { return true })

			assert.EqualError(t, err, "not implemented for (int)")
		})

		t.Run(fmt.Sprintf("%s should return err if there is a type mismatch between predicate function's argument and input slice", fnName), func(t *testing.T) {
			in := []string{"hello", "world"}

			_, err := fn(in, func(int) bool { return true })

			assert.EqualError(t, err, "predicate function's argument (int) has to be (string)")
		})

		t.Run(fmt.Sprintf("%s should return true if predicate passes for all element in input slice", fnName), func(t *testing.T) {
			in := []int{1, 3, 5, 7, 9, 11, 13}

			output, err := fn(in, func(elem int) bool { return elem%2 == 1 })

			assert.NoError(t, err)
			assert.True(t, output)
		})

		t.Run(fmt.Sprintf("%s should return false if predicate fails for at least one of the elements in input slice", fnName), func(t *testing.T) {
			in := []int{1, 2, 5, 7, 9, 11, 13}

			output, err := fn(in, func(num int) bool { return num%2 == 1 })

			assert.NoError(t, err)
			assert.False(t, output)
		})

		t.Run(fmt.Sprintf("%s should support structs", fnName), func(t *testing.T) {
			type person struct {
				name string
				age  int
			}
			in := []person{
				{name: "John", age: 12},
				{name: "Doe", age: 25},
			}

			{
				output, err := fn(in, func(person person) bool { return person.age < 18 })

				assert.NoError(t, err)
				assert.False(t, output)
			}
			{
				output, err := fn(in, func(person person) bool { return person.age < 30 })

				assert.NoError(t, err)
				assert.True(t, output)
			}
		})
	}
}

func ExampleAll() {
	input := []int{0, 1, 2, 3, 4}

	output, _ := godash.All(input, func(num int) bool {
		return num >= 0
	})

	fmt.Println(output)

	// Output: true
}

func ExampleEvery() {
	input := []int{0, 1, 2, 3, 4}

	output, _ := godash.Every(input, func(num int) bool {
		return num >= 0
	})

	fmt.Println(output)

	// Output: true
}
