package obj

import (
	"fmt"
	"reflect"
)

// ErrMismatchType returned when field of source can't be mapped to destination due to mismatched types.
var ErrMismatchType error = fmt.Errorf("type mismatch")
var ErrUnsupportedType error = fmt.Errorf("type unsupported")

type Mapper struct{}

// NewMapper creates a new instance of Mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// Map copies src field values to dst fields. Fields must have the same name.
// Sample usage:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/bryan-t/goeasy/obj"
//	)
//
//	type UserDTO struct {
//		ID   int
//		Name string
//	}
//
//	type User struct {
//		ID   int
//		Name string
//	}
//
//	func main() {
//		dto := UserDTO{
//			ID:   1,
//			Name: "John",
//		}
//
//		user := User{}
//		mapper := obj.NewMapper()
//		err := mapper.Map(dto, &user)
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Printf("user: %+v\n", user)
//	}
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
			return ErrMismatchType
		}
		dst.SetBool(src.Bool())
	case reflect.Int:
		if src.Type().Kind() != reflect.Int {
			return ErrMismatchType
		}
		dst.SetInt(src.Int())
	case reflect.Int8:
		if src.Type().Kind() != reflect.Int8 {
			return ErrMismatchType
		}
		dst.SetInt(src.Int())
	case reflect.Int16:
		if src.Type().Kind() != reflect.Int16 {
			return ErrMismatchType
		}
		dst.SetInt(src.Int())
	case reflect.Int32:
		if src.Type().Kind() != reflect.Int32 {
			return ErrMismatchType
		}
		dst.SetInt(src.Int())
	case reflect.Int64:
		if src.Type().Kind() != reflect.Int64 {
			return ErrMismatchType
		}
		dst.SetInt(src.Int())
	case reflect.Uint:
		if src.Type().Kind() != reflect.Uint {
			return ErrMismatchType
		}
		dst.SetUint(src.Uint())
	case reflect.Uint8:
		if src.Type().Kind() != reflect.Uint8 {
			return ErrMismatchType
		}
		dst.SetUint(src.Uint())
	case reflect.Uint16:
		if src.Type().Kind() != reflect.Uint16 {
			return ErrMismatchType
		}
		dst.SetUint(src.Uint())
	case reflect.Uint32:
		if src.Type().Kind() != reflect.Uint32 {
			return ErrMismatchType
		}
		dst.SetUint(src.Uint())
	case reflect.Uint64:
		if src.Type().Kind() != reflect.Uint64 {
			return ErrMismatchType
		}
		dst.SetUint(src.Uint())
	case reflect.Uintptr:
		return nil // ignore
	case reflect.Float32:
		if src.Type().Kind() != reflect.Float32 {
			return ErrMismatchType
		}
		dst.SetFloat(src.Float())
	case reflect.Float64:
		if src.Type().Kind() != reflect.Float64 {
			return ErrMismatchType
		}
		dst.SetFloat(src.Float())
	case reflect.Complex64:
		if src.Type().Kind() != reflect.Complex64 {
			return ErrMismatchType
		}
		dst.SetComplex(src.Complex())
	case reflect.Complex128:
		if src.Type().Kind() != reflect.Complex128 {
			return ErrMismatchType
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
			return ErrMismatchType
		}
		dst.SetString(src.String())
	case reflect.Struct:
		if src.Type().Kind() != reflect.Struct {
			return ErrMismatchType
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
