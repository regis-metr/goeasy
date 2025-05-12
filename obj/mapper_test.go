package obj

import (
	"fmt"
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
	src := struct {
		AllTypes testAllTypes
	}{
		AllTypes: testAllTypes{
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
		},
	}

	mapper := NewMapper()
	dst := struct {
		AllTypes testAllTypes
	}{
		AllTypes: testAllTypes{},
	}
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

func TestMapSliceDst(t *testing.T) {
	type IntStruct struct {
		Int int
	}
	tests := []struct {
		name string
		src  interface{}
		dst  []IntStruct
		err  error
	}{
		{
			name: "Array source",
			src:  [3]IntStruct{{1}, {2}, {3}},
		},
		{
			name: "Wrong source type",
			src:  1,
			err:  ErrMismatchType,
		},
		{
			name: "Slice source",
			src:  [3]IntStruct{{1}, {2}},
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

func TestMapMapDst(t *testing.T) {
	type IntStruct struct {
		Int int
	}
	tests := []struct {
		name string
		src  interface{}
		dst  map[int]IntStruct
		err  error
	}{
		{
			name: "Success",
			src:  map[int]IntStruct{1: {1}, 2: {2}},
		},
		{
			name: "Wrong source type",
			src:  1,
			err:  ErrMismatchType,
		},
		{
			name: "Wrong key type",
			src:  map[string]IntStruct{"1": {1}},
			err:  ErrMismatchType,
		},
		{
			name: "Wrong value type",
			src:  map[int]int{1: 1},
			err:  ErrMismatchType,
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
			assert.Equal(t, test.src, test.dst)
		})
	}
}

func TestMapWithFieldMaps(t *testing.T) {
	type IntStruct struct {
		IntField int
	}
	type WrappedAllTypes struct {
		AllTypes testAllTypes
	}
	tests := []struct {
		name      string
		src       any
		dst       any
		cfg       []FieldMapConfig
		expected  any
		configErr error
		mapErr    error
	}{
		{
			name: "Map source to dest field name",
			src:  testAllTypes{Int: 1},
			dst:  &IntStruct{},
			cfg: []FieldMapConfig{{
				Source:      "Int",
				Destination: "IntField",
			}},
			expected: &IntStruct{1},
		},
		{
			name: "Map source to dest field name with func",
			src:  testAllTypes{Int: 1},
			dst:  &IntStruct{},
			cfg: []FieldMapConfig{{
				Source:              "Int",
				Destination:         "IntField",
				GetDestinationValue: func(source any) (any, error) { return 2, nil },
			}},
			expected: &IntStruct{2},
		},
		{
			name: "Map source to dest field name with func error",
			src:  testAllTypes{Int: 1},
			dst:  &IntStruct{},
			cfg: []FieldMapConfig{{
				Source:              "Int",
				Destination:         "IntField",
				GetDestinationValue: func(source any) (any, error) { return nil, fmt.Errorf("Test Error") },
			}},
			mapErr: fmt.Errorf("Test Error"),
		},
		{
			name: "With func",
			src:  testAllTypes{Int: 1},
			dst:  &testAllTypes{},
			cfg: []FieldMapConfig{{
				Destination:         "Int",
				GetDestinationValue: func(source any) (any, error) { return 2, nil },
			}},
			expected: &testAllTypes{Int: 2},
		},
		{
			name: "Struct within struct",
			src:  WrappedAllTypes{testAllTypes{Int: 1}},
			dst:  &WrappedAllTypes{testAllTypes{}},
			cfg: []FieldMapConfig{{
				Destination:         "Int",
				GetDestinationValue: func(source any) (any, error) { return 2, nil },
			}},
			expected: &WrappedAllTypes{testAllTypes{Int: 2}},
		},
	}
	// TODO: add error cases

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			var err error
			if reflect.TypeOf(test.dst).Elem() == reflect.TypeOf(testAllTypes{}) {
				err = ConfigureFieldMaps[testAllTypes, testAllTypes](mapper, test.cfg...)
			} else if reflect.TypeOf(test.dst).Elem() == reflect.TypeOf(WrappedAllTypes{}) {
				err = ConfigureFieldMaps[testAllTypes, testAllTypes](mapper, test.cfg...)
			} else {
				err = ConfigureFieldMaps[testAllTypes, IntStruct](mapper, test.cfg...)
			}
			if test.configErr != nil || err != nil {
				assert.Equal(t, test.configErr, err)
				return
			}

			err = mapper.Map(test.src, test.dst)
			if test.mapErr != nil || err != nil {
				assert.Equal(t, test.mapErr, err)
				return
			}
			assert.Equal(t, test.expected, test.dst)
		})
	}
}

func TestConfigureFieldMapsSourceNotStruct(t *testing.T) {
	mapper := NewMapper()
	cfg := FieldMapConfig{
		Source:      "Int",
		Destination: "Int",
	}
	err := ConfigureFieldMaps[int, testAllTypes](mapper, cfg)
	assert.Equal(t, fmt.Errorf("sourceT and destinationT must be structs"), err)
}

func TestConfigureFieldMapsDestinationNotStruct(t *testing.T) {
	mapper := NewMapper()
	cfg := FieldMapConfig{
		Source:      "Int",
		Destination: "Int",
	}
	err := ConfigureFieldMaps[testAllTypes, int](mapper, cfg)
	assert.Equal(t, fmt.Errorf("sourceT and destinationT must be structs"), err)
}

func TestConfigureFieldMapsDestinationFieldEmpty(t *testing.T) {
	mapper := NewMapper()
	cfg := FieldMapConfig{
		Source: "Int",
	}
	err := ConfigureFieldMaps[testAllTypes, testAllTypes](mapper, cfg)
	assert.Equal(t, fmt.Errorf("destination field names must be provided"), err)
}

// AI generated code start
type testUserDTO struct {
	ID             int
	withGetterName string
}

type testUser struct {
	ID   int
	Name string
}

func (u testUserDTO) GetName() string {
	return "Mr. " + u.withGetterName
}

func TestMapWithGetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUser{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.ID, user.ID, "ID not equal")
	assert.Equal(t, "Mr. John", user.Name, "Name not equal")
}

type testUserWithSetter struct {
	withSetterID   int
	withSetterName string
}

func (u *testUserWithSetter) SetID(id int) {
	u.withSetterID = id
}

func (u *testUserWithSetter) SetName(name string) {
	u.withSetterName = name
}

func TestMapWithSetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithSetter{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.ID, user.withSetterID, "ID not equal")
	assert.Equal(t, dto.GetName(), user.withSetterName, "Name not equal")
}

type testUserWithDifferentSetter struct {
	withSetterID   int
	withSetterName string
}

func (u *testUserWithDifferentSetter) SetID(id int64) {
	u.withSetterID = int(id)
}

func (u *testUserWithDifferentSetter) SetName(name string) {
	u.withSetterName = name
}

func TestMapWithDifferentSetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithDifferentSetter{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Equal(t, ErrMismatchType, err, "Error not equal")
	assert.Equal(t, 0, user.withSetterID, "ID not equal")
	assert.Equal(t, "", user.withSetterName, "Name not equal")
}

type testUserWithGetterAndSetter struct {
	withGetterSetterID   int
	withGetterSetterName string
}

func (u *testUserWithGetterAndSetter) SetID(id int) {
	u.withGetterSetterID = id
}

func (u *testUserWithGetterAndSetter) SetName(name string) {
	u.withGetterSetterName = name
}

func (u testUserDTO) GetID() int {
	return u.ID
}

func TestMapWithGetterAndSetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithGetterAndSetter{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.GetID(), user.withGetterSetterID, "ID not equal")
	assert.Equal(t, dto.GetName(), user.withGetterSetterName, "Name not equal")
}

type testUserWithSetterNoField struct {
	withSetterID int
}

func (u *testUserWithSetterNoField) SetID(id int) {
	u.withSetterID = id
}

func TestMapWithSetterNoField(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithSetterNoField{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.ID, user.withSetterID, "ID not equal")
}

func TestMapWithSetterAndFieldMap(t *testing.T) {

	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithSetter{}
	mapper := NewMapper()
	err := ConfigureFieldMaps[testUserDTO, testUserWithSetter](mapper, FieldMapConfig{
		Source:      "Name",
		Destination: "Name",
		GetDestinationValue: func(source any) (any, error) {
			return "Mr. " + source.(string), nil
		},
	})

	assert.Nil(t, err, "ConfigureFieldMaps returned an error")

	err = mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.ID, user.withSetterID, "ID not equal")
	assert.Equal(t, "Mr. Mr. John", user.withSetterName, "Name not equal")
}

// AI generated code start
func TestMapWithInterfaceField(t *testing.T) {
	type StructWithInterface struct {
		Data interface{}
	}

	tests := []struct {
		name     string
		src      StructWithInterface
		dst      StructWithInterface
		expected interface{}
		err      error
	}{
		{
			name:     "String interface",
			src:      StructWithInterface{Data: "test string"},
			dst:      StructWithInterface{},
			expected: "test string",
		},
		{
			name:     "Int interface",
			src:      StructWithInterface{Data: 42},
			dst:      StructWithInterface{},
			expected: 42,
		},
		{
			name:     "Struct interface",
			src:      StructWithInterface{Data: testAllTypes{Int: 1, String: "test"}},
			dst:      StructWithInterface{},
			expected: testAllTypes{Int: 1, String: "test"},
		},
		{
			name:     "Nil interface",
			src:      StructWithInterface{Data: nil},
			dst:      StructWithInterface{},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			assert.Equal(t, test.err, err, "Error not equal")
			assert.Equal(t, test.expected, test.dst.Data, "Data not equal")
		})
	}
}

// AI generated code end

// AI generated code start
func TestMapWithNestedInterfaceField(t *testing.T) {
	type InnerStruct struct {
		Data interface{}
	}

	type OuterStruct struct {
		Inner interface{}
	}

	tests := []struct {
		name     string
		src      OuterStruct
		dst      OuterStruct
		expected interface{}
		err      error
	}{
		{
			name: "Nested string interface",
			src:  OuterStruct{Inner: InnerStruct{Data: "test string"}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: "test string",
			},
		},
		{
			name: "Nested int interface",
			src:  OuterStruct{Inner: InnerStruct{Data: 42}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: 42,
			},
		},
		{
			name: "Nested struct interface",
			src:  OuterStruct{Inner: InnerStruct{Data: testAllTypes{Int: 1, String: "test"}}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: testAllTypes{Int: 1, String: "test"},
			},
		},
		{
			name: "Nested nil interface",
			src:  OuterStruct{Inner: InnerStruct{Data: nil}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			assert.Equal(t, test.err, err, "Error not equal")
			assert.Equal(t, test.expected, test.dst.Inner, "Inner not equal")
		})
	}
}

// AI generated code end

// AI generated code start
func TestMapNotAddressable(t *testing.T) {
	type TestStruct struct {
		Field int
	}

	src := TestStruct{Field: 42}
	dst := TestStruct{}

	mapper := NewMapper()
	err := mapper.Map(src, dst) // Passing non-addressable value as destination

	assert.Equal(t, ErrNotAddresable, err, "Error not equal")
}

// AI generated code end

// AI generated code start
func TestMapWithInitializedInterfaceField(t *testing.T) {
	type StructWithInterface struct {
		Data interface{}
	}

	src := StructWithInterface{Data: "test string"}
	dst := StructWithInterface{Data: ""} // Initialized with the same type as the source

	mapper := NewMapper()
	err := mapper.Map(src, &dst)

	assert.Equal(t, ErrNotAddresable, err, "Error not equal")
	assert.NotEqual(t, src.Data, dst.Data, "Data equal")
}

// AI generated code end
// AI generated code start
func TestMapWithNestedStructAndCustomFunction(t *testing.T) {
	type InnerStruct struct {
		Value int
	}

	type OuterStruct struct {
		Inner interface{}
	}

	src := OuterStruct{
		Inner: InnerStruct{Value: 42},
	}

	dst := OuterStruct{}

	mapper := NewMapper()
	err := ConfigureFieldMaps[OuterStruct, OuterStruct](mapper, FieldMapConfig{
		Source:      "Inner",
		Destination: "Inner",
		GetDestinationValue: func(source any) (any, error) {
			if inner, ok := source.(InnerStruct); ok {
				return InnerStruct{Value: inner.Value * 2}, nil
			}
			fmt.Println(reflect.TypeOf(source))
			return nil, fmt.Errorf("unexpected type")
		},
	})

	assert.Nil(t, err, "ConfigureFieldMaps returned an error")

	err = mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, InnerStruct{Value: 84}, dst.Inner, "Inner not equal")
}

// AI generated code end
