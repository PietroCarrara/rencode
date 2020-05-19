package rencode

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
