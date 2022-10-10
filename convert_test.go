// Modified from: database/sql/convert_test.go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style

package rel

import (
	"database/sql"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var someTime = time.Unix(123, 0)
var answer int64 = 42

type (
	userDefined       float64
	userDefinedSlice  []int
	userDefinedString string
)

type conversionTest struct {
	s, d any // source and destination

	// following are used if they're non-zero
	wantint    int64
	wantuint   uint64
	wantstr    string
	wantbytes  []byte
	wantraw    sql.RawBytes
	wantf32    float32
	wantf64    float64
	wanttime   time.Time
	wantbool   bool // used if d is of type *bool
	wanterr    string
	wantiface  any
	wantptr    *int64 // if non-nil, *d's pointed value must be equal to *wantptr
	wantnil    bool   // if true, *d must be *int64(nil)
	wantusrdef userDefined
	wantusrstr userDefinedString
}

// Target variables for scanning into.
var (
	scanstr    string
	scanbytes  []byte
	scanraw    sql.RawBytes
	scanint    int
	scanint8   int8
	scanint16  int16
	scanint32  int32
	scanuint8  uint8
	scanuint16 uint16
	scanbool   bool
	scanf32    float32
	scanf64    float64
	scantime   time.Time
	scanptr    *int64
	scaniface  any
)

func conversionTests() []conversionTest {
	// Return a fresh instance to test so "go test -count 2" works correctly.
	return []conversionTest{
		// Exact conversions (destination pointer type matches source type)
		{s: "foo", d: &scanstr, wantstr: "foo"},
		{s: 123, d: &scanint, wantint: 123},
		{s: someTime, d: &scantime, wanttime: someTime},
		{s: someTime.UTC(), d: &scantime, wanttime: someTime.UTC()},

		// To strings
		{s: "string", d: &scanstr, wantstr: "string"},
		{s: []byte("byteslice"), d: &scanstr, wantstr: "byteslice"},
		{s: true, d: &scanstr, wantstr: "true"},
		{s: 123, d: &scanstr, wantstr: "123"},
		{s: int8(123), d: &scanstr, wantstr: "123"},
		{s: int64(123), d: &scanstr, wantstr: "123"},
		{s: uint8(123), d: &scanstr, wantstr: "123"},
		{s: uint16(123), d: &scanstr, wantstr: "123"},
		{s: uint32(123), d: &scanstr, wantstr: "123"},
		{s: uint64(123), d: &scanstr, wantstr: "123"},
		{s: 1.5, d: &scanstr, wantstr: "1.5"},
		{s: float32(1.5), d: &scanstr, wantstr: "1.5"},

		// From time.Time:
		{s: time.Unix(1, 0).UTC(), d: &scanstr, wantstr: "1970-01-01T00:00:01Z"},
		{s: time.Unix(1453874597, 0).In(time.FixedZone("here", -3600*8)), d: &scanstr, wantstr: "2016-01-26T22:03:17-08:00"},
		{s: time.Unix(1, 2).UTC(), d: &scanstr, wantstr: "1970-01-01T00:00:01.000000002Z"},
		{s: time.Time{}, d: &scanstr, wantstr: "0001-01-01T00:00:00Z"},
		{s: time.Unix(1, 2).UTC(), d: &scanbytes, wantbytes: []byte("1970-01-01T00:00:01.000000002Z")},
		{s: time.Unix(1, 2).UTC(), d: &scaniface, wantiface: time.Unix(1, 2).UTC()},

		// To []byte
		{s: nil, d: &scanbytes, wantbytes: nil},
		{s: "string", d: &scanbytes, wantbytes: []byte("string")},
		{s: []byte("byteslice"), d: &scanbytes, wantbytes: []byte("byteslice")},
		{s: 123, d: &scanbytes, wantbytes: []byte("123")},
		{s: int8(123), d: &scanbytes, wantbytes: []byte("123")},
		{s: int64(123), d: &scanbytes, wantbytes: []byte("123")},
		{s: uint8(123), d: &scanbytes, wantbytes: []byte("123")},
		{s: uint16(123), d: &scanbytes, wantbytes: []byte("123")},
		{s: uint32(123), d: &scanbytes, wantbytes: []byte("123")},
		{s: uint64(123), d: &scanbytes, wantbytes: []byte("123")},
		{s: 1.5, d: &scanbytes, wantbytes: []byte("1.5")},
		{s: userDefinedString("user string"), d: &scanbytes, wantbytes: []byte("user string")},

		// To sql.RawBytes
		{s: nil, d: &scanraw, wantraw: nil},
		{s: []byte("byteslice"), d: &scanraw, wantraw: sql.RawBytes("byteslice")},
		{s: "string", d: &scanraw, wantraw: sql.RawBytes("string")},
		{s: 123, d: &scanraw, wantraw: sql.RawBytes("123")},
		{s: int8(123), d: &scanraw, wantraw: sql.RawBytes("123")},
		{s: int64(123), d: &scanraw, wantraw: sql.RawBytes("123")},
		{s: uint8(123), d: &scanraw, wantraw: sql.RawBytes("123")},
		{s: uint16(123), d: &scanraw, wantraw: sql.RawBytes("123")},
		{s: uint32(123), d: &scanraw, wantraw: sql.RawBytes("123")},
		{s: uint64(123), d: &scanraw, wantraw: sql.RawBytes("123")},
		{s: 1.5, d: &scanraw, wantraw: sql.RawBytes("1.5")},
		{s: userDefinedString("user string"), d: &scanraw, wantbytes: sql.RawBytes("user string")},
		// time.Time has been placed here to check that the sql.RawBytes slice gets
		// correctly reset when calling time.Time.AppendFormat.
		{s: time.Unix(2, 5).UTC(), d: &scanraw, wantraw: sql.RawBytes("1970-01-01T00:00:02.000000005Z")},

		// To integers
		{s: "255", d: &scanuint8, wantuint: 255},
		{s: []byte("255"), d: &scanuint8, wantuint: 255},
		{s: "256", d: &scanuint8, wanterr: "converting driver.Value type string (\"256\") to a uint8: value out of range"},
		{s: "256", d: &scanuint16, wantuint: 256},
		{s: "-1", d: &scanint, wantint: -1},
		{s: "foo", d: &scanint, wanterr: "converting driver.Value type string (\"foo\") to a int: invalid syntax"},

		// int64 to smaller integers
		{s: int64(5), d: &scanuint8, wantuint: 5},
		{s: int64(256), d: &scanuint8, wanterr: "converting driver.Value type int64 (\"256\") to a uint8: value out of range"},
		{s: int64(256), d: &scanuint16, wantuint: 256},
		{s: int64(65536), d: &scanuint16, wanterr: "converting driver.Value type int64 (\"65536\") to a uint16: value out of range"},

		// True bools
		{s: true, d: &scanbool, wantbool: true},
		{s: "True", d: &scanbool, wantbool: true},
		{s: "TRUE", d: &scanbool, wantbool: true},
		{s: "1", d: &scanbool, wantbool: true},
		{s: 1, d: &scanbool, wantbool: true},
		{s: int64(1), d: &scanbool, wantbool: true},
		{s: uint16(1), d: &scanbool, wantbool: true},

		// False bools
		{s: false, d: &scanbool, wantbool: false},
		{s: "false", d: &scanbool, wantbool: false},
		{s: "FALSE", d: &scanbool, wantbool: false},
		{s: "0", d: &scanbool, wantbool: false},
		{s: 0, d: &scanbool, wantbool: false},
		{s: int64(0), d: &scanbool, wantbool: false},
		{s: uint16(0), d: &scanbool, wantbool: false},

		// Not bools
		{s: "yup", d: &scanbool, wanterr: `sql/driver: couldn't convert "yup" into type bool`},
		{s: 2, d: &scanbool, wanterr: `sql/driver: couldn't convert 2 into type bool`},

		// Floats
		{s: float64(1.5), d: &scanf64, wantf64: float64(1.5)},
		{s: int64(1), d: &scanf64, wantf64: float64(1)},
		{s: float64(1.5), d: &scanf32, wantf32: float32(1.5)},
		{s: "1.5", d: &scanf32, wantf32: float32(1.5)},
		{s: "foo", d: &scanf32, wanterr: "converting driver.Value type string (\"foo\") to a float32: invalid syntax"},
		{s: "1.5", d: &scanf64, wantf64: float64(1.5)},
		{s: "foo", d: &scanf64, wanterr: "converting driver.Value type string (\"foo\") to a float64: invalid syntax"},

		// Pointers
		{s: any(nil), d: &scanptr, wantnil: true},
		{s: int64(42), d: &scanptr, wantptr: &answer},

		// To any
		{s: float64(1.5), d: &scaniface, wantiface: float64(1.5)},
		{s: int64(1), d: &scaniface, wantiface: int64(1)},
		{s: "str", d: &scaniface, wantiface: "str"},
		{s: []byte("byteslice"), d: &scaniface, wantiface: []byte("byteslice")},
		{s: true, d: &scaniface, wantiface: true},
		{s: nil, d: &scaniface},
		{s: []byte(nil), d: &scaniface, wantiface: []byte(nil)},

		// To a user-defined type
		{s: 1.5, d: new(userDefined), wantusrdef: 1.5},
		{s: int64(123), d: new(userDefined), wantusrdef: 123},
		{s: "1.5", d: new(userDefined), wantusrdef: 1.5},
		{s: []byte{1, 2, 3}, d: new(userDefinedSlice), wanterr: `unsupported Scan, storing driver.Value type []uint8 into type *rel.userDefinedSlice`},
		{s: "str", d: new(userDefinedString), wantusrstr: "str"},
		{s: []byte("byte"), d: new(userDefinedString), wantusrstr: "byte"},

		// Other errors
		{s: complex(1, 2), d: &scanstr, wanterr: `unsupported Scan, storing driver.Value type complex128 into type *string`},
		{s: complex(1, 2), d: &scanbytes, wanterr: `unsupported Scan, storing driver.Value type complex128 into type *[]uint8`},
	}
}

func intPtrValue(intptr any) any {
	return reflect.Indirect(reflect.Indirect(reflect.ValueOf(intptr))).Int()
}

func intValue(intptr any) int64 {
	return reflect.Indirect(reflect.ValueOf(intptr)).Int()
}

func uintValue(intptr any) uint64 {
	return reflect.Indirect(reflect.ValueOf(intptr)).Uint()
}

func float64Value(ptr any) float64 {
	return *(ptr.(*float64))
}

func float32Value(ptr any) float32 {
	return *(ptr.(*float32))
}

func timeValue(ptr any) time.Time {
	return *(ptr.(*time.Time))
}

func TestConversions(t *testing.T) {
	for n, ct := range conversionTests() {
		err := convertAssign(ct.d, ct.s)
		errstr := ""
		if err != nil {
			errstr = err.Error()
		}
		errf := func(format string, args ...any) {
			base := fmt.Sprintf("convertAssign #%d: for %v (%T) -> %T, ", n, ct.s, ct.s, ct.d)
			t.Errorf(base+format, args...)
		}
		if errstr != ct.wanterr {
			errf("got error %q, want error %q", errstr, ct.wanterr)
		}
		if ct.wantstr != "" && ct.wantstr != scanstr {
			errf("want string %q, got %q", ct.wantstr, scanstr)
		}
		if ct.wantbytes != nil && string(ct.wantbytes) != string(scanbytes) {
			errf("want byte %q, got %q", ct.wantbytes, scanbytes)
		}
		if ct.wantraw != nil && string(ct.wantraw) != string(scanraw) {
			errf("want sql.RawBytes %q, got %q", ct.wantraw, scanraw)
		}
		if ct.wantint != 0 && ct.wantint != intValue(ct.d) {
			errf("want int %d, got %d", ct.wantint, intValue(ct.d))
		}
		if ct.wantuint != 0 && ct.wantuint != uintValue(ct.d) {
			errf("want uint %d, got %d", ct.wantuint, uintValue(ct.d))
		}
		if ct.wantf32 != 0 && ct.wantf32 != float32Value(ct.d) {
			errf("want float32 %v, got %v", ct.wantf32, float32Value(ct.d))
		}
		if ct.wantf64 != 0 && ct.wantf64 != float64Value(ct.d) {
			errf("want float32 %v, got %v", ct.wantf64, float64Value(ct.d))
		}
		if bp, boolTest := ct.d.(*bool); boolTest && *bp != ct.wantbool && ct.wanterr == "" {
			errf("want bool %v, got %v", ct.wantbool, *bp)
		}
		if !ct.wanttime.IsZero() && !ct.wanttime.Equal(timeValue(ct.d)) {
			errf("want time %v, got %v", ct.wanttime, timeValue(ct.d))
		}
		if ct.wantnil && *ct.d.(**int64) != nil {
			errf("want nil, got %v", intPtrValue(ct.d))
		}
		if ct.wantptr != nil {
			if *ct.d.(**int64) == nil {
				errf("want pointer to %v, got nil", *ct.wantptr)
			} else if *ct.wantptr != intPtrValue(ct.d) {
				errf("want pointer to %v, got %v", *ct.wantptr, intPtrValue(ct.d))
			}
		}
		if ifptr, ok := ct.d.(*any); ok {
			if !reflect.DeepEqual(ct.wantiface, scaniface) {
				errf("want interface %#v, got %#v", ct.wantiface, scaniface)
				continue
			}
			if srcBytes, ok := ct.s.([]byte); ok {
				dstBytes := (*ifptr).([]byte)
				if len(srcBytes) > 0 && &dstBytes[0] == &srcBytes[0] {
					errf("copy into any didn't copy []byte data")
				}
			}
		}
		if ct.wantusrdef != 0 && ct.wantusrdef != *ct.d.(*userDefined) {
			errf("want userDefined %f, got %f", ct.wantusrdef, *ct.d.(*userDefined))
		}
		if len(ct.wantusrstr) != 0 && ct.wantusrstr != *ct.d.(*userDefinedString) {
			errf("want userDefined %q, got %q", ct.wantusrstr, *ct.d.(*userDefinedString))
		}
	}
}

// Tests that assigning to sql.RawBytes doesn't allocate (and also works).
func TestRawBytesAllocs(t *testing.T) {
	var tests = []struct {
		name string
		in   any
		want string
	}{
		{"uint64", uint64(12345678), "12345678"},
		{"uint32", uint32(1234), "1234"},
		{"uint16", uint16(12), "12"},
		{"uint8", uint8(1), "1"},
		{"uint", uint(123), "123"},
		{"int", int(123), "123"},
		{"int8", int8(1), "1"},
		{"int16", int16(12), "12"},
		{"int32", int32(1234), "1234"},
		{"int64", int64(12345678), "12345678"},
		{"float32", float32(1.5), "1.5"},
		{"float64", float64(64), "64"},
		{"bool", false, "false"},
		{"time", time.Unix(2, 5).UTC(), "1970-01-01T00:00:02.000000005Z"},
	}

	buf := make(sql.RawBytes, 10)
	test := func(name string, in any, want string) {
		if err := convertAssign(&buf, in); err != nil {
			t.Fatalf("%s: convertAssign = %v", name, err)
		}
		match := len(buf) == len(want)
		if match {
			for i, b := range buf {
				if want[i] != b {
					match = false
					break
				}
			}
		}
		if !match {
			t.Fatalf("%s: got %q (len %d); want %q (len %d)", name, buf, len(buf), want, len(want))
		}
	}

	n := testing.AllocsPerRun(100, func() {
		for _, tt := range tests {
			test(tt.name, tt.in, tt.want)
		}
	})

	// The numbers below are only valid for 64-bit interface word sizes,
	// and gc. With 32-bit words there are more convT2E allocs, and
	// with gccgo, only pointers currently go in interface data.
	// So only care on amd64 gc for now.
	measureAllocs := runtime.GOARCH == "amd64" && runtime.Compiler == "gc"

	if n > 0.5 && measureAllocs {
		t.Fatalf("allocs = %v; want 0", n)
	}

	// This one involves a convT2E allocation, string -> any
	n = testing.AllocsPerRun(100, func() {
		test("string", "foo", "foo")
	})
	if n > 1.5 && measureAllocs {
		t.Fatalf("allocs = %v; want max 1", n)
	}
}

// https://golang.org/issues/13905
func TestUserDefinedBytes(t *testing.T) {
	type userDefinedBytes []byte
	var u userDefinedBytes
	v := []byte("foo")

	convertAssign(&u, v)
	if &u[0] == &v[0] {
		t.Fatal("userDefinedBytes got potentially dirty driver memory")
	}
}

func TestAssignZero(t *testing.T) {
	vbool := true
	assignZero(&vbool)
	if vbool != false {
		t.Error("vbool is not zero")
	}

	vstring := "a"
	assignZero(&vstring)
	if vstring != "" {
		t.Error("vstring is not zero")
	}

	vint := int(1)
	assignZero(&vint)
	if vint != 0 {
		t.Error("vint is not zero")
	}

	vint8 := int8(1)
	assignZero(&vint8)
	if vint8 != 0 {
		t.Error("vint8 is not zero")
	}

	vint16 := int16(1)
	assignZero(&vint16)
	if vint16 != 0 {
		t.Error("vint16 is not zero")
	}

	vint32 := int32(1)
	assignZero(&vint32)
	if vint32 != 0 {
		t.Error("vint32 is not zero")
	}

	vint64 := int64(1)
	assignZero(&vint64)
	if vint64 != 0 {
		t.Error("vint64 is not zero")
	}

	vuint := uint(1)
	assignZero(&vuint)
	if vuint != 0 {
		t.Error("vuint is not zero")
	}

	vuint8 := uint8(1)
	assignZero(&vuint8)
	if vuint8 != 0 {
		t.Error("vuint8 is not zero")
	}

	vuint16 := uint16(1)
	assignZero(&vuint16)
	if vuint16 != 0 {
		t.Error("vuint16 is not zero")
	}

	vuint32 := uint32(1)
	assignZero(&vuint32)
	if vuint32 != 0 {
		t.Error("vuint32 is not zero")
	}

	vuint64 := uint64(1)
	assignZero(&vuint64)
	if vuint64 != 0 {
		t.Error("vuint64 is not zero")
	}

	vuintptr := uintptr(1)
	assignZero(&vuintptr)
	if vuintptr != 0 {
		t.Error("vuintptr is not zero")
	}

	vfloat32 := float32(1)
	assignZero(&vfloat32)
	if vfloat32 != 0 {
		t.Error("vfloat32 is not zero")
	}

	vfloat64 := float64(1)
	assignZero(&vfloat64)
	if vfloat64 != 0 {
		t.Error("vfloat64 is not zero")
	}

	vinterface := any("a")
	assignZero(&vinterface)
	if vinterface != nil {
		t.Error("vinterface is not zero")
	}

	vbytes := []byte("a")
	assignZero(&vbytes)
	if vbytes != nil {
		t.Error("vbytes is not zero")
	}

	rawBytes := sql.RawBytes("a")
	assignZero(&rawBytes)
	if rawBytes != nil {
		t.Error("rawBytes is not zero")
	}

	vuserDefined := userDefined(1)
	assignZero(&vuserDefined)
	if vuserDefined != 0 {
		t.Error("vuserDefined is not zero")
	}
}
