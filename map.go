package godash

import (
	"fmt"
	"reflect"
)

// Map applies mapperFn on each item of in and puts it in out
func Map(in, out, mapperFn interface{}) error {
	mapper := reflect.ValueOf(mapperFn)
	if mapper.Kind() != reflect.Func {
		return fmt.Errorf("mapperFn has to be a function")
	}
	input := reflect.ValueOf(in)
	output := reflect.ValueOf(out)
	if !output.Elem().CanSet() {
		return fmt.Errorf("cannot set out. Pass a reference to set output")
	}

	if input.Kind() == reflect.Slice {
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
