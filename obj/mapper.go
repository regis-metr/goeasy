package obj

import (
	"fmt"
	"reflect"
)

var MismatchTypeError error = fmt.Errorf("type mismatch")
var UnsupportedTypeError error = fmt.Errorf("type unsupported")

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(src interface{}, dst interface{}) error {
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	if dstValue.Type().Kind() != reflect.Pointer {
		return fmt.Errorf("destination should be pointer")
	}
	return m.mapValue(srcValue, dstValue.Elem())

}

func (m *Mapper) mapValue(src reflect.Value, dst reflect.Value) error {
	if src.Type().Kind() == reflect.Pointer {
		return m.mapValue(src.Elem(), dst)
	}

	switch dst.Type().Kind() {
	case reflect.Bool:
		if src.Type().Kind() != reflect.Bool {
			return MismatchTypeError
		}
		dst.SetBool(src.Bool())
	case reflect.Int:
		if src.Type().Kind() != reflect.Int {
			return MismatchTypeError
		}
		dst.SetInt(src.Int())
	case reflect.Int8:
		if src.Type().Kind() != reflect.Int8 {
			return MismatchTypeError
		}
		dst.SetInt(src.Int())
	case reflect.Int16:
		if src.Type().Kind() != reflect.Int16 {
			return MismatchTypeError
		}
		dst.SetInt(src.Int())
	case reflect.Int32:
		if src.Type().Kind() != reflect.Int32 {
			return MismatchTypeError
		}
		dst.SetInt(src.Int())
	case reflect.Int64:
		if src.Type().Kind() != reflect.Int64 {
			return MismatchTypeError
		}
		dst.SetInt(src.Int())
	case reflect.Uint:
		if src.Type().Kind() != reflect.Uint {
			return MismatchTypeError
		}
		dst.SetUint(src.Uint())
	case reflect.Uint8:
		if src.Type().Kind() != reflect.Uint8 {
			return MismatchTypeError
		}
		dst.SetUint(src.Uint())
	case reflect.Uint16:
		if src.Type().Kind() != reflect.Uint16 {
			return MismatchTypeError
		}
		dst.SetUint(src.Uint())
	case reflect.Uint32:
		if src.Type().Kind() != reflect.Uint32 {
			return MismatchTypeError
		}
		dst.SetUint(src.Uint())
	case reflect.Uint64:
		if src.Type().Kind() != reflect.Uint64 {
			return MismatchTypeError
		}
		dst.SetUint(src.Uint())
	case reflect.Uintptr:
		return nil // ignore
	case reflect.Float32:
		if src.Type().Kind() != reflect.Float32 {
			return MismatchTypeError
		}
		dst.SetFloat(src.Float())
	case reflect.Float64:
		if src.Type().Kind() != reflect.Float64 {
			return MismatchTypeError
		}
		dst.SetFloat(src.Float())
	case reflect.Complex64:
		if src.Type().Kind() != reflect.Complex64 {
			return MismatchTypeError
		}
		dst.SetComplex(src.Complex())
	case reflect.Complex128:
		if src.Type().Kind() != reflect.Complex128 {
			return MismatchTypeError
		}
		dst.SetComplex(src.Complex())
	case reflect.Array:
		return nil // TODO
	case reflect.Chan:
		return nil // ignore
	case reflect.Func:
		return nil // ignore
	case reflect.Interface:
		return nil // TODO: get underlying type
	case reflect.Map:
		return nil // TODO
	case reflect.Pointer:
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		return m.mapValue(src, dst.Elem())
	case reflect.Slice:
		return nil // TODO
	case reflect.String:
		if src.Type().Kind() != reflect.String {
			return MismatchTypeError
		}
		dst.SetString(src.String())
	case reflect.Struct:
		if src.Type().Kind() != reflect.Struct {
			return MismatchTypeError
		}
		for i := 0; i < dst.NumField(); i++ {
			dstField := dst.Field(i)
			srcField := src.FieldByName(dst.Type().Field(i).Name)
			err := m.mapValue(srcField, dstField)
			if err != nil {
				return err
			}
		}
	case reflect.UnsafePointer:
		return nil // ignore

	}
	return nil // ignore
}
