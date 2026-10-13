package reflector

import (
	"fmt"
	"reflect"
	"strings"
)

// CheckFieldsForNilInterfaces - метод, проверяющий отсутствие nil-интерфейса в структуре.
func CheckFieldsForNilInterfaces(structure any) (err error) {
	value := reflect.ValueOf(structure)

	if value.Type().Kind() != reflect.Struct {
		return fmt.Errorf("%w: only struct", ErrNotSupportType)
	}

	var nilFieldNames []string

	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).Kind() != reflect.Interface {
			continue
		}

		if value.Field(i).IsNil() {
			nilFieldNames = append(nilFieldNames, value.Type().Field(i).Name)
		}
	}

	if len(nilFieldNames) != 0 {
		return fmt.Errorf("%w: %s", ErrFoundedNilInterfaceInStructure, strings.Join(nilFieldNames, ", "))
	}

	return nil
}
