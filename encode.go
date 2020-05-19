package rencode

import (
	"encoding/binary"
	"math"
	"strconv"
)

var chrIntMap = map[int]byte{
	1: chrInt1,
	2: chrInt2,
	4: chrInt4,
	8: chrInt8,
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
