package rencode

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
)

var chrIntMap = map[int]byte{
	1: chrInt1,
	2: chrInt2,
	4: chrInt4,
	8: chrInt8,
}

func Encode(data interface{}) ([]byte, error) {
	typeof := reflect.TypeOf(data)
	valueof := reflect.ValueOf(data)

	if data == nil {
		return encodeNil(), nil
	}

	switch typeof.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt(valueof.Int())
	case reflect.Float32, reflect.Float64:
		return encodeFloat(valueof.Float())
	case reflect.Bool:
		return encodeBool(valueof.Bool()), nil
	case reflect.String:
		return encodeString(valueof.String())
	case reflect.Array, reflect.Slice:
		return encodeSlice(valueof.Interface())
	case reflect.Map:
		return encodeMap(data)
	}

	return nil, fmt.Errorf("can't encode value of type %s, kind %s", typeof, typeof.Kind())
}

func encodeInt(n int64) ([]byte, error) {
	buff := make([]byte, 9)

	if 0 <= n && n < intPosFixedCount {
		size := putInt(buff, n)
		return buff[:size], nil
	} else if -intNegFixedCount <= n && n < 0 {
		size := putInt(buff, intNegFixedStart-1-n)
		return buff[:size], nil
	} else {
		size := putInt(buff[1:], n)
		buff[0] = chrIntMap[size]
		return buff[:size+1], nil
	}
}

func encodeFloat(n float64) ([]byte, error) {
	buf := make([]byte, 9)

	if n <= math.MaxFloat32 {
		buf[0] = chrFloat32
		bits := math.Float32bits(float32(n))
		binary.BigEndian.PutUint32(buf[1:], bits)
		return buf[:4+1], nil
	}

	buf[0] = chrFloat64
	bits := math.Float64bits(n)
	binary.BigEndian.PutUint64(buf[1:], bits)
	return buf[:8+1], nil
}

func encodeBool(n bool) []byte {
	buff := make([]byte, 1)

	if n {
		buff[0] = chrTrue
	} else {
		buff[0] = chrFalse
	}

	return buff
}

func encodeNil() []byte {
	return []byte{chrNone}
}

func encodeString(s string) ([]byte, error) {
	var bytes []byte

	if len(s) < strFixedCount {
		bytes = append(bytes, byte(strFixedStart+len(s)))
		bytes = append(bytes, s...)
		return bytes, nil
	}

	size := strconv.Itoa(len(s))
	bytes = append(bytes, size...)
	bytes = append(bytes, ':')
	bytes = append(bytes, s...)
	return bytes, nil
}

func encodeSlice(data interface{}) ([]byte, error) {
	var bytes []byte

	typeof := reflect.TypeOf(data)
	valueof := reflect.ValueOf(data)

	if typeof.Kind() != reflect.Array && typeof.Kind() != reflect.Slice {
		panic("encoding non-slice type on encodeSlice")
	}

	length := valueof.Len()

	if length < listFixedCount {
		bytes = append(bytes, byte(listFixedStart+length))
	} else {
		bytes = append(bytes, chrList)
	}

	for i := 0; i < length; i++ {
		val, err := Encode(valueof.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, val...)
	}

	if length >= listFixedCount {
		bytes = append(bytes, chrTerm)
	}

	return bytes, nil
}

func encodeMap(data interface{}) ([]byte, error) {
	var bytes []byte

	typeof := reflect.TypeOf(data)
	valueof := reflect.ValueOf(data)

	if typeof.Kind() != reflect.Map {
		panic("encoding non-map type on encodeMap")
	}

	keys := valueof.MapKeys()

	if len(keys) < dictFixedCount {
		bytes = append(bytes, byte(dictFixedStart+len(keys)))
	} else {
		bytes = append(bytes, chrDict)
	}

	// Sort keys to make the ordering constant
	// bewteen encodings (helps out when testing)
	sort.Slice(keys, func(i, j int) bool {
		iVal := fmt.Sprint(keys[i].Interface())
		jVal := fmt.Sprint(keys[j].Interface())

		return iVal < jVal
	})

	for _, key := range keys {
		val := valueof.MapIndex(key)

		bin, err := Encode(key.Interface())
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, bin...)

		bin, err = Encode(val.Interface())
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, bin...)
	}

	if len(keys) >= dictFixedCount {
		bytes = append(bytes, chrTerm)
	}

	return bytes, nil
}
