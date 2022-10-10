// Modified from: database/sql/convert.go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style

package rel

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var _, localTimeOffset = time.Now().Local().Zone()

// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
// dest will be set to zero value if src is nil.
// this function assumes dest will never be nil.
func convertAssign(dest, src any) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			*d = s
			return nil
		case *[]byte:
			*d = []byte(s)
			return nil
		case *sql.RawBytes:
			*d = append((*d)[:0], s...)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			*d = string(s)
			return nil
		case *any:
			*d = cloneBytes(s)
			return nil
		case *[]byte:
			*d = cloneBytes(s)
			return nil
		case *sql.RawBytes:
			*d = s
			return nil
		}
	case time.Time:
		switch d := dest.(type) {
		case *time.Time:
			// make sure timezone equal for test assertion.
			if _, offset := s.Zone(); offset == localTimeOffset {
				*d = s.Local()
			} else {
				*d = s
			}
			return nil
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		case *sql.RawBytes:
			*d = s.AppendFormat((*d)[:0], time.RFC3339Nano)
			return nil
		}
	case nil:
		assignZero(dest)
		return nil
	}

	var sv reflect.Value

	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			if s, ok := asString(src); ok {
				*d = s
				return nil
			}
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *sql.RawBytes:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes([]byte(*d)[:0], sv); ok {
			*d = sql.RawBytes(b)
			return nil
		}
	case *bool:
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool)
		}
		return err
	case *any:
		*d = src
		return nil
	}

	dpv := reflect.ValueOf(dest)

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	// The following conversions use a string value as an intermediate representation
	// to convert between various numeric types.
	//
	// This also allows scanning into user defined types such as "type Int int64".
	// For symmetry, also check for string destination types.
	if s, ok := asString(src); ok {
		switch dv.Kind() {
		case reflect.Ptr:
			dv.Set(reflect.New(dv.Type().Elem()))
			return convertAssign(dv.Interface(), src)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
			if err != nil {
				// The errors that ParseInt returns have concrete type *NumError
				err = err.(*strconv.NumError).Err
				return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
			}
			dv.SetInt(i64)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
			if err != nil {
				// The errors that ParseUint returns have concrete type *NumError
				err = err.(*strconv.NumError).Err
				return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
			}
			dv.SetUint(u64)
			return nil
		case reflect.Float32, reflect.Float64:
			f64, err := strconv.ParseFloat(s, dv.Type().Bits())
			if err != nil {
				// The errors that ParseFloat returns have concrete type *NumError
				err = err.(*strconv.NumError).Err
				return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
			}
			dv.SetFloat(f64)
			return nil
		case reflect.String:
			dv.SetString(s)
			return nil
		}
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func asString(src any) (string, bool) {
	switch v := src.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), true
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64), true
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32), true
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), true
	}
	return "", false
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}

func assignZero(dest any) {
	switch d := dest.(type) {
	case *bool:
		*d = false
	case *string:
		*d = ""
	case *int:
		*d = 0
	case *int8:
		*d = 0
	case *int16:
		*d = 0
	case *int32:
		*d = 0
	case *int64:
		*d = 0
	case *uint:
		*d = 0
	case *uint8:
		*d = 0
	case *uint16:
		*d = 0
	case *uint32:
		*d = 0
	case *uint64:
		*d = 0
	case *uintptr:
		*d = 0
	case *float32:
		*d = 0
	case *float64:
		*d = 0
	case *any:
		*d = nil
	case *[]byte:
		*d = nil
	case *sql.RawBytes:
		*d = nil
	default:
		rv := reflect.ValueOf(dest)
		rv.Elem().Set(reflect.Zero(rv.Type().Elem()))
	}
}
