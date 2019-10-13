package godash

import (
	"fmt"
	"reflect"
)

// Any checks if predicate returns truthy for any element of collection. Iteration is stopped once predicate returns truthy.
// Currently, input of type slice is supported
//
// Validations:
//
// 1. Predicate function should take one argument and return one value
// 2. Predicate function should return a bool value
// 3. Predicate function's argument should be of the same type as the elements of the input slice
//
// Validation errors are returned to the caller
func Any(in, predicateFn interface{}) (bool, error) {
	var output bool
	input := reflect.ValueOf(in)
	predicate := reflect.ValueOf(predicateFn)

	if predicate.Kind() != reflect.Func {
		return output, fmt.Errorf("predicateFn has to be a function")
	}

	predicateFnType := predicate.Type()
	if predicateFnType.NumIn() != 1 {
		return output, fmt.Errorf("predicate function has to take only one argument")
	}

	if predicateFnType.NumOut() != 1 {
		return output, fmt.Errorf("predicate function should return only one return value")
	}

	if predicateFnType.Out(0).Kind() != reflect.Bool {
		return output, fmt.Errorf("predicate function should return a boolean value")
	}

	inputKind := input.Kind()
	if inputKind == reflect.Slice {
		inputSliceElemType := input.Type().Elem
		predicateFnArgType := predicateFnType.In(0)
		if inputSliceElemType() != predicateFnArgType {
			return output, fmt.Errorf("predicate function's argument (%s) has to be (%s)", predicateFnArgType, inputSliceElemType())
		}

		for i := 0; i < input.Len(); i++ {
			arg := input.Index(i)
			returnValue := predicate.Call([]reflect.Value{arg})[0]
			if returnValue.Bool() {
				return true, nil
			}
		}

		return output, nil
	}

	return output, fmt.Errorf("not implemented for (%s)", inputKind)
}
