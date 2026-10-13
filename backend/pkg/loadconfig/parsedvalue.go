package loadconfig

import (
	"reflect"

	"github.com/masharpik/vtb-hack-back/pkg/cast"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

func findParser(tp reflect.Type) func(string) (reflect.Value, error) {
	method, found := tp.MethodByName("Parse")

	if !found ||
		method.Type.NumIn() != 2 ||
		method.Type.In(1).Kind() != reflect.String ||
		method.Type.NumOut() != 2 ||
		method.Type.Out(0) != tp ||
		method.Type.Out(1) != errType {
		return nil
	}

	return func(s string) (reflect.Value, error) {
		results := method.Func.Call([]reflect.Value{reflect.New(tp).Elem(), reflect.ValueOf(s)})
		ierr := results[1].Interface()
		if ierr != nil {
			return reflect.Value{}, ierr.(error)
		}

		return results[0], nil
	}
}

func SetterForType(tp reflect.Type) (setterFnc, bool) {
	parser := findParser(tp)

	if parser == nil {
		ptrParser := findParser(reflect.PointerTo(tp))
		if ptrParser == nil {
			return DefaultSetter, false
		}
		//have pointer parser
		parser = func(s string) (reflect.Value, error) {
			pv, err := ptrParser(s)
			if err != nil {
				return reflect.Value{}, err
			}

			if pv.IsNil() {
				return reflect.Value{}, cast.ErrNilValue
			}

			return pv.Elem(), nil
		}
	}

	setter := func(val any, dst reflect.Value) error {
		str, err := cast.As[string](val)
		if err != nil {
			return err
		}

		rv, err := parser(str)
		if err != nil {
			return err
		}

		dst.Set(rv)

		return nil
	}

	return setter, true
}

func tryRegisterByParser(fi *flagInfo) error {
	setter, byParser := SetterForType(fi.value.Type())
	if !byParser {
		return nil
	}

	return registerFlag(ReflectValueVarPWithSetter(setter, true), fi, setter)
}
