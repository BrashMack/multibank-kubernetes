// Package loadconfig - позволяет загрузкить конфиг файл в глобальную переременную нужного типа.
package loadconfig

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"reflect"
	"slices"
	"time"
	"unsafe"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	et "github.com/go-playground/validator/v10/translations/en"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/masharpik/vtb-hack-back/pkg/cast"
)

var (
	// ErrNotStructDst - неверный тип приемника конфига.
	ErrNotStructDst = errors.New("dst must be pointer to struct")
	// ErrUnsupportedType - тип приемника не поддерживается.
	ErrUnsupportedType = errors.New("unsupported type")
	validate           = validator.New()
)

type setterFnc func(val any, dst reflect.Value) error

type paramValue struct {
	name  string
	value any
}

// ValidatorError - Тип ошибки валидации поля.
type ValidatorError struct {
	Message string
	Key     string
}

func (s ValidatorError) Error() string {
	return fmt.Sprintf("%s: %s", s.Key, s.Message)
}

// MustTo - парсит конфиг в указанный приемник либо кидает панику.
func MustTo[T any](dst *T) {
	if err := To(dst); err != nil {
		panic(err)
	}
}

//nolint:revive // To - парсит конфиг в переменную нужного типа.
func To[T any](dst *T) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v\n", err)
		}
	}()

	return ToCtx(context.Background(), dst)
}

//nolint:revive // To - парсит конфиг в переменную нужного типа с поддержкой контекста.
func ToCtx[T any](ctx context.Context, dst *T) error {
	if dst == nil {
		return ErrNotStructDst
	}

	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	dstRV := reflect.ValueOf(dst).Elem()

	if dstRV.Kind() != reflect.Struct {
		return ErrNotStructDst
	}

	allParams := make(params)

	if err := allParams.defineFlagsForStruct(dstRV, ""); err != nil {
		return err
	}
	viper.AutomaticEnv()
	viper.BindEnv("CONFIG")
	viper.BindEnv("SECRETS")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("can't bind pflags: %w", err)
	}
	if path := viper.GetString("config"); len(path) > 0 {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigFile("config.yaml")
	}

	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("can't load config: %w", err)
	}

	if secretsPath := viper.GetString("secrets"); len(secretsPath) > 0 {
		viper.SetConfigFile(secretsPath)
	} else {
		viper.SetConfigFile("secrets.yaml")
	}

	if err := viper.MergeInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("can't load secrets: %w", err)
	}

	structsForValidate := []paramValue{{name: "", value: dstRV}}

	keys := make([]string, 0, len(allParams))

	for k := range allParams {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, n := range keys {
		v := allParams[n]

		if v.isStruct {
			structsForValidate = append(structsForValidate, paramValue{name: v.name, value: v.value.Addr().Interface()})
		}

		if !v.registred {
			continue
		}

		val := viper.Get(n)

		if val == nil {
			continue
		}

		if err := v.setValue(val, v.value); err != nil {
			return fmt.Errorf("can't set value %v of key %s: %w", val, n, err)
		}
	}

	var (
		errList error
		errMap  = make(map[string]struct{})
	)

	for _, sv := range structsForValidate {
		err := validate.StructCtx(ctx, sv.value)
		if err == nil {
			continue
		}

		var errs validator.ValidationErrors

		if errors.As(err, &errs) {
			for key, val := range errs.Translate(trans) {
				if sv.name != "" {
					key = sv.name + "." + key
				}

				if _, exists := errMap[val]; !exists {
					errMap[val] = struct{}{}
					errList = multierror.Append(errList, ValidatorError{Message: val, Key: key})
				}
			}

			continue
		}

		panic(err)
	}

	if errList != nil {
		return fmt.Errorf("config validation error: %w", errList)
	}

	return nil
}

func defaultFieldSetter[T any](val any, dst reflect.Value) error {
	v, err := cast.As[T](val)
	if err != nil {
		return err
	}

	dst.Set(reflect.ValueOf(v).Convert(dst.Type()))

	return nil
}

var reflectValueType = reflect.TypeOf(reflect.Value{})

func registerFlag[T any](
	call func(p *T, name, shorthand string, value T, usage string), fi *flagInfo,
	setter ...setterFnc,
) error {
	var defTypedValue T

	if fi.registred {
		return nil
	}

	fi.registred = true

	if reflect.TypeOf(defTypedValue) == reflectValueType { //real type from fi.value.Type()
		*(*reflect.Value)(unsafe.Pointer(&defTypedValue)) = reflect.New(fi.value.Type()).Elem()
	}

	switch len(setter) {
	case 0:
		fi.setValue = defaultFieldSetter[T]
	case 1:
		fi.setValue = setter[0]
	default:
		panic("registerFlag allow only one setter")
	}

	if fi.defaultValue != nil {
		if err := fi.setValue(*fi.defaultValue, reflect.ValueOf(&defTypedValue).Elem()); err != nil {
			return fmt.Errorf("can't parse default value param %s: %w", fi.name, err)
		}
	}

	var ptr *T

	if reflect.TypeOf(defTypedValue) == reflectValueType {
		*(**reflect.Value)(unsafe.Pointer(&ptr)) = &fi.value
	} else {
		ptr = cast.MustAddressOf[T](fi.value)
	}

	call(ptr, fi.name, fi.short, defTypedValue, fi.description)

	return nil
}

func (s params) defineFlagsForStruct(dst reflect.Value, basePath string) error {
	tp := dst.Type()
	fieldCount := tp.NumField()

	for fieldIndex := 0; fieldIndex < fieldCount; fieldIndex++ {
		fieldInfo := tp.Field(fieldIndex)
		fieldValue := dst.Field(fieldIndex)

		if !fieldInfo.IsExported() {
			continue
		}

		paramName := fieldInfo.Name

		if jn, ok := fieldInfo.Tag.Lookup("json"); ok {
			paramName = jn
		}

		if basePath != "" {
			paramName = basePath + "." + paramName
		}

		if err := s.bindField(fieldInfo, fieldValue, paramName, basePath); err != nil {
			return err
		}
	}

	return nil
}

type params map[string]*flagInfo

type flagInfo struct {
	name         string
	short        string
	defaultValue *string
	description  string
	value        reflect.Value
	isStruct     bool
	registred    bool
	setValue     setterFnc
}

var (
	typeDuration = cast.ReflectType[time.Duration]()
	typeIP       = cast.ReflectType[net.IP]()
	typeIPNet    = cast.ReflectType[net.IPNet]()
	typeIPMask   = cast.ReflectType[net.IPMask]()
)

//nolint:revive // can't simplify
func (s params) bindField(fieldInfo reflect.StructField, fieldValue reflect.Value, paramName, parent string) (err error) {
	fi := &flagInfo{
		name:         paramName,
		short:        "",
		defaultValue: nil,
		description:  "",
		value:        fieldValue,
		isStruct:     false,
		registred:    false,
	}

	fi.short, _ = fieldInfo.Tag.Lookup("short")
	fi.description, _ = fieldInfo.Tag.Lookup("description")

	if defaultValueStr, hasDefaultValue := fieldInfo.Tag.Lookup("default"); hasDefaultValue {
		fi.defaultValue = &defaultValueStr
	}

	defer func() {
		if err != nil {
			err = fmt.Errorf("bind config param %s to field %v error: %w", paramName, fieldInfo, err)

			return
		}

		s[paramName] = fi
	}()

	switch fieldInfo.Type {
	case typeDuration:
		return registerFlag(pflag.DurationVarP, fi)
	case typeIP:
		return registerFlag(pflag.IPVarP, fi)
	case typeIPNet:
		return registerFlag(pflag.IPNetVarP, fi)
	case typeIPMask:
		return registerFlag(pflag.IPMaskVarP, fi)
	default:
		tryRegisterByParser(fi) //Try find method Parse[T](string)(T, error)

		if fi.registred {
			return
		}
	}

	switch fieldInfo.Type.Kind() {
	case reflect.Bool:
		return registerFlag(pflag.BoolVarP, fi)
	case reflect.Int:
		return registerFlag(pflag.IntVarP, fi)
	case reflect.Int8:
		return registerFlag(pflag.Int8VarP, fi)
	case reflect.Int16:
		return registerFlag(pflag.Int16VarP, fi)
	case reflect.Int32:
		return registerFlag(pflag.Int32VarP, fi)
	case reflect.Int64:
		if fieldInfo.Type == typeDuration {
			return registerFlag(pflag.DurationVarP, fi)
		}

		return registerFlag(pflag.Int64VarP, fi)
	case reflect.Uint:
		return registerFlag(pflag.UintVarP, fi)
	case reflect.Uint8:
		return registerFlag(pflag.Uint8VarP, fi)
	case reflect.Uint16:
		return registerFlag(pflag.Uint16VarP, fi)
	case reflect.Uint32:
		return registerFlag(pflag.Uint32VarP, fi)
	case reflect.Uint64:
		return registerFlag(pflag.Uint64VarP, fi)
	case reflect.Float32:
		return registerFlag(pflag.Float32VarP, fi)
	case reflect.Float64:
		return registerFlag(pflag.Float64VarP, fi)
	case reflect.String:
		return registerFlag(pflag.StringVarP, fi)
	case reflect.Struct:
		for i := 0; i < fieldInfo.Type.NumField(); i++ {
			if fieldInfo.Type.Field(i).IsExported() {
				fi.isStruct = true //struct with one or more exported field
				break
			}
		}

		if !fi.registred && !fi.isStruct {
			break //use reflectValue type
		}

		return s.defineFlagsForStruct(fieldValue, paramName)
	case reflect.Slice:
		elem := fieldInfo.Type.Elem()
		switch elem {
		case typeDuration:
			return registerFlag(pflag.DurationSliceVarP, fi)
		case typeIP:
			return registerFlag(pflag.IPSliceVarP, fi)
		case typeIPNet:
			return registerFlag(pflag.IPNetSliceVarP, fi)
		}

		switch elem.Kind() {
		case reflect.Uint:
			return registerFlag(pflag.UintSliceVarP, fi)
		case reflect.String:
			return registerFlag(pflag.StringSliceVarP, fi)
		case reflect.Int:
			return registerFlag(pflag.IntSliceVarP, fi)
		case reflect.Int64:
			return registerFlag(pflag.Int64SliceVarP, fi)
		case reflect.Int32:
			return registerFlag(pflag.Int32SliceVarP, fi)
		case reflect.Float32:
			return registerFlag(pflag.Float32SliceVarP, fi)
		case reflect.Float64:
			return registerFlag(pflag.Float32SliceVarP, fi)
		case reflect.Bool:
			return registerFlag(pflag.BoolSliceVarP, fi)
		}
	case reflect.Array, reflect.Map, reflect.Pointer:
	case reflect.Chan,
		reflect.Complex128, // todo: implement
		reflect.Complex64,  // todo: implement
		reflect.Func,
		reflect.Interface, reflect.Invalid, reflect.Uintptr, reflect.UnsafePointer:
		fallthrough
	default:
		return fmt.Errorf("%w: %v", ErrUnsupportedType, fieldInfo.Type)
	}

	return registerFlag(ReflectValueVarPWithSetter(DefaultSetter, false), fi, DefaultSetter)
}

var trans ut.Translator

func init() {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ = uni.GetTranslator("en")

	if err := et.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}
}
