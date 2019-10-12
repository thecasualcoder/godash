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
	if err := validateMapperFunction(mapper); err != nil {
		return err
	}

	if input.Kind() == reflect.Slice {
		if input.Type().Elem().Kind() != mapper.Type().In(0).Kind() {
			return fmt.Errorf("mapper function's first argument has to be the type of element of input slice")
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


func validateMapperFunction(mapper reflect.Value) error {
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

	return nil
}
