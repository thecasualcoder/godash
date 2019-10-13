package godash_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/godash"
	"testing"
)

func TestAny(t *testing.T) {

	t.Run("should return err if predicate is not a function", func(t *testing.T) {
		in := []int{1, 2, 3}

		_, err := godash.Any(in, "not a func")

		assert.EqualError(t, err, "predicateFn has to be a function")
	})

	t.Run("should return err if predicate function do not take exactly one argument", func(t *testing.T) {
		in := []int{1, 2, 3}

		{
			_, err := godash.Any(in, func() {})

			assert.EqualError(t, err, "predicate function has to take only one argument")
		}
		{
			_, err := godash.Any(in, func(int, int) {})

			assert.EqualError(t, err, "predicate function has to take only one argument")
		}
	})

	t.Run("should return err if predicate function do not return exactly one value", func(t *testing.T) {
		in := []int{1, 2, 3}

		{
			_, err := godash.Any(in, func(int) {})

			assert.EqualError(t, err, "predicate function should return only one return value")
		}
		{
			_, err := godash.Any(in, func(int) (bool, bool) { return true, true })

			assert.EqualError(t, err, "predicate function should return only one return value")

		}
	})

	t.Run("should return err if predicate function's return value is not a boolean", func(t *testing.T) {
		in := []int{1, 2, 3}

		_, err := godash.Any(in, func(int) int { return 0 })

		assert.EqualError(t, err, "predicate function should return a boolean value")
	})

	t.Run("should return err if input is not a slice", func(t *testing.T) {
		in := 1

		_, err := godash.Any(in, func(int) bool { return true })

		assert.EqualError(t, err, "not implemented for (int)")
	})

	t.Run("should return err if there is a type mismatch between predicate function's argument and input slice", func(t *testing.T) {
		in := []string{"hello", "world"}

		_, err := godash.Any(in, func(int) bool { return true })

		assert.EqualError(t, err, "predicate function's argument (int) has to be (string)")
	})

	t.Run("should return true if predicate passes for any of the element in input slice", func(t *testing.T) {
		in := []int{1, 1, 2, 3, 5, 8, 13}

		output, err := godash.Any(in, func(elem int) bool { return elem%5 == 0 })

		assert.NoError(t, err)
		assert.True(t, output)
	})

	t.Run("should return false if predicate fails for all the elements in input slice", func(t *testing.T) {
		in := []int{1, 1, 2, 3, 5, 8, 13}

		output, err := godash.Any(in, func(num int) bool { return num%6 == 0 })

		assert.NoError(t, err)
		assert.False(t, output)
	})

	t.Run("should support structs", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}
		in := []person{
			{name: "John", age: 12},
			{name: "Doe", age: 25},
		}

		{
			output, err := godash.Any(in, func(person person) bool { return person.age < 18 })

			assert.NoError(t, err)
			assert.True(t, output)
		}
		{
			output, err := godash.Any(in, func(person person) bool { return person.age < 30 })

			assert.NoError(t, err)
			assert.True(t, output)
		}
	})
}

func ExampleAny() {
	input := []int{0, 1, 2, 3, 4}

	output, _ := godash.Any(input, func(num int) bool {
		return num%3 == 0
	})

	fmt.Println(output)

	// Output: true
}
