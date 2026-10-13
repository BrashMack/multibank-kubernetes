package reflector

import "reflect"

// DepthNew - метод, создающий указатель на структуру и инициализирующий их поля значением по умолчанию.
func DepthNew[T any]() *T {
	elem := new(T)

	value := reflect.ValueOf(elem)

	handle(value)

	return elem
}

func handle(value reflect.Value) {
	switch value.Kind() {
	case reflect.Pointer:
		fillPointer(value)
	case reflect.Struct:
		fillStruct(value)
	default:
		return
	}
}

func fillPointer(value reflect.Value) {
	if value.IsNil() {
		newValue := reflect.New(value.Type().Elem())

		if value.CanSet() {
			value.Set(newValue)
		}
	}

	handle(value.Elem())
}

func fillStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		handle(value.Field(i))
	}
}
