package obj

import (
	"math"
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

func TestMapNotAllowedTypes(t *testing.T) {
	tests := []struct {
		name string
		src  interface{}
		dst  interface{}
		exp  interface{}
	}{}
	mapper := NewMapper()
	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			err := mapper.Map(scenario.src, scenario.dst)
			assert.Equal(t, ErrMismatchType, err, "Error is not mismatch")
			assert.Equal(t, scenario.exp, scenario.dst, "Expected and destination is not equal")
			assert.NotEqual(t, scenario.dst, scenario.src, "Source is copied to destination")
		})
	}
}
