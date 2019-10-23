package godash

import (
	"fmt"
	"reflect"
)

// Find out elements.
//
// Input of type slice is supported as of now.
// Output is a elements are matched.
// PredicateFn function is applied on each element of input to determine to find element until it finds the element
//
// Validations:
//
//	1. Input's element type and Output should be of same type
//	2. Predicate function can take one argument and return one argument
//	3. Predicate's return argument is always boolean.
//	4. Predicate's input argument should be output element type.
//
// Validation errors are returned to the caller
func Find(in, out, predicateFn interface{}) error {
	input := reflect.ValueOf(in)
	output := reflect.ValueOf(out)
	inputTypeElem := input.Type().Elem()
	if inputTypeElem != output.Elem().Type() {
		return fmt.Errorf("input slice (%s) and output (%s) should be of the same Type", inputTypeElem, output.Elem().Type())
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
			if inputTypeElem.Kind() != predicate.Type().In(0).Kind() {
				return fmt.Errorf(
					"predicate function's first argument has to be the type (%s) instead of (%s)",
					inputTypeElem,
					predicate.Type().In(0),
				)
			}
		}
		for i := 0; i < input.Len(); i++ {
			arg := input.Index(i)

			returnValues := predicate.Call([]reflect.Value{arg})
			predicatePassed := returnValues[0].Bool()

			if predicatePassed {
				output.Elem().Set(arg)
				return nil
			}
		}
		return fmt.Errorf("element not found")
	}
	return fmt.Errorf("not implemented")
}
