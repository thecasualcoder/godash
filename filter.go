package godash

import (
	"fmt"
	"reflect"
)

// Filter out elements that fail the predicate.
//
// Input of type slice is supported as of now.
// Output is a slice in which filtered-in elements are stored.
// PredicateFn function is applied on each element of input to determine to filter or not
//
// Validations:
//
//	1. Input and Output's slice should be of same type
//	2. Predicate function can take one argument and return one argument
//	3. Predicate's return argument is always boolean.
//	4. Predicate's input argument should be input/output slice's element type.
//
// Validation errors are returned to the caller.
func Filter(in, out, predicateFn interface{}) error {
	input := reflect.ValueOf(in)

	output := reflect.ValueOf(out)
	if err := validateOut(output); err != nil {
		return err
	}
	if input.Type() != output.Elem().Type() {
		return fmt.Errorf("input(%s) and output(%s) should be of the same Type", input.Type(), output.Elem().Type())
	}

	predicate := reflect.ValueOf(predicateFn)
	if predicate.Type().NumOut() != 1 {
		return fmt.Errorf("predicate function should return only one return value - a boolean")
	}
	if predicateType := predicate.Type().Out(0).Kind(); predicateType != reflect.Bool {
		return fmt.Errorf("predicate function should return only a (boolean) and not a (%s)", predicateType)
	}

	if input.Kind() == reflect.Slice {
		{
			if input.Type().Elem().Kind() != predicate.Type().In(0).Kind() {
				return fmt.Errorf(
					"predicate function's first argument has to be the type (%s) instead of (%s)",
					input.Type().Elem(),
					predicate.Type().In(0),
				)
			}
		}

		result := reflect.MakeSlice(output.Elem().Type(), 0, input.Len())
		for i := 0; i < input.Len(); i++ {
			arg := input.Index(i)

			returnValues := predicate.Call([]reflect.Value{arg})
			predicatePassed := returnValues[0].Bool()

			if predicatePassed {
				result = reflect.Append(result, arg)
			}
		}
		output.Elem().Set(result)

		return nil
	}
	return fmt.Errorf("not implemented")
}
