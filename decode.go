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
		return 0, fmt.Errorf("can't assign value to %v", ref)
	}

	ptr := val.Elem()
	chr := buf[0]

	switch chr {
	case chrList:
		return decodeSliceVariable(buf, ptr)
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

	if intPosFixedStart <= chr && chr <= intPosFixedStart+intPosFixedCount ||
		intNegFixedStart <= chr && chr < intNegFixedStart+intNegFixedCount {

		return decodeIntSmall(buf, ptr)
	}

	if strFixedStart <= chr && chr < strFixedStart+strFixedCount {
		return decodeStringFixed(buf, ptr)
	}

	if '1' <= chr && chr <= '9' {
		return decodeStringVariable(buf, ptr)
	}

	if listFixedStart <= chr && chr <= listFixedStart+listFixedCount-1 {
		return decodeSliceFixed(buf, ptr)
	}

	panic("this line should be theoretically impossible to be executed")
}

func decodeIntStr(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if chr != chrInt {
		return 0, fmt.Errorf("expected chr byte to be %d, but found %d", chrInt1, chr)
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

	val.Set(reflect.ValueOf(integer).Convert(val.Type()))

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

	val.Set(reflect.ValueOf(integer).Convert(val.Type()))

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

	val.Set(reflect.ValueOf(integer).Convert(val.Type()))

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

	val.Set(reflect.ValueOf(integer).Convert(val.Type()))

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

	val.Set(reflect.ValueOf(integer).Convert(val.Type()))

	return 9, nil
}

func decodeIntSmall(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if intPosFixedStart <= chr && chr <= intPosFixedStart+intPosFixedCount {
		integer := int8(chr) - intPosFixedStart
		val.Set(reflect.ValueOf(integer).Convert(val.Type()))
		return 1, nil
	}

	if intNegFixedStart <= chr && chr < intNegFixedStart+intNegFixedCount {
		integer := (int8(chr) - intNegFixedStart + 1) * -1
		val.Set(reflect.ValueOf(integer).Convert(val.Type()))
		return 1, nil
	}

	return 0, fmt.Errorf(
		"expected chr byte to be in range [%d, %d] or [%d, %d], but found %d",
		intPosFixedStart,
		intPosFixedStart+intPosFixedCount,
		intNegFixedStart,
		intNegFixedStart+intNegFixedCount,
		chr,
	)
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

	val.Set(reflect.ValueOf(float).Convert(val.Type()))

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

	val.Set(reflect.ValueOf(float).Convert(val.Type()))

	return 9, nil
}

func decodeStringFixed(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if !(strFixedStart <= chr && chr <= strFixedStart+strFixedCount) {
		return 0, fmt.Errorf(
			"expected chr byte to be in range [%d, %d], but found %d",
			strFixedStart,
			strFixedStart+strFixedCount,
			chr,
		)
	}

	length := int(chr - strFixedStart)
	bytes := string(buf[1 : length+1])

	val.Set(reflect.ValueOf(bytes).Convert(val.Type()))

	return length + 1, nil
}

func decodeStringVariable(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if !('1' <= chr && chr <= '9') {
		return 0, fmt.Errorf("expected chr byte to be in range [%d, %d], but found %d", '1', '9', chr)
	}

	lenStr := ""

	for buf[0] != ':' {
		lenStr += string(buf[0])
		buf = buf[1:]

		if len(buf) <= 0 {
			return 0, fmt.Errorf("incomplete stream: could not find ':' while decoding string")
		}
	}

	buf = buf[1:] // Skip ':'

	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return 0, err
	}

	str := string(buf[:length])
	val.Set(reflect.ValueOf(str).Convert(val.Type()))

	return len(lenStr) + 1 + length, nil // length string + ':' + actual string
}

func decodeSliceFixed(buf []byte, val reflect.Value) (int, error) {
	chr := buf[0]

	if !(listFixedStart <= chr && chr <= listFixedStart+listFixedCount-1) {
		return 0, fmt.Errorf(
			"expected chr byte to be in range [%d, %d], but found %d",
			listFixedStart,
			listFixedStart+listFixedCount,
			chr,
		)
	}

	typ := val.Type()
	kind := typ.Kind()
	if kind != reflect.Slice && kind != reflect.Array && kind != reflect.Interface {
		return 0, fmt.Errorf("can't decode list into type \"%s\"", typ)
	}

	if kind == reflect.Interface {
		typ = reflect.SliceOf(typ)
	}

	buf = buf[1:]
	length := int(chr - listFixedStart)
	slice := reflect.MakeSlice(typ, length, length)
	bytes := 0

	for i := 0; i < length; i++ {
		var data interface{}

		n, err := Decode(buf, &data)
		if err != nil {
			return 0, err
		}
		slice.Index(i).Set(reflect.ValueOf(data))

		bytes += n
		buf = buf[n:]
	}

	val.Set(slice)

	return bytes + 1, nil
}

func decodeSliceVariable(buf []byte, val reflect.Value) (int, error) {

	if buf[0] != chrList {
		return 0, fmt.Errorf("expected chr byte %d, found %d", chrList, buf[0])
	}

	typ := val.Type()
	kind := typ.Kind()
	if kind != reflect.Slice && kind != reflect.Array && kind != reflect.Interface {
		return 0, fmt.Errorf("can't decode list into type \"%s\"", typ)
	}

	if kind == reflect.Interface {
		typ = reflect.SliceOf(typ)
	}

	slice := reflect.MakeSlice(typ, 0, 0)
	bytes := 0
	buf = buf[1:]

	for buf[0] != chrTerm {
		var data interface{}
		n, err := Decode(buf, &data)

		if err != nil {
			return 0, err
		}

		slice = reflect.Append(slice, reflect.ValueOf(data))

		buf = buf[n:]
		bytes += n
	}

	val.Set(slice)

	return bytes + 2, nil // bytes + chrList + chrTerm
}
