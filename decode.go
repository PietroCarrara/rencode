package rencode

import (
	"fmt"
	"reflect"
)

func Decode(buf []byte, ref interface{}) (int, error) {
	val := reflect.ValueOf(ref)

	if val.Kind() != reflect.Ptr || val.IsNil() {
		return 0, fmt.Errorf("can't assign value to %v", val.Interface())
	}

	ptr := val.Elem()

	switch buf[0] {
	case chrList:
		// return decodeSlice
	case chrDict:
		// return decodeMap
	case chrInt:
		// return decodeInt
	case chrInt1:
		// return decodeInt8
	case chrInt2:
		// return decodeInt16
	case chrInt4:
		// return decodeInt32
	case chrInt8:
		// return decodeInt64
	case chrFloat32:
		// return decodeFloat32
	case chrFloat64:
		// return decodeFloat64
	case chrTrue:
		ptr.Set(reflect.ValueOf(true))
		return 1, nil
	case chrFalse:
		ptr.Set(reflect.ValueOf(false))
		return 1, nil
	case chrNone:
		ptr.Set(reflect.Zero(ptr.Type()))
		return 1, nil
	}

	panic("TODO")
}
