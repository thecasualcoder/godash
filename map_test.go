package godash_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/godash"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("support primitive types", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := make([]int, 0)

		err := godash.Map(in, &out, func(element int) int {
			return element * element
		})

		expected := []int{1, 4, 9}
		assert.NoError(t, err)
		assert.Equal(t, expected, out)
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
		out := make([]string, 0)
		expected := []string{"john", "doe"}

		err := godash.Map(in, &out, func(p person) string {
			return p.name
		})

		assert.NoError(t, err)
		assert.Equal(t, expected, out)
	})
}
