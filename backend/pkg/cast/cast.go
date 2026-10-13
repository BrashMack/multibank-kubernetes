// Package cast - умное преобразование любого типа в любой другой.
package cast

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unsafe"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/exp/constraints"
)

// Options - опциональные настройки конертации.
type Options struct {
	UseDefaultTag  bool
	UseJSONTag     bool
	NameNormalizer func(string) string
	NilStringValue *string
	TimeLayout     string
	StrictBindMode bool
}

func (s Options) WithStrictBindMode() Options {
	s.StrictBindMode = true

	return s
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type interfaceHeader struct {
	Type unsafe.Pointer
	Data unsafe.Pointer
}

type customConverter struct {
	srcType   reflect.Type
	convertTo func(src any, value reflect.Value, opt Options) error
}

var (
	// DefaultOptions - опции понвертации используемые по умолчанию.
	DefaultOptions = Options{
		UseDefaultTag:  true,
		UseJSONTag:     true,
		NameNormalizer: DefaultNameNormalizer,
		NilStringValue: &DefaultNilStringValue,
		TimeLayout:     time.RFC3339Nano,
		StrictBindMode: false,
	}
	// ErrCastFailure - преобразование не удалось.
	ErrCastFailure = errors.New("can't cast")
	// ErrNilValue - не удалось записать nil в тип который такое значение не поддерживает.
	ErrNilValue = errors.New("can't assign nil value to nonrefference type")
	// ErrTypeNotSupported - преобразование между запрошеными типами не поддерживается.
	ErrTypeNotSupported = errors.New("type not supported")
	// ErrDstNotAddressable - назначение преобразвание должно меть тип указателя или быть адресуемым.
	ErrDstNotAddressable = errors.New("destination value mast be pointer or addressable")
	//ErrNoBindInStrictMode - выдается если не все поля исхоного объекта существуют в целевой структуре.
	ErrExtraSrcKeyInStrictMode = errors.New("can't bind source key in StrictBindMode")
	//ErrNoBindInStrictMode - выдается если не все поля исхоного объекта существуют в целевой структуре.
	ErrExtraDstKeyInStrictMode = errors.New("can't bind destination key in StrictBindMode")
	// DefaultNilStringValue - предствление nil в виде строки.
	DefaultNilStringValue = fmt.Sprint(nil)
	customConvertors      = make(map[reflect.Type]*customConverter)
	customConvertorsMutex sync.RWMutex
)

// CustomCaster - функция реализующая кастомный конвертер в нужный тип.
type CustomCaster[S, D any] func(src S, opt Options) (D, error)

// Caster - интерфейс который дожен реализовать тип для объявления кастомного преобразователя.
type Caster[SourceType any] interface {
	CastFrom(SourceType, Options) error
}

// WrapCastError - если ошибка не nil  - то оборачивает ее текст в ErrCastFailure,
// иначе выдает результат равный ргументам.
func WrapCastError[T any](v T, err error) (res T, resErr error) { //nolint:ireturn // generic type
	if err != nil {
		return res, fmt.Errorf("%w: %s", ErrCastFailure, err.Error())
	}

	return v, nil
}

// MustRegisterCaster - регистрирует кастомный кастер или бросает панику.
func MustRegisterCaster[S, D any](caster CustomCaster[S, D]) {
	customConvertorsMutex.Lock()
	defer customConvertorsMutex.Unlock()

	srcType := reflect.TypeOf((*S)(nil)).Elem()
	dstType := reflect.TypeOf((*D)(nil)).Elem()

	customConvertors[dstType] = &customConverter{
		srcType: srcType,
		convertTo: func(src any, value reflect.Value, opt Options) (err error) {
			*value.Addr().Interface().(*D), err = caster(src.(S), opt)

			return err
		}}
}

// normalizeToNoPointer - выдает тип без указател и подсчитывает число удаленных указаталей.
func normalizeToNoPointer(tp reflect.Type) (res reflect.Type, qty int) {
	for tp.Kind() == reflect.Pointer {
		qty++

		tp = tp.Elem()
	}

	return tp, qty
}

func getCustomConvertor(dstType reflect.Type) (*customConverter, bool) {
	customConvertorsMutex.RLock()
	defer customConvertorsMutex.RUnlock()

	c, f := customConvertors[dstType]

	return c, f
}

var (
	casterType       = reflect.TypeOf((*Caster[struct{}])(nil)).Elem()
	casterMethodName = casterType.Method(0).Name
	errType          = reflect.TypeOf((*error)(nil)).Elem()
	optType          = reflect.TypeOf((*Options)(nil)).Elem()
)

func tryAsCaster(normalizedType reflect.Type) (res *customConverter) {
	dstType := reflect.PtrTo(normalizedType)

	method, f := dstType.MethodByName(casterMethodName)
	if !f ||
		method.Type.NumOut() != 1 ||
		method.Type.Out(0) != errType ||
		method.Type.NumIn() != 3 ||
		method.Type.In(2) != optType {
		return nil // not found or wrong signature
	}

	return &customConverter{
		srcType: method.Type.In(1),
		convertTo: func(src any, value reflect.Value, opt Options) error {
			if !value.CanAddr() {
				return ErrDstNotAddressable
			}
			value = value.Addr()

			err := method.Func.Call([]reflect.Value{value, reflect.ValueOf(src), reflect.ValueOf(opt)})[0].Interface()
			if err == nil {
				return nil
			}

			return err.(error) //nolint:forcetypeassert //checked by sources
		}}
}

// ReflectType выдает reflect.Type для желаемого типа.
func ReflectType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func findCustomConverter(dstType reflect.Type) (converter *customConverter) {
	var found bool

	if converter, found = getCustomConvertor(dstType); found {
		return converter
	}

	defer func() {
		customConvertorsMutex.Lock()
		defer customConvertorsMutex.Unlock()

		customConvertors[dstType] = converter
	}()

	normalizedType, ptrQty := normalizeToNoPointer(dstType)
	convNormalized, fndNormalized := getCustomConvertor(normalizedType)

	if !fndNormalized {
		defer func() {
			customConvertorsMutex.Lock()
			defer customConvertorsMutex.Unlock()

			customConvertors[normalizedType] = convNormalized
		}()

		convNormalized = tryAsCaster(normalizedType)
	}

	if convNormalized == nil {
		return nil
	}

	return &customConverter{ // generate pointer to normalized converter
		srcType: convNormalized.srcType,
		convertTo: func(src any, value reflect.Value, opt Options) error {
			tmpVal := reflect.New(normalizedType).Elem()
			err := convNormalized.convertTo(src, tmpVal, opt)

			if err != nil {
				return err
			}

			for ptrQty > 0 {
				value.Set(reflect.New(value.Type().Elem()))
				value = value.Elem()
				ptrQty--
			}
			value.Set(tmpVal)

			return nil
		}}
}

// DefaultNameNormalizer - значение по умолчанию для Options.NameNormalizer
// позволяет сопоставить названия полей с полями структуры.
func DefaultNameNormalizer(name string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r):
			return unicode.ToLower(r)
		default:
			return -1
		}
	}, name)
}

// As - преобразует любое значение в желаемый тип или возвращает ошибку.
func As[T any](src any) (T, error) { // nolint: ireturn // generic type
	return AsOpt[T](src, DefaultOptions)
}

// AsOpt - преобразует любое значение  с указыными опцими  в желаемый тип или возвращает ошибку.
func AsOpt[T any](src any, opt Options) (res T, err error) { // nolint: ireturn //generic type
	return res, ConvertTo(src, &res, opt)
}

// MustAs - преобразует любое значение в желаемый тип или кидает панику.
func MustAs[T any](src any) T { // nolint: ireturn //generic type
	res, err := As[T](src)

	if err != nil {
		panic(err)
	}

	return res
}

// MustAsOpt - преобразует любое значение с указыными опцими в желаемый тип или кидает панику.
func MustAsOpt[T any](src any, opt Options) T { // nolint: ireturn //generic type
	res, err := AsOpt[T](src, opt)

	if err != nil {
		panic(err)
	}

	return res
}

func derefference(v any, opt Options) any {
	if v == nil {
		return nil
	}

	var kd reflect.Kind

	rv := reflect.ValueOf(v)

	for {
		kd = rv.Type().Kind()

		switch kd { // nolint: exhaustive // no need all elements check
		case reflect.Invalid:
			return nil
		case reflect.String:
			ns := opt.NilStringValue

			if ns != nil && rv.String() == *ns {
				return nil
			}
		case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
			if rv.IsNil() {
				return nil
			}
		}

		if kd == reflect.Pointer {
			rv = rv.Elem()

			continue
		}

		return rv.Interface()
	}
}

func setIfNoError[T any](dst reflect.Value, src result[T]) error {
	if src.Error != nil {
		return src.Error
	}

	set(dst, src.Value)

	return nil
}

func set[T any](dst reflect.Value, v T) {
	*(*T)(dst.Addr().UnsafePointer()) = v
}

// MustSet - записывает значение в указанный приемник или кидает панику если тип приемника не совпадает.
func MustSet[T any](dst reflect.Value, v T) {
	*MustAddressOf[T](dst) = v
}

// MustAddressOf - преобразет приемник типа reflect.Value в типизированный указатель на это значение.
func MustAddressOf[T any](dst reflect.Value) *T {
	st := dst.Type()
	dt := reflect.TypeOf((*T)(nil)).Elem()

	if st.Kind() != dt.Kind() || st.Size() != dt.Size() {
		panic(fmt.Sprintf("can't use *%v as *%v", st, dt))
	}

	return (*T)(dst.Addr().UnsafePointer())
}

func setBytes(dst reflect.Value, src unsafe.Pointer, count int) {
	srcBytes := *(*[]byte)(unsafe.Pointer(&sliceHeader{Data: src, Cap: count, Len: count}))
	dstBytes := *(*[]byte)(unsafe.Pointer(&sliceHeader{Data: dst.Addr().UnsafePointer(), Cap: count, Len: count}))

	copy(dstBytes, srcBytes)
}

type result[T any] struct {
	Value T
	Error error
}

func valueKind(src any) (reflect.Value, reflect.Kind) {
	v := reflect.ValueOf(src)

	if src == nil {
		return v, reflect.Invalid
	}

	return v, v.Type().Kind()
}

func newResult[T any](value T, err error) result[T] {
	return result[T]{Value: value, Error: err}
}

func toBool(src any) result[bool] { // nolint: revive // can't simplify
	switch sv, tp := valueKind(src); tp {
	case reflect.Invalid:
		return newResult(false, ErrNilValue)
	case reflect.Bool:
		return result[bool]{Value: sv.Bool(), Error: nil}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch sv.Int() {
		case 0:
			return result[bool]{Value: false, Error: nil}
		case 1:
			return result[bool]{Value: true, Error: nil}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch sv.Uint() {
		case 0:
			return result[bool]{Value: false, Error: nil}
		case 1:
			return result[bool]{Value: true, Error: nil}
		}
	case reflect.Float32, reflect.Float64:
		switch sv.Float() {
		case 0:
			return result[bool]{Value: false, Error: nil}
		case 1:
			return result[bool]{Value: true, Error: nil}
		}
	case reflect.String:
		switch sv.String() {
		case "1", "true", "TRUE", "True", "yes", "YES":
			return result[bool]{Value: true, Error: nil}
		case "0", "false", "FALSE", "False", "no", "NO":
			return result[bool]{Value: false, Error: nil}
		}
	case reflect.Array, reflect.Chan,
		reflect.Complex128, reflect.Complex64,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer,
		reflect.Slice, reflect.Struct, reflect.UnsafePointer:
	}

	return newResult(false, ErrCastFailure)
}

func toInt[T constraints.Signed](src any) result[T] {
	switch sv, tp := valueKind(src); tp {
	case reflect.Invalid:
		return newResult[T](0, ErrNilValue)
	case reflect.Bool:
		var v int64
		if sv.Bool() {
			v = 1
		}

		return newResult(castInt[T](v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newResult(castInt[T](sv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, err := uInt2Int(sv.Uint())
		if err != nil {
			return newResult[T](0, err)
		}

		return newResult(castInt[T](v))
	case reflect.Float32, reflect.Float64:
		f := sv.Float()
		v := int64(f)

		if float64(v) != f {
			return newResult[T](0, ErrCastFailure)
		}

		return newResult(castInt[T](v))
	case reflect.String:
		v, err := strconv.ParseInt(sv.String(), 10, reflect.TypeOf(T(0)).Bits())

		if err != nil {
			return newResult[T](0, ErrCastFailure)
		}

		return newResult(castInt[T](v))
	case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Func,
		reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
	}

	return newResult[T](0, ErrCastFailure)
}

func toUint[T constraints.Unsigned](src any) result[T] {
	switch sv, tp := valueKind(src); tp {
	case reflect.Invalid:
		return newResult[T](0, ErrNilValue)
	case reflect.Bool:
		var v uint64
		if sv.Bool() {
			v = 1
		}

		return newResult(castUInt[T](v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := int2Uint(sv.Int())
		if err != nil {
			return newResult[T](0, err)
		}

		return newResult(castUInt[T](v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return newResult(castUInt[T](sv.Uint()))
	case reflect.Float32, reflect.Float64:
		f := sv.Float()
		v := uint64(f)

		if float64(v) != f {
			return newResult[T](0, ErrCastFailure)
		}

		return newResult(castUInt[T](v))
	case reflect.String:
		v, err := strconv.ParseUint(sv.String(), 10, reflect.TypeOf(T(0)).Bits())

		if err != nil {
			return newResult[T](0, ErrCastFailure)
		}

		return newResult(castUInt[T](v))
	case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
	}

	return newResult[T](0, ErrCastFailure)
}

func toString(src any, opt Options) result[string] {
	switch val := src.(type) {
	case time.Time:
		return newResult(val.Format(DefaultOptions.TimeLayout), nil)
	case []byte:
		return newResult(base64.StdEncoding.EncodeToString(val), nil)
	case fmt.Stringer:
		return newResult(val.String(), nil)
	}

	switch sv, tp := valueKind(src); tp {
	case reflect.Invalid:
		if ns := opt.NilStringValue; ns != nil {
			return newResult(*ns, nil)
		}

		return newResult("", ErrNilValue)
	case reflect.Array, reflect.Interface, reflect.Map, reflect.Slice, reflect.Struct:
		res, err := json.Marshal(sv.Interface())
		if err != nil {
			return newResult("", ErrCastFailure)
		}

		return newResult(string(res), nil)
	case reflect.String:
		return newResult(sv.String(), nil)
	case reflect.Bool, reflect.Complex128, reflect.Complex64, reflect.Int,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return newResult(fmt.Sprintf("%v", src), nil)
	case reflect.Chan, reflect.Uintptr, reflect.Pointer, reflect.Func, reflect.UnsafePointer:
	}

	return newResult("", ErrCastFailure)
}

func toFloat[T constraints.Float](src any) result[T] {
	switch sv, tp := valueKind(src); tp {
	case reflect.Invalid:
		return newResult[T](0, ErrNilValue)
	case reflect.Bool:
		var v float64
		if sv.Bool() {
			v = 1
		}

		return newResult(T(v), nil)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newResult(T(sv.Int()), nil)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return newResult(T(sv.Uint()), nil)
	case reflect.Float32, reflect.Float64:
		return newResult(T(sv.Float()), nil)
	case reflect.String:
		v, err := strconv.ParseFloat(sv.String(), reflect.TypeOf(T(0)).Bits())

		if err != nil {
			return newResult[T](0, ErrCastFailure)
		}

		return newResult(T(v), nil)
	case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
	}

	return newResult[T](0, ErrCastFailure)
}

func toComplex[T constraints.Complex](src any) result[T] {
	switch sv, tp := valueKind(src); tp {
	case reflect.Invalid:
		return newResult[T](0, ErrNilValue)
	case reflect.Complex64, reflect.Complex128:
		return newResult(T(sv.Complex()), ErrCastFailure)
	case reflect.String:
		v, err := strconv.ParseComplex(sv.String(), reflect.TypeOf(T(0)).Bits())

		if err != nil {
			return newResult[T](0, ErrCastFailure)
		}

		return newResult(T(v), nil)
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Float32, reflect.Float64, reflect.Func, reflect.Int,
		reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface, reflect.Map, reflect.Pointer,
		reflect.Slice, reflect.Struct, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8,
		reflect.Uintptr, reflect.UnsafePointer:
	}

	return newResult[T](0, ErrCastFailure)
}

func toSlice(src any, elem reflect.Type, opt Options) result[sliceHeader] {
	if src == nil {
		return result[sliceHeader]{Value: sliceHeader{Data: nil, Len: 0, Cap: 0}, Error: nil}
	}

	switch sv, tp := valueKind(src); tp {
	case reflect.String:
		var (
			list []interface{}
			err  error
		)

		if err = json.Unmarshal([]byte(sv.String()), &list); err != nil {
			return newResult(sliceHeader{}, ErrCastFailure)
		}

		res := reflect.MakeSlice(reflect.SliceOf(elem), len(list), len(list))

		for i, elemValue := range list {
			dst := reflect.New(elem).Elem()

			if err = ConvertToReflectValue(derefference(elemValue, opt), dst, opt); err != nil {
				return newResult(sliceHeader{}, ErrCastFailure)
			}

			res.Index(i).Set(dst)
		}

		return newResult(sliceHeader{Data: res.UnsafePointer(), Len: res.Len(), Cap: res.Cap()}, nil)
	case reflect.Array, reflect.Slice:
		l := sv.Len()
		res := reflect.MakeSlice(reflect.SliceOf(elem), l, l)

		var (
			err       error
			elemValue interface{}
		)

		for i := 0; i < l; i++ {
			elemValue = sv.Index(i).Interface()
			dst := reflect.New(elem).Elem()

			if err = ConvertToReflectValue(derefference(elemValue, opt), dst, opt); err != nil {
				return newResult(sliceHeader{}, ErrCastFailure)
			}

			res.Index(i).Set(dst)
		}

		return newResult(sliceHeader{Data: res.UnsafePointer(), Len: res.Len(), Cap: res.Cap()}, nil)
	case reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64,
		reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface,
		reflect.Invalid, reflect.Map, reflect.Pointer, reflect.Struct, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
	}

	return newResult(sliceHeader{}, ErrCastFailure)
}

var anyType = reflect.TypeOf((*interface{})(nil)).Elem()

func toInterface(src any, dstTp reflect.Type) result[interfaceHeader] {
	if src == nil {
		return result[interfaceHeader]{Value: interfaceHeader{Data: nil, Type: nil}, Error: nil}
	}

	if dstTp == anyType {
		return result[interfaceHeader]{Value: *(*interfaceHeader)(unsafe.Pointer(&src)), Error: nil}
	}

	return newResult(interfaceHeader{}, ErrCastFailure)
}

func toMap(src any, // nolint: revive //can't simplify
	key reflect.Type,
	elem reflect.Type,
	opt Options) result[unsafe.Pointer] {
	if src == nil {
		return result[unsafe.Pointer]{Value: nil, Error: nil}
	}

	switch sv, tp := valueKind(src); tp {
	case reflect.String:
		var (
			tmpMap map[string]interface{}
			err    error
		)

		if err = json.Unmarshal([]byte(sv.String()), &tmpMap); err != nil {
			return newResult[unsafe.Pointer](nil, ErrCastFailure)
		}

		res := reflect.MakeMap(reflect.MapOf(key, elem))

		for srcKey, srcVal := range tmpMap {
			dstKey := reflect.New(key).Elem()
			dstVal := reflect.New(elem).Elem()

			if err = ConvertToReflectValue(derefference(srcKey, opt), dstKey, opt); err != nil {
				return newResult[unsafe.Pointer](nil, ErrCastFailure)
			}

			if err = ConvertToReflectValue(derefference(srcVal, opt), dstVal, opt); err != nil {
				return newResult[unsafe.Pointer](nil, ErrCastFailure)
			}

			res.SetMapIndex(dstKey, dstVal)
		}

		return newResult(res.UnsafePointer(), nil)
	case reflect.Struct:
		res := reflect.MakeMap(reflect.MapOf(key, elem))
		tp := sv.Type()
		l := tp.NumField()

		var err error

		for i := 0; i < l; i++ {
			fi := tp.Field(i)

			if !fi.IsExported() {
				continue
			}

			keyDst := reflect.New(key).Elem()
			valDst := reflect.New(elem).Elem()

			if err = ConvertToReflectValue(opt.NameNormalizer(fi.Name), keyDst, opt); err != nil {
				return newResult[unsafe.Pointer](nil, ErrCastFailure)
			}

			if err = ConvertToReflectValue(sv.Field(i).Interface(), valDst, opt); err != nil {
				return newResult[unsafe.Pointer](nil, ErrCastFailure)
			}

			res.SetMapIndex(keyDst, valDst)
		}

		return newResult(res.UnsafePointer(), nil)
	case reflect.Map:
		res := reflect.MakeMap(reflect.MapOf(key, elem))
		iter := sv.MapRange()

		var err error

		for iter.Next() {
			dstKey := reflect.New(key).Elem()
			dstVal := reflect.New(elem).Elem()

			if err = ConvertToReflectValue(derefference(iter.Key().Interface(), opt), dstKey, opt); err != nil {
				return newResult[unsafe.Pointer](nil, ErrCastFailure)
			}

			if err = ConvertToReflectValue(derefference(iter.Value().Interface(), opt), dstVal, opt); err != nil {
				return newResult[unsafe.Pointer](nil, ErrCastFailure)
			}

			res.SetMapIndex(dstKey, dstVal)
		}

		return newResult(res.UnsafePointer(), nil)
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Pointer, reflect.Slice, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
	}

	return newResult[unsafe.Pointer](nil, ErrCastFailure)
}

func setStruct(src any, dst reflect.Value, opt Options) error { // nolint: revive //can't simplify
	dstTp := dst.Type()
	bindStatus := make([]bool, dstTp.NumField())

	dstFieldAliases := make(map[string]int)

	findDstFiled := func(name string) (reflect.Value, bool) {
		index, found := dstFieldAliases[opt.NameNormalizer(name)]
		if !found {
			return reflect.Value{}, false
		}

		bindStatus[index] = true

		return dst.Field(index), true

	}

	fieldDefault := make(map[reflect.Value]string)

	var extraSourceKeys []string

	setDefaults := func() error {
		for fv, defValue := range fieldDefault {
			if err := ConvertToReflectValue(defValue, fv, opt); err != nil {
				return err
			}
		}

		var err error

		if opt.StrictBindMode {
			for index, ok := range bindStatus {
				if !ok {
					err = multierror.Append(err, fmt.Errorf("%w, key: %s", ErrExtraDstKeyInStrictMode, dstTp.Field(index).Name))
				}
			}

			for _, name := range extraSourceKeys {
				err = multierror.Append(err, fmt.Errorf("%w, key: %s", ErrExtraSrcKeyInStrictMode, name))
			}
		}

		return err
	}

	for dstFiledIndex, l := 0, dstTp.NumField(); dstFiledIndex < l; dstFiledIndex++ {
		field := dst.Field(dstFiledIndex)

		if !field.CanSet() {
			bindStatus[dstFiledIndex] = true
			continue
		}

		fi := dstTp.Field(dstFiledIndex)

		if opt.UseJSONTag {
			if name, ok := fi.Tag.Lookup("json"); ok {
				dstFieldAliases[opt.NameNormalizer(name)] = dstFiledIndex
			}
		}

		if opt.UseDefaultTag {
			if defValue, ok := fi.Tag.Lookup("default"); ok {
				fieldDefault[field] = defValue
				bindStatus[dstFiledIndex] = true
			}
		}

		dstFieldAliases[opt.NameNormalizer(fi.Name)] = dstFiledIndex
	}

	switch sv, tp := valueKind(src); tp {
	case reflect.String:
		var newSRC map[string]any

		if err := ConvertToReflectValue(sv.String(), reflect.ValueOf(&newSRC).Elem(), opt); err != nil {
			return err
		}

		return setStruct(newSRC, dst, opt)
	case reflect.Struct:
		l := sv.NumField()
		tp := sv.Type()

		for i := 0; i < l; i++ {
			sf := sv.Field(i)
			tf := tp.Field(i)

			if !tf.IsExported() {
				continue
			}

			if fv, ok := findDstFiled(tf.Name); ok {
				delete(fieldDefault, fv)

				if err := ConvertToReflectValue(sf.Interface(), fv, opt); err != nil {
					return err
				}
			} else {
				extraSourceKeys = append(extraSourceKeys, tf.Name)
			}
		}

		return setDefaults()
	case reflect.Map:
		iter := sv.MapRange()
		for iter.Next() {
			var keyName string

			if err := ConvertToReflectValue(iter.Key().Interface(), reflect.ValueOf(&keyName).Elem(), opt); err != nil {
				return err
			}

			if fv, ok := findDstFiled(keyName); ok {
				delete(fieldDefault, fv)

				if err := ConvertToReflectValue(iter.Value().Interface(), fv, opt); err != nil {
					return err
				}
			} else {
				extraSourceKeys = append(extraSourceKeys, keyName)
			}
		}

		return setDefaults()
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Pointer, reflect.Slice, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
	}

	return ErrCastFailure
}

// AsReflectType преобразует значение в нужный reflect.Type и возвращает результат в reflect.Value.
func AsReflectType(src any, tp reflect.Type, opt Options) (reflect.Value, error) {
	val := reflect.New(tp).Elem()

	if err := ConvertToReflectValue(src, val, opt); err != nil {
		return reflect.Value{}, err
	}

	return val, nil
}

// ConvertTo преобразуе и сораняет значение в указанный приемник.
func ConvertTo[T any](src any, dst *T, opt Options) (err error) {
	dstRv := reflect.ValueOf(dst).Elem()

	return ConvertToReflectValue(src, dstRv, opt)
}

// ConvertToReflectValue преобразуе и сораняет значение в reflect.Value.
func ConvertToReflectValue(src any, // nolint: revive, gocyclo //can't simplify
	dstRv reflect.Value,
	opt Options) (err error) {
	src = derefference(src, opt)

	defer func(s any, d reflect.Value) {
		if err != nil {
			err = fmt.Errorf("cast value %v of type %T to type %v error: %w", s, s, d.Type(), err)
		}
	}(src, dstRv)

	dstTp := dstRv.Type()

	if converter := findCustomConverter(dstTp); converter != nil {
		if reflect.TypeOf(src) != converter.srcType {
			rv, err := AsReflectType(src, converter.srcType, opt)
			if err != nil {
				return fmt.Errorf("can't prepare source for custom converter: %w", err)
			}

			src = rv.Interface()
		}

		return converter.convertTo(src, dstRv, opt)
	}

	switch dstTp.Kind() {
	case reflect.Bool:
		return setIfNoError(dstRv, toBool(src))
	case reflect.Int:
		return setIfNoError(dstRv, toInt[int](src))
	case reflect.Int8:
		return setIfNoError(dstRv, toInt[int8](src))
	case reflect.Int16:
		return setIfNoError(dstRv, toInt[int16](src))
	case reflect.Int32:
		return setIfNoError(dstRv, toInt[int32](src))
	case reflect.Int64:
		return setIfNoError(dstRv, toInt[int64](src))
	case reflect.Uint:
		return setIfNoError(dstRv, toUint[uint](src))
	case reflect.Uint8:
		return setIfNoError(dstRv, toUint[uint8](src))
	case reflect.Uint16:
		return setIfNoError(dstRv, toUint[uint16](src))
	case reflect.Uint32:
		return setIfNoError(dstRv, toUint[uint32](src))
	case reflect.Uint64:
		return setIfNoError(dstRv, toUint[uint64](src))
	case reflect.Uintptr:
		return setIfNoError(dstRv, toUint[uintptr](src))
	case reflect.Float32:
		return setIfNoError(dstRv, toFloat[float32](src))
	case reflect.Float64:
		return setIfNoError(dstRv, toFloat[float64](src))
	case reflect.Complex64:
		return setIfNoError(dstRv, toComplex[complex64](src))
	case reflect.Complex128:
		return setIfNoError(dstRv, toComplex[complex128](src))
	case reflect.String:
		return setIfNoError(dstRv, toString(src, opt))
	case reflect.Array:
		slice := toSlice(src, dstTp.Elem(), opt)
		if slice.Error != nil {
			return slice.Error
		}

		if slice.Value.Data == nil {
			return fmt.Errorf("%w %v", ErrNilValue, dstTp)
		}

		if dstRv.Len() != slice.Value.Len {
			return fmt.Errorf("%w. Array length don't matched %v, but need %d elemets",
				ErrCastFailure, dstTp, slice.Value.Len)
		}

		setBytes(dstRv, slice.Value.Data, slice.Value.Len*int(dstTp.Elem().Size()))

		return nil
	case reflect.Map:
		return setIfNoError(dstRv, toMap(src, dstTp.Key(), dstTp.Elem(), opt))
	case reflect.Pointer:
		if src == nil {
			set[unsafe.Pointer](dstRv, nil)

			return nil
		}

		ptrRes := reflect.New(dstTp.Elem())
		set(dstRv, ptrRes.UnsafePointer())

		return ConvertToReflectValue(src, ptrRes.Elem(), opt)
	case reflect.Slice:
		return setIfNoError(dstRv, toSlice(src, dstTp.Elem(), opt))
	case reflect.Struct:
		return setStruct(src, dstRv, opt)
	case reflect.Interface:
		return setIfNoError(dstRv, toInterface(src, dstTp))
	case reflect.Chan, reflect.Func, reflect.Invalid, reflect.UnsafePointer:
	}

	return ErrTypeNotSupported
}

func castInt[T constraints.Signed](v int64) (res T, err error) { // nolint: ireturn //generic type
	bits := unsafe.Sizeof(res) << 3
	maxInt := int64((uint64(1) << (bits - 1)) - 1)
	minInt := -maxInt - 1

	if v < minInt || v > maxInt {
		return 0, fmt.Errorf("can't convert %v to %T, %w", v, res, ErrCastFailure)
	}

	return T(v), nil
}

func castUInt[T constraints.Unsigned](v uint64) (res T, err error) { // nolint: ireturn //generic type
	if v > uint64(^T(0)) {
		return 0, fmt.Errorf("can't convert %v to %T, %w", v, res, ErrCastFailure)
	}

	return T(v), nil
}

const int64Max = ^uint64(0) >> 1

func uInt2Int(v uint64) (res int64, err error) {
	if v > int64Max {
		return 0, fmt.Errorf("can't convert %v to %T, %w", v, res, ErrCastFailure)
	}

	return int64(v), nil
}

func int2Uint(v int64) (res uint64, err error) {
	if v < 0 {
		return 0, fmt.Errorf("can't convert %v to %T, %w", v, res, ErrCastFailure)
	}

	return uint64(v), nil
}
