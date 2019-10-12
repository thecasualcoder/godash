package godash

import (
	"fmt"
	"reflect"
)

// Map applies mapperFn on each item of in and puts it in out
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
	if mapperFnType.NumIn() != 1 {
		return fmt.Errorf("mapper function has to take only one argument")
	}

	if mapperFnType.NumOut() != 1 {
		return fmt.Errorf("mapper function should return only one return value")
	}

	if input.Kind() == reflect.Slice {
		if output.Kind() != reflect.Slice {
			return fmt.Errorf("output should be a slice for input of type slice")
		}
		if input.Type().Elem().Kind() != mapper.Type().In(0).Kind() {
			return fmt.Errorf("mapper function's first argument has to be the type of element of input slice")
		}
		if output.Type().Elem().Elem().Kind() != mapper.Type().Out(0).Kind() {
			return fmt.Errorf("mapper function's return type has to be the type of element of output slice")
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
	return fmt.Errorf("not implemented")
}
