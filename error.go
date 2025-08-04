package housekeeper

import "reflect"

func findError(vals []reflect.Value) error {
	for _, val := range vals {
		if isErrorType(val.Type()) {
			if val.IsNil() {
				return nil
			} else {
				return val.Interface().(error)
			}
		}
	}
	return nil
}
