package housekeeper

import "reflect"

var (
	errorType = reflect.TypeFor[error]()
)

func findError(vals []reflect.Value) error {
	for _, val := range vals {
		if val.Type().AssignableTo(errorType) {
			if val.IsNil() {
				return nil
			} else {
				return val.Interface().(error)
			}
		}
	}
	return nil
}
