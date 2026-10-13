package loadconfig

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	"github.com/spf13/pflag"

	"github.com/masharpik/vtb-hack-back/pkg/cast"
)

var _ pflag.Value = &reflectValue{}

type reflectValue struct {
	Value        reflect.Value
	Inputs       []string
	Setter       setterFnc
	StringerPref bool
}

func (s *reflectValue) String() string {
	val := s.Value.Interface()

	if s.StringerPref {
		stringer, ok := val.(fmt.Stringer)

		if !ok {
			stringer, ok = s.Value.Addr().Interface().(fmt.Stringer)
		}

		if ok && !(slices.Contains([]reflect.Kind{reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Pointer,
			reflect.Slice,
			reflect.UnsafePointer}, s.Value.Kind()) && s.Value.IsNil()) {
			return stringer.String()
		}
	}

	return cast.MustAs[string](val)
}

func (s *reflectValue) Set(val string) error {
	if v, err := strconv.Unquote(val); err == nil {
		val = v
	}

	s.Inputs = append(s.Inputs, val)

	var err1 error

	if s.Value.Type().Kind() == reflect.Slice {
		item := reflect.New(s.Value.Type().Elem()).Elem()
		itemSetter, _ := SetterForType(item.Type())

		if err := itemSetter(val, item); err == nil {
			s.Value = reflect.Append(s.Value, item)

			return nil
		}
	}

	if len(s.Inputs) == 1 {
		if err1 = s.Setter(val, s.Value); err1 == nil {
			return nil
		}
	}

	if err2 := s.Setter(s.Inputs, s.Value); err2 != nil {
		if err1 != nil {
			return fmt.Errorf("can't set value %v to flag of type %s: %w", val, s.Value.Type().String(), err1)
		}

		return fmt.Errorf("can't set value %v to flag of type %s: %w", s.Inputs, s.Value.Type().String(), err2)
	}

	return nil
}

func (s *reflectValue) Type() string {
	return s.Value.Type().String()
}

func sliceSetter(val any, dst reflect.Value) error {
	items, err := cast.As[[]any](val)
	if err != nil {
		return err
	}

	item := reflect.New(dst.Type().Elem()).Elem()
	itemSetter, _ := SetterForType(item.Type())

	for _, itemVal := range items {
		if err := itemSetter(itemVal, item); err != nil {
			return err
		}

		dst.Set(reflect.Append(dst, item))
	}

	return nil
}

func DefaultSetter(val any, dst reflect.Value) error {
	if rv, ok := dst.Interface().(reflect.Value); ok {
		dst = rv
	}

	if dst.Type().Kind() == reflect.Slice {
		if err := sliceSetter(val, dst); err == nil {
			return nil
		}
	}

	return cast.ConvertToReflectValue(val, dst, cast.DefaultOptions)
}

func ReflectValueVarPWithSetter(setter setterFnc, stringerPref bool) func(dst *reflect.Value, name, shorthand string, defValue reflect.Value, usage string) {
	return func(dst *reflect.Value, name, shorthand string, defValue reflect.Value, usage string) {

		if defValue.IsValid() {
			dst.Set(defValue)
		}

		pflag.CommandLine.VarP(&reflectValue{
			Value:        *dst,
			Inputs:       nil,
			Setter:       setter,
			StringerPref: stringerPref,
		}, name, shorthand, usage)
	}

}
