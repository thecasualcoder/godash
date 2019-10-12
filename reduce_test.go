package godash_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/godash"
	"strconv"
	"testing"
)

func TestReduce(t *testing.T) {
	t.Run("support primitive types", func(t *testing.T) {
		{
			in := []int{1, 2, 3}
			var out int

			err := godash.Reduce(in, &out, func(acc, element int) int {
				return acc + element
			})

			expected := 6
			assert.NoError(t, err)
			assert.Equal(t, expected, out)
		}
		{
			in := []int{1, 2, 3}
			var out string

			err := godash.Reduce(in, &out, func(acc string, element int) string {
				return acc + strconv.Itoa(element)
			})

			expected := "123"
			assert.NoError(t, err)
			assert.Equal(t, expected, out)
		}
	})

	t.Run("support structs", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		in := []person{
			{name: "john", age: 20},
			{name: "doe", age: 23},
		}
		out := 0
		expected := 43

		err := godash.Reduce(in, &out, func(acc int, p person) int {
			return acc + p.age
		})

		assert.NoError(t, err)
		assert.Equal(t, expected, out)
	})

	add := func(acc, element int) int {
		return acc + element
	}

	t.Run("should not panic if output is nil", func(t *testing.T) {
		in := []int{1, 2, 3}
		{
			var out int

			err := godash.Reduce(in, out, add)

			assert.EqualError(t, err, "cannot set out. Pass a reference to set output")
		}

		{
			err := godash.Reduce(in, nil, add)

			assert.EqualError(t, err, "output is nil. Pass a reference to set output")
		}
	})

	t.Run("should not accept reducer function that are not functions", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out int

		err := godash.Reduce(in, &out, 7)

		assert.EqualError(t, err, "reduceFn has to be a (func) and not (int)")
	})

	t.Run("should not accept reducer function that do not take exactly two argument", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out int

		{
			err := godash.Reduce(in, &out, func() int { return 0 })
			assert.EqualError(t, err, "reduceFn has to take exactly 2 arguments and not 0 argument(s)")
		}

		{
			err := godash.Reduce(in, &out, func(int) int { return 0 })
			assert.EqualError(t, err, "reduceFn has to take exactly 2 arguments and not 1 argument(s)")
		}
	})

	t.Run("should not accept reducer function that do not return exactly one value", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out int

		{
			err := godash.Reduce(in, &out, func(int, int) {})
			assert.EqualError(t, err, "reduceFn should have only one return value and not 0 return type(s)")
		}

		{
			err := godash.Reduce(in, &out, func(int, int) (int, int) { return 0, 0 })
			assert.EqualError(t, err, "reduceFn should have only one return value and not 2 return type(s)")
		}
	})

	t.Run("should accept reducer function whose first argument's kind should be output's kind", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out int

		{
			err := godash.Reduce(in, &out, func(string, int) int { return 0 })
			assert.EqualError(t, err, "reduceFn's first argument's type(string) has to be the type of out(int)")
		}

		{
			err := godash.Reduce(in, &out, func(int, int) int { return 0 })
			assert.NoError(t, err)
		}
	})

	t.Run("should accept reducer function whose second argument's kind should be input slice's element kind", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out string

		{
			err := godash.Reduce(in, &out, func(string, string) string { return "" })
			assert.EqualError(t, err, "reduceFn's second argument's type(string) has to be the type of element of input slice(int)")
		}

		{
			err := godash.Reduce(in, &out, func(string, int) string { return "" })
			assert.NoError(t, err)
		}
	})

	t.Run("should accept reducer function whose return kind should be output's kind", func(t *testing.T) {
		in := []int{1, 2, 3}
		var out string

		{
			err := godash.Reduce(in, &out, func(string, int) int { return 0 })
			assert.EqualError(t, err, "reduceFn's return type(int) has to be the type of out(string)")
		}

		{
			err := godash.Reduce(in, &out, func(string, int) string { return "" })
			assert.NoError(t, err)
		}
	})
}
