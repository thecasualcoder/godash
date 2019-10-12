package godash

import (
	"fmt"
	"reflect"
)

// Reduce reduces the given collection using given reduce function
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
			return fmt.Errorf("reducer function's first argument has to be the type of out")
		}
		if input.Type().Elem().Kind() != reducerFnType.In(1).Kind() {
			return fmt.Errorf("reducer function's second argument has to be the type of element of input slice")
		}
		if outputKind != reducerFnType.Out(0).Kind() {
			return fmt.Errorf("reducer function's return type has to be the type of out")
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
		return fmt.Errorf("reduceFn has to be a function")
	}
	if reducerFnType.NumIn() != 2 {
		return fmt.Errorf("reducer function has to take exactly two argument")
	}
	if reducerFnType.NumOut() != 1 {
		return fmt.Errorf("reducer function should return only one return value")
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
