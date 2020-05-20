package rencode

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strconv"
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
		return decodeIntStr(buf, ptr)
	case chrInt1:
		return decodeInt8(buf, ptr)
	case chrInt2:
		return decodeInt16(buf, ptr)
	case chrInt4:
		return decodeInt32(buf, ptr)
	case chrInt8:
		return decodeInt64(buf, ptr)
	case chrFloat32:
		return decodeFloat32(buf, ptr)
	case chrFloat64:
		return decodeFloat64(buf, ptr)
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

func decodeIntStr(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if chr != chrInt {
		return 0, fmt.Errorf("decodeInt8: expected chr byte to be %d, but found %d", chrInt1, chr)
	}

	length := len(buf)
	index := -1

	for i := 0; i < length; i++ {
		if buf[i] == chrTerm {
			index = i
			break
		}
	}

	if index == -1 {
		return 0, fmt.Errorf("could not find chrTerm on stream")
	}

	integer, err := strconv.ParseInt(string(buf[1:index]), 10, 64)
	if err != nil {
		return 0, err
	}

	val.SetInt(integer)

	return index + 1, nil
}

func decodeInt8(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if chr != chrInt1 {
		return 0, fmt.Errorf("expected chr byte to be %d, but found %d", chrInt1, chr)
	}

	if len(buf[1:]) < 1 {
		return 0, fmt.Errorf("incomplete stream: decoding 1 bytes but found %d", len(buf[1:]))
	}

	integer := int8(buf[1])

	val.SetInt(int64(integer))

	return 2, nil

}

func decodeInt16(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if chr != chrInt2 {
		return 0, fmt.Errorf("expected chr byte to be %d, but found %d", chrInt2, chr)
	}

	if len(buf[1:]) < 2 {
		return 0, fmt.Errorf("incomplete stream: decoding 2 bytes but found %d", len(buf[1:]))
	}

	integer := int16(binary.BigEndian.Uint16(buf[1:]))

	val.SetInt(int64(integer))

	return 3, nil

}

func decodeInt32(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if chr != chrInt4 {
		return 0, fmt.Errorf("expected chr byte to be %d, but found %d", chrInt4, chr)
	}

	if len(buf[1:]) < 4 {
		return 0, fmt.Errorf("incomplete stream: decoding 4 bytes but found %d", len(buf[1:]))
	}

	integer := int32(binary.BigEndian.Uint32(buf[1:]))

	val.SetInt(int64(integer))

	return 5, nil

}

func decodeInt64(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if chr != chrInt8 {
		return 0, fmt.Errorf("expected chr byte to be %d, but found %d", chrInt8, chr)
	}

	if len(buf[1:]) < 8 {
		return 0, fmt.Errorf("incomplete stream: decoding 8 bytes but found %d", len(buf[1:]))
	}

	integer := int64(binary.BigEndian.Uint64(buf[1:]))

	val.SetInt(integer)

	return 9, nil
}

func decodeFloat32(buf []byte, val reflect.Value) (int, error) {
	if buf[0] != chrFloat32 {
		return 0, fmt.Errorf("expected chr byte %d, found %d", chrFloat32, buf[0])
	}

	if len(buf[1:]) < 4 {
		return 0, fmt.Errorf("incomplete stream: decoding 4 bytes but found %d", len(buf[1:]))
	}

	bits := binary.BigEndian.Uint32(buf[1:])
	float := math.Float32frombits(bits)

	val.SetFloat(float64(float))

	return 5, nil
}

func decodeFloat64(buf []byte, val reflect.Value) (int, error) {
	if buf[0] != chrFloat64 {
		return 0, fmt.Errorf("expected chr byte %d, found %d", chrFloat64, buf[0])
	}

	if len(buf[1:]) < 8 {
		return 0, fmt.Errorf("incomplete stream: decoding 4 bytes but found %d", len(buf[1:]))
	}

	bits := binary.BigEndian.Uint64(buf[1:])
	float := math.Float64frombits(bits)

	val.SetFloat(float)

	return 9, nil
}
