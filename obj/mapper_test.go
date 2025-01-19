package obj

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testAllTypes struct {
	Bool       bool
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	Pointer    *int
	String     string
}

func TestMap(t *testing.T) {
	i := 1
	src := testAllTypes{
		Bool:       true,
		Int:        math.MaxInt,
		Int8:       math.MaxInt8,
		Int16:      math.MaxInt16,
		Int32:      math.MaxInt32,
		Int64:      math.MaxInt64,
		Uint:       math.MaxUint,
		Uint8:      math.MaxUint8,
		Uint16:     math.MaxUint16,
		Uint32:     math.MaxUint32,
		Uint64:     math.MaxUint64,
		Float32:    math.MaxFloat32,
		Float64:    math.MaxFloat64,
		Complex64:  complex(math.MaxFloat32, math.MaxFloat32),
		Complex128: complex(math.MaxFloat64, math.MaxFloat64),
		Pointer:    &i,
		String:     "test",
	}

	mapper := NewMapper()
	dst := testAllTypes{}
	err := mapper.Map(src, &dst)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src, dst, "Not equal")
}

func TestMapNoEquivalentField(t *testing.T) {
	mapper := NewMapper()
	i := 1
	src := testAllTypes{
		Bool:       true,
		Int:        math.MaxInt,
		Int8:       math.MaxInt8,
		Int16:      math.MaxInt16,
		Int32:      math.MaxInt32,
		Int64:      math.MaxInt64,
		Uint:       math.MaxUint,
		Uint8:      math.MaxUint8,
		Uint16:     math.MaxUint16,
		Uint32:     math.MaxUint32,
		Uint64:     math.MaxUint64,
		Float32:    math.MaxFloat32,
		Float64:    math.MaxFloat64,
		Complex64:  complex(math.MaxFloat32, math.MaxFloat32),
		Complex128: complex(math.MaxFloat64, math.MaxFloat64),
		Pointer:    &i,
		String:     "test",
	}
	dst := struct {
		Int  int
		Test string
	}{}

	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src.Int, dst.Int, "Int are not equal")
	assert.Empty(t, dst.Test)
}

func TestMapArrayDst(t *testing.T) {
	type IntStruct struct {
		Int int
	}
	tests := []struct {
		name string
		src  interface{}
		dst  [3]IntStruct
		err  error
	}{
		{
			name: "Equal length array",
			src:  [3]IntStruct{{1}, {2}, {3}},
		},
		{
			name: "Longer array",
			src:  [4]IntStruct{{1}, {2}, {3}, {4}},
			err:  ErrInsufficientCapacity,
		},
		{
			name: "Shorter array",
			src:  [3]IntStruct{{1}, {2}},
		},
		{
			name: "Equal length slice",
			src:  []IntStruct{{1}, {2}, {3}},
		},
		{
			name: "Longer slice",
			src:  []IntStruct{{1}, {2}, {3}, {4}},
			err:  ErrInsufficientCapacity,
		},
		{
			name: "Shorter slice",
			src:  make([]IntStruct, 2, 5),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			if err != nil || test.err != nil {
				assert.Equal(t, test.err, err, "Error not equal")
				return
			}
			srcV := reflect.ValueOf(test.src)
			for i := range srcV.Len() {
				v := srcV.Index(i).Field(0)
				assert.Equal(t, int(v.Int()), test.dst[i].Int, "Has unequal value")
			}
		})
	}
}

func TestMapMismatchType(t *testing.T) {
	i := 1
	src := struct {
		Int    int
		String string
		Ptr    *int
	}{Ptr: &i}
	tests := []struct {
		name string
		dst  interface{}
	}{
		{
			name: "Bool",
			dst: &struct {
				Int bool
			}{},
		},
		{
			name: "Int",
			dst: &struct {
				String int
			}{},
		},
		{
			name: "Int8",
			dst: &struct {
				Int int8
			}{},
		},
		{
			name: "Int16",
			dst: &struct {
				Int int8
			}{},
		},
		{
			name: "Int32",
			dst: &struct {
				Int int8
			}{},
		},
		{
			name: "Int64",
			dst: &struct {
				Int int64
			}{},
		},
		{
			name: "Uint",
			dst: &struct {
				Int uint
			}{},
		},
		{
			name: "Uint8",
			dst: &struct {
				Int uint8
			}{},
		},
		{
			name: "Uint16",
			dst: &struct {
				Int uint16
			}{},
		},
		{
			name: "Uint32",
			dst: &struct {
				Int uint32
			}{},
		},
		{
			name: "Uint64",
			dst: &struct {
				Int uint64
			}{},
		},
		{
			name: "Float32",
			dst: &struct {
				Int float32
			}{},
		},
		{
			name: "Float64",
			dst: &struct {
				Int float64
			}{},
		},
		{
			name: "Complex64",
			dst: &struct {
				Int complex64
			}{},
		},
		{
			name: "Complex128",
			dst: &struct {
				Int complex128
			}{},
		},
		{
			name: "Pointer",
			dst: &struct {
				Ptr *int8
			}{},
		},
		{
			name: "String",
			dst: &struct {
				Int string
			}{},
		},
	}

	for _, test := range tests {
		//t.Run(test.name, func(t *testing.T) {
		mapper := NewMapper()
		err := mapper.Map(src, test.dst)
		assert.Equal(t, ErrMismatchType, err)
		//})
	}
}
