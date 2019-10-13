package godash_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/godash"
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("should filter elements that fail predicate", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8}
		var output []int

		err := godash.Filter(input, &output, func(a int) bool {
			return a%2 == 0
		})
		expected := []int{2, 4, 6, 8}

		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})

	t.Run("should struct types", func(t *testing.T) {
		type person struct {
			age int
		}
		input := []person{
			{30},
			{20},
			{40},
			{10},
		}
		var output []person

		err := godash.Filter(input, &output, func(p person) bool {
			return p.age > 20
		})
		expected := []person{
			{30},
			{40},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})

	t.Run("should validate predicate's arg", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8}
		var output []int

		err := godash.Filter(input, &output, func(a string) bool {
			return a == ""
		})

		assert.EqualError(t, err, "predicate function's first argument has to be the type (int) instead of (string)")
	})

	t.Run("should validate predicate's return type", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8}
		var output []int

		{
			err := godash.Filter(input, &output, func(a int) int {
				return a
			})
			assert.EqualError(t, err, "predicate function should return only a (boolean) and not a (int)")
		}
		{
			err := godash.Filter(input, &output, func(int) (int, bool) {
				return 1, true
			})
			assert.EqualError(t, err, "predicate function should return only one return value - a boolean")
		}
	})

	t.Run("should validate output's type", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8}
		var output []string

		err := godash.Filter(input, &output, func(a int) bool {
			return a == 0
		})

		assert.EqualError(t, err, "input([]int) and output([]string) should be of the same Type")
	})

	t.Run("should not panic if output is nil", func(t *testing.T) {
		in := []int{1, 2, 3}
		{
			var out []int

			err := godash.Filter(in, out, func(int) bool { return true })

			assert.EqualError(t, err, "output is nil. Pass a reference to set output")
		}
		{
			err := godash.Filter(in, nil, func(int) bool { return true })

			assert.EqualError(t, err, "output is nil. Pass a reference to set output")
		}
	})
}

func ExampleFilter() {
	input := []string{
		"rhythm",
		"of",
		"life",
	}
	var output []string

	_ = godash.Filter(input, &output, func(in string) bool {
		return len(in) > 3
	})

	fmt.Println(output)

	// Output:
	// [rhythm life]
}
