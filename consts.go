package rencode

import "reflect"

const (
	defaultFloatBits = reflect.Float32
	maxIntLength     = reflect.Int64

	chrList    byte = 59
	chrDict    byte = 60
	chrInt     byte = 61
	chrInt1    byte = 62
	chrInt2    byte = 63
	chrInt4    byte = 64
	chrInt8    byte = 65
	chrFloat32 byte = 66
	chrFloat64 byte = 44
	chrTrue    byte = 67
	chrFalse   byte = 68
	chrNone    byte = 69
	chrTerm    byte = 127

	intPosFixedStart = 0
	intPosFixedCount = 44

	dictFixedStart = 102
	dictFixedCount = 25

	intNegFixedStart = 70
	intNegFixedCount = 32

	strFixedStart = 128
	strFixedCount = 64

	listFixedStart = strFixedStart + strFixedCount
	listFixedCount = 64
)
