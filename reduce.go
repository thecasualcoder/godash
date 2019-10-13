package godash

import (
	"fmt"
	"reflect"
)

// Reduce can accept a reducer and apply the reducer on each element
// of the input slice while providing an accumulator to save the reduce output.
//
// Input of type slice is supported as of now.
// Output is the accumulator.
// ReduceFn is the reducer function.
//
// Whatever ReduceFn returns is fed as accumulator for the next iteration.
// Reduction happens from left-to-right.
//
// Reduce does the following validations:
//
//	1. Reducer function should accept exactly 2 arguments and return 1 argument
//  2. Reducer function's second argument should be the same type as input slice's element type
// 	3. Reducer function's return type should be the same as that of the accumulator
//
// Validation errors are returned to the caller.
func Reduce(in, out, reduceFn interface{}) error {
	input := reflect.ValueOf(in)
	output := reflect.ValueOf(out)
	if err := isReferenceType(output); err != nil {
		return err
	}

	reducer := reflect.ValueOf(reduceFn)
	if err := validateReducer(reducer); err != nil {
		return err
	}

	if input.Kind() == reflect.Slice {
		outputKind := output.Elem().Kind()
		reducerFnType := reducer.Type()
		if outputKind != reducerFnType.In(0).Kind() {
			return fmt.Errorf("reduceFn's first argument's type(%s) has to be the type of out(%s)", reducerFnType.In(0).Kind(), outputKind)
		}
		if input.Type().Elem().Kind() != reducerFnType.In(1).Kind() {
			return fmt.Errorf("reduceFn's second argument's type(%s) has to be the type of element of input slice(%s)", reducerFnType.In(1).Kind(), input.Type().Elem().Kind())
		}
		if outputKind != reducerFnType.Out(0).Kind() {
			return fmt.Errorf("reduceFn's return type(%s) has to be the type of out(%s)", reducerFnType.Out(0).Kind(), outputKind)
		}

		result := output.Elem()
		for i := 0; i < input.Len(); i++ {
			arg := input.Index(i)
			returnValues := reducer.Call([]reflect.Value{result, arg})

			result = returnValues[0]
		}
		output.Elem().Set(result)

		return nil
	}
	return fmt.Errorf("not implemented")
}

func validateReducer(reducer reflect.Value) error {
	reducerFnType := reducer.Type()
	if reducer.Kind() != reflect.Func {
		return fmt.Errorf("reduceFn has to be a (func) and not (%s)", reducer.Kind())
	}
	if reducerFnType.NumIn() != 2 {
		return fmt.Errorf("reduceFn has to take exactly 2 arguments and not %d argument(s)", reducerFnType.NumIn())
	}
	if reducerFnType.NumOut() != 1 {
		return fmt.Errorf("reduceFn should have only one return value and not %d return type(s)", reducerFnType.NumOut())
	}
	return nil
}

func isReferenceType(output reflect.Value) error {
	zeroValue := reflect.Value{}
	if output == zeroValue {
		return fmt.Errorf("output is nil. Pass a reference to set output")
	}
	if output.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot set out. Pass a reference to set output")
	}
	return nil
}
