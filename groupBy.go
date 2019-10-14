package godash

import (
	"fmt"
	"reflect"
)

// GroupBy creates an object composed of keys generated from the results of running
// each element of slice throught iteration. The order of grouped values
// is determined by the order they occur in `collection`. The corresponding
// value of each key is an array of elements responsible for generating the
// key.
//
// Validations:
//
// 1. Group function should take one argument and return one value
// 2. Group function should return a one value
// 3. Group function's argument should be of the same type as the elements of the input slice
// 4. Output should be a map
// 5. The key type's output should be of the same type as the predicate function output
// 6. The value type's output (as sliceValueOutput) should be a slice
// 7. The element type of the sliceValueOutput should be of the same type as the element type of the input slice
// Validation errors are returned to the caller
func GroupBy(in, out, groupFn interface{}) error {

	input := reflect.ValueOf(in)
	output := reflect.ValueOf(out)
	group := reflect.ValueOf(groupFn)

	if group.Kind() != reflect.Func {
		return fmt.Errorf("groupFn (%s) has to be a function", group.Kind())
	}

	groupFnType := group.Type()

	if groupFnType.NumIn() != 1 {
		return fmt.Errorf("group function has to take only one argument")
	}

	if groupFnType.NumOut() != 1 {
		return fmt.Errorf("group function should return only one return value")
	}

	outputType := output.Elem().Type()
	outputKind := output.Elem().Kind()

	if outputKind != reflect.Map {
		return fmt.Errorf("output has to be a map")
	}

	if groupFnType.Out(0).Kind() != outputType.Key().Kind() {
		return fmt.Errorf("group function should return the type of key's output")
	}

	outputSliceType := outputType.Elem()
	outputSliceKind := outputSliceType.Kind()

	if outputSliceKind != reflect.Slice {
		return fmt.Errorf("The type of value's output should be a slice")
	}

	inputKind := input.Kind()

	if inputKind == reflect.Slice {
		inputSliceElemType := input.Type().Elem()
		groupFnArgType := groupFnType.In(0)

		if inputSliceElemType != outputSliceType.Elem() {
			return fmt.Errorf("The type of element of value's slice has to be a same the type of input")
		}

		if inputSliceElemType != groupFnArgType {
			return fmt.Errorf("group function's argument (%s) has to be (%s)", groupFnArgType, inputSliceElemType)
		}

		result := reflect.MakeMap(outputType)

		for i := 0; i < input.Len(); i++ {
			arg := input.Index(i)
			argValues := []reflect.Value{arg}
			returnValue := group.Call(argValues)[0]
			slice := result.MapIndex(returnValue)

			if slice.IsValid() {
				slice = reflect.Append(slice, argValues[0])
				result.SetMapIndex(returnValue, slice)

			} else {
				slice := reflect.MakeSlice(outputSliceType, 0, input.Len())
				slice = reflect.Append(slice, argValues[0])
				result.SetMapIndex(returnValue, slice)
			}
		}
		output.Elem().Set(result)

		return nil
	}
	return fmt.Errorf("not implemented for (%s)", inputKind)
}
