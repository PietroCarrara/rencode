package rencode

import (
	"encoding/binary"
	"math"
)

// Encodes a int into a byte slice and returns how many bytes were used
// If the slice is too small this function panics
// Takes at most 8 bytes
func putInt(buf []byte, n int64) int {
	size := -1

	if math.MinInt8 <= n && n <= math.MaxInt8 {
		size = 8 / 8
		buf[0] = byte(n)
	} else if math.MinInt16 <= n && n <= math.MaxInt16 {
		size = 16 / 8
		binary.BigEndian.PutUint16(buf, uint16(n))
	} else if math.MinInt32 <= n && n <= math.MaxInt32 {
		size = 32 / 8
		binary.BigEndian.PutUint32(buf, uint32(n))
	} else if math.MinInt64 <= n && n <= math.MaxInt64 {
		size = 64 / 8
		binary.BigEndian.PutUint64(buf, uint64(n))
	}

	if size <= 0 {
		panic("int value outside of range")
	}

	return size
}
