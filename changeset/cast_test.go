package changeset

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCast(t *testing.T) {
	var entity struct {
		Field1 int `db:"field1"`
		Field2 string
		Field3 bool
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": true,
		"field4": "ignore please",
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": true,
	}

	expectedSchema := map[string]Field{
		"field1": Field{0, reflect.TypeOf(0)},
		"field2": Field{"", reflect.TypeOf("")},
		"field3": Field{false, reflect.TypeOf(false)},
	}

	ch := Cast(entity, params, []string{"field1", "field2", "field3"})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedSchema, ch.schema)

	ch = Cast(&entity, params, []string{"field1", "field2", "field3"})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedSchema, ch.schema)
}

func TestCastError(t *testing.T) {
	var entity struct {
		Field1 int
	}

	params := map[string]interface{}{
		"field1": "1",
	}

	ch := Cast(entity, params, []string{"field1"})
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field1 is invalid", ch.Error().Error())
}

func TestCastPanic(t *testing.T) {
	params := map[string]interface{}{
		"field1": "1",
	}

	assert.Panics(t, func() {
		Cast("entity", params, []string{"field1"})
	})
}

func TestCastBasic(t *testing.T) {
	var entity struct {
		Field1  bool
		Field2  int
		Field3  int8
		Field4  int16
		Field5  int32
		Field6  int64
		Field7  uint
		Field8  uint8
		Field9  uint16
		Field10 uint32
		Field11 uint64
		Field12 uintptr
		Field13 float32
		Field14 float64
		Field15 string
	}

	params := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedChanges := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedSchema := map[string]Field{
		"field1":  Field{false, reflect.TypeOf(false)},
		"field2":  Field{0, reflect.TypeOf(0)},
		"field3":  Field{int8(0), reflect.TypeOf(int8(0))},
		"field4":  Field{int16(0), reflect.TypeOf(int16(0))},
		"field5":  Field{int32(0), reflect.TypeOf(int32(0))},
		"field6":  Field{int64(0), reflect.TypeOf(int64(0))},
		"field7":  Field{uint(0), reflect.TypeOf(uint(0))},
		"field8":  Field{uint8(0), reflect.TypeOf(uint8(0))},
		"field9":  Field{uint16(0), reflect.TypeOf(uint16(0))},
		"field10": Field{uint32(0), reflect.TypeOf(uint32(0))},
		"field11": Field{uint64(0), reflect.TypeOf(uint64(0))},
		"field12": Field{uintptr(0), reflect.TypeOf(uintptr(0))},
		"field13": Field{float32(0), reflect.TypeOf(float32(0))},
		"field14": Field{float64(0), reflect.TypeOf(float64(0))},
		"field15": Field{"", reflect.TypeOf("")},
	}

	ch := Cast(entity, params, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedSchema, ch.schema)
}

func TestCastBasicWithValue(t *testing.T) {
	entity := struct {
		Field1  bool
		Field2  int
		Field3  int8
		Field4  int16
		Field5  int32
		Field6  int64
		Field7  uint
		Field8  uint8
		Field9  uint16
		Field10 uint32
		Field11 uint64
		Field12 uintptr
		Field13 float32
		Field14 float64
		Field15 string
	}{
		true,
		1,
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		uintptr(1),
		float32(1),
		float64(1),
		"1",
	}

	params := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedChanges := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedSchema := map[string]Field{
		"field1":  Field{true, reflect.TypeOf(false)},
		"field2":  Field{1, reflect.TypeOf(0)},
		"field3":  Field{int8(1), reflect.TypeOf(int8(0))},
		"field4":  Field{int16(1), reflect.TypeOf(int16(0))},
		"field5":  Field{int32(1), reflect.TypeOf(int32(0))},
		"field6":  Field{int64(1), reflect.TypeOf(int64(0))},
		"field7":  Field{uint(1), reflect.TypeOf(uint(0))},
		"field8":  Field{uint8(1), reflect.TypeOf(uint8(0))},
		"field9":  Field{uint16(1), reflect.TypeOf(uint16(0))},
		"field10": Field{uint32(1), reflect.TypeOf(uint32(0))},
		"field11": Field{uint64(1), reflect.TypeOf(uint64(0))},
		"field12": Field{uintptr(1), reflect.TypeOf(uintptr(0))},
		"field13": Field{float32(1), reflect.TypeOf(float32(0))},
		"field14": Field{float64(1), reflect.TypeOf(float64(0))},
		"field15": Field{"1", reflect.TypeOf("")},
	}

	ch := Cast(entity, params, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedSchema, ch.schema)
}

func TestCastPtr(t *testing.T) {
	var entity struct {
		Field1  *bool
		Field2  *int
		Field3  *int8
		Field4  *int16
		Field5  *int32
		Field6  *int64
		Field7  *uint
		Field8  *uint8
		Field9  *uint16
		Field10 *uint32
		Field11 *uint64
		Field12 *uintptr
		Field13 *float32
		Field14 *float64
		Field15 *string
	}

	params := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedChanges := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedSchema := map[string]Field{
		"field1":  Field{nil, reflect.TypeOf(false)},
		"field2":  Field{nil, reflect.TypeOf(0)},
		"field3":  Field{nil, reflect.TypeOf(int8(0))},
		"field4":  Field{nil, reflect.TypeOf(int16(0))},
		"field5":  Field{nil, reflect.TypeOf(int32(0))},
		"field6":  Field{nil, reflect.TypeOf(int64(0))},
		"field7":  Field{nil, reflect.TypeOf(uint(0))},
		"field8":  Field{nil, reflect.TypeOf(uint8(0))},
		"field9":  Field{nil, reflect.TypeOf(uint16(0))},
		"field10": Field{nil, reflect.TypeOf(uint32(0))},
		"field11": Field{nil, reflect.TypeOf(uint64(0))},
		"field12": Field{nil, reflect.TypeOf(uintptr(0))},
		"field13": Field{nil, reflect.TypeOf(float32(0))},
		"field14": Field{nil, reflect.TypeOf(float64(0))},
		"field15": Field{nil, reflect.TypeOf("")},
	}

	ch := Cast(entity, params, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedSchema, ch.schema)
}

func TestCastPtrWithValue(t *testing.T) {
	var (
		vbool    = true
		vint     = int(1)
		vint8    = int8(1)
		vint16   = int16(1)
		vint32   = int32(1)
		vint64   = int64(1)
		vuint    = uint(1)
		vuint8   = uint8(1)
		vuint16  = uint16(1)
		vuint32  = uint32(1)
		vuint64  = uint64(1)
		vuintptr = uintptr(1)
		vfloat32 = float32(1)
		vfloat64 = float64(1)
		vstring  = ""
	)
	entity := struct {
		Field1  *bool
		Field2  *int
		Field3  *int8
		Field4  *int16
		Field5  *int32
		Field6  *int64
		Field7  *uint
		Field8  *uint8
		Field9  *uint16
		Field10 *uint32
		Field11 *uint64
		Field12 *uintptr
		Field13 *float32
		Field14 *float64
		Field15 *string
	}{
		&vbool,
		&vint,
		&vint8,
		&vint16,
		&vint32,
		&vint64,
		&vuint,
		&vuint8,
		&vuint16,
		&vuint32,
		&vuint64,
		&vuintptr,
		&vfloat32,
		&vfloat64,
		&vstring,
	}

	params := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedChanges := map[string]interface{}{
		"field1":  true,
		"field2":  2,
		"field3":  3,
		"field4":  4,
		"field5":  5,
		"field6":  6,
		"field7":  7,
		"field8":  8,
		"field9":  9,
		"field10": 10,
		"field11": 11,
		"field12": 12,
		"field13": 13,
		"field14": 14,
		"field15": "15",
	}

	expectedSchema := map[string]Field{
		"field1":  Field{vbool, reflect.TypeOf(false)},
		"field2":  Field{vint, reflect.TypeOf(0)},
		"field3":  Field{vint8, reflect.TypeOf(int8(0))},
		"field4":  Field{vint16, reflect.TypeOf(int16(0))},
		"field5":  Field{vint32, reflect.TypeOf(int32(0))},
		"field6":  Field{vint64, reflect.TypeOf(int64(0))},
		"field7":  Field{vuint, reflect.TypeOf(uint(0))},
		"field8":  Field{vuint8, reflect.TypeOf(uint8(0))},
		"field9":  Field{vuint16, reflect.TypeOf(uint16(0))},
		"field10": Field{vuint32, reflect.TypeOf(uint32(0))},
		"field11": Field{vuint64, reflect.TypeOf(uint64(0))},
		"field12": Field{vuintptr, reflect.TypeOf(uintptr(0))},
		"field13": Field{vfloat32, reflect.TypeOf(float32(0))},
		"field14": Field{vfloat64, reflect.TypeOf(float64(0))},
		"field15": Field{vstring, reflect.TypeOf("")},
	}

	ch := Cast(entity, params, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedSchema, ch.schema)
}

func TestCastSlice(t *testing.T) {
	var entity struct {
		Field1  []bool
		Field2  []int
		Field3  []int8
		Field4  []int16
		Field5  []int32
		Field6  []int64
		Field7  []uint
		Field8  []uint8
		Field9  []uint16
		Field10 []uint32
		Field11 []uint64
		Field12 []uintptr
		Field13 []float32
		Field14 []float64
		Field15 []string
	}

	params := map[string]interface{}{
		"field1":  []bool{true},
		"field2":  []int{2},
		"field3":  []int8{3},
		"field4":  []int16{4},
		"field5":  []int32{5},
		"field6":  []int64{6},
		"field7":  []uint{7},
		"field8":  []uint8{8},
		"field9":  []uint16{9},
		"field10": []uint32{10},
		"field11": []uint64{11},
		"field12": []uintptr{12},
		"field13": []float32{13},
		"field14": []float64{14},
		"field15": []string{"15"},
	}

	expectedChanges := map[string]interface{}{
		"field1":  []bool{true},
		"field2":  []int{2},
		"field3":  []int8{3},
		"field4":  []int16{4},
		"field5":  []int32{5},
		"field6":  []int64{6},
		"field7":  []uint{7},
		"field8":  []uint8{8},
		"field9":  []uint16{9},
		"field10": []uint32{10},
		"field11": []uint64{11},
		"field12": []uintptr{12},
		"field13": []float32{13},
		"field14": []float64{14},
		"field15": []string{"15"},
	}

	expectedSchema := map[string]Field{
		"field1":  Field{[]bool(nil), reflect.TypeOf([]bool{})},
		"field2":  Field{[]int(nil), reflect.TypeOf([]int{})},
		"field3":  Field{[]int8(nil), reflect.TypeOf([]int8{})},
		"field4":  Field{[]int16(nil), reflect.TypeOf([]int16{})},
		"field5":  Field{[]int32(nil), reflect.TypeOf([]int32{})},
		"field6":  Field{[]int64(nil), reflect.TypeOf([]int64{})},
		"field7":  Field{[]uint(nil), reflect.TypeOf([]uint{})},
		"field8":  Field{[]uint8(nil), reflect.TypeOf([]uint8{})},
		"field9":  Field{[]uint16(nil), reflect.TypeOf([]uint16{})},
		"field10": Field{[]uint32(nil), reflect.TypeOf([]uint32{})},
		"field11": Field{[]uint64(nil), reflect.TypeOf([]uint64{})},
		"field12": Field{[]uintptr(nil), reflect.TypeOf([]uintptr{})},
		"field13": Field{[]float32(nil), reflect.TypeOf([]float32{})},
		"field14": Field{[]float64(nil), reflect.TypeOf([]float64{})},
		"field15": Field{[]string(nil), reflect.TypeOf([]string{})},
	}

	ch := Cast(entity, params, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, ch.Changes(), expectedChanges)
	assert.Equal(t, expectedSchema, ch.schema)
}

func TestCastSliceWithValue(t *testing.T) {
	entity := struct {
		Field1  []bool
		Field2  []int
		Field3  []int8
		Field4  []int16
		Field5  []int32
		Field6  []int64
		Field7  []uint
		Field8  []uint8
		Field9  []uint16
		Field10 []uint32
		Field11 []uint64
		Field12 []uintptr
		Field13 []float32
		Field14 []float64
		Field15 []string
	}{
		[]bool{true},
		[]int{1},
		[]int8{1},
		[]int16{1},
		[]int32{1},
		[]int64{1},
		[]uint{1},
		[]uint8{1},
		[]uint16{1},
		[]uint32{1},
		[]uint64{1},
		[]uintptr{1},
		[]float32{1},
		[]float64{1},
		[]string{"1"},
	}

	params := map[string]interface{}{
		"field1":  []bool{true},
		"field2":  []int{2},
		"field3":  []int8{3},
		"field4":  []int16{4},
		"field5":  []int32{5},
		"field6":  []int64{6},
		"field7":  []uint{7},
		"field8":  []uint8{8},
		"field9":  []uint16{9},
		"field10": []uint32{10},
		"field11": []uint64{11},
		"field12": []uintptr{12},
		"field13": []float32{13},
		"field14": []float64{14},
		"field15": []string{"15"},
	}

	expectedChanges := map[string]interface{}{
		"field1":  []bool{true},
		"field2":  []int{2},
		"field3":  []int8{3},
		"field4":  []int16{4},
		"field5":  []int32{5},
		"field6":  []int64{6},
		"field7":  []uint{7},
		"field8":  []uint8{8},
		"field9":  []uint16{9},
		"field10": []uint32{10},
		"field11": []uint64{11},
		"field12": []uintptr{12},
		"field13": []float32{13},
		"field14": []float64{14},
		"field15": []string{"15"},
	}

	expectedSchema := map[string]Field{
		"field1":  Field{[]bool{true}, reflect.TypeOf([]bool{})},
		"field2":  Field{[]int{1}, reflect.TypeOf([]int{})},
		"field3":  Field{[]int8{1}, reflect.TypeOf([]int8{})},
		"field4":  Field{[]int16{1}, reflect.TypeOf([]int16{})},
		"field5":  Field{[]int32{1}, reflect.TypeOf([]int32{})},
		"field6":  Field{[]int64{1}, reflect.TypeOf([]int64{})},
		"field7":  Field{[]uint{1}, reflect.TypeOf([]uint{})},
		"field8":  Field{[]uint8{1}, reflect.TypeOf([]uint8{})},
		"field9":  Field{[]uint16{1}, reflect.TypeOf([]uint16{})},
		"field10": Field{[]uint32{1}, reflect.TypeOf([]uint32{})},
		"field11": Field{[]uint64{1}, reflect.TypeOf([]uint64{})},
		"field12": Field{[]uintptr{1}, reflect.TypeOf([]uintptr{})},
		"field13": Field{[]float32{1}, reflect.TypeOf([]float32{})},
		"field14": Field{[]float64{1}, reflect.TypeOf([]float64{})},
		"field15": Field{[]string{"1"}, reflect.TypeOf([]string{})},
	}

	ch := Cast(entity, params, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, ch.Changes(), expectedChanges)
	assert.Equal(t, expectedSchema, ch.schema)
}
