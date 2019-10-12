package godash

import (
	"fmt"
	"reflect"
)


func validateOut(output reflect.Value) error {
	zeroValue := reflect.Value{}
	if output == zeroValue {
		return fmt.Errorf("output is nil. Pass a reference to set output")
	}

	if output.IsNil() {
		return fmt.Errorf("output is nil. Pass a reference to set output")
	}

	if !output.Elem().CanSet() {
		return fmt.Errorf("cannot set out. Pass a reference to set output")
	}

	return nil
}
