package godash

import (
	"fmt"
	"reflect"
)

// Map applies mapperFn on each item of in and puts it in out.
// Currently, input and output for type slice is supported.
//
// Validations:
//
//	1. Mapper function should take in one argument and return one argument
//	2. Mapper function's argument should be of the same type of each element of input slice.
//	3. Mapper function's output should be of the same type of each element of output slice.
//
// Validation failures are returned as error by the godash.Map to the caller.
func Map(in, out, mapperFn interface{}) error {
	input := reflect.ValueOf(in)
	output := reflect.ValueOf(out)
	if err := validateOut(output); err != nil {
		return err
	}

	mapper := reflect.ValueOf(mapperFn)
	if mapper.Kind() != reflect.Func {
		return fmt.Errorf("mapperFn has to be a function")
	}

	mapperFnType := mapper.Type()

	if mapperFnType.NumOut() != 1 {
		return fmt.Errorf("mapper function should return only one return value")
	}

	if input.Kind() == reflect.Slice {
		if output.Elem().Kind() != reflect.Slice {
			return fmt.Errorf("output should be a slice for input of type slice")
		}

		if mapperFnType.NumIn() != 1 {
			return fmt.Errorf("mapper function has to take only one argument")
		}

		if input.Type().Elem() != mapper.Type().In(0) {
			return fmt.Errorf("mapper function's first argument (%s) has to be (%s)", mapper.Type().In(0), input.Type().Elem())
		}
		if output.Elem().Type().Elem() != mapper.Type().Out(0) {
			return fmt.Errorf("mapper function's return type has to be (%s) but is (%s)", mapper.Type().Out(0), output.Elem().Type().Elem())
		}

		result := reflect.MakeSlice(output.Elem().Type(), 0, input.Len())
		for i := 0; i < input.Len(); i++ {
			arg := input.Index(i)

			returnValues := mapper.Call([]reflect.Value{arg})

			result = reflect.Append(result, returnValues[0])
		}
		output.Elem().Set(result)

		return nil
	}

	if input.Kind() == reflect.Map {
		if output.Elem().Kind() != reflect.Slice {
			return fmt.Errorf("output should be a slice for input of type slice")
		}

		if mapperFnType.NumIn() != 2 {
			return fmt.Errorf("mapper function has to take exactly two arguments")
		}

		if mapper.Type().In(0) != input.Type().Key() {
			return fmt.Errorf("mapper function's first argument (%s) has to be (%s)", mapper.Type().In(0), input.Type().Key())
		}
		if mapper.Type().In(1) != input.Type().Elem() {
			return fmt.Errorf("mapper function's second argument (%s) has to be (%s)", mapper.Type().In(1), input.Type().Elem())
		}
		if mapper.Type().Out(0) != output.Elem().Type().Elem() {
			return fmt.Errorf("mapper function's return type has to be (%s) but is (%s)", mapper.Type().Out(0), output.Elem().Type().Elem())
		}

		result := reflect.MakeSlice(output.Elem().Type(), 0, input.Len())
		for _, key :=  range input.MapKeys() {
			value := input.MapIndex(key)

			returnValues := mapper.Call([]reflect.Value{key, value})

			result = reflect.Append(result, returnValues[0])
		}
		output.Elem().Set(result)

		return nil
	}
	return fmt.Errorf("not implemented")
}
