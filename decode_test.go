package rencode

import (
	"fmt"
	"math"
	"testing"
)

func TestDecodingBool(t *testing.T) {
	t.Parallel()

	tru := false
	fal := true

	nTru, errTru := Decode([]byte{chrTrue}, &tru)
	nFal, errFal := Decode([]byte{chrFalse}, &fal)

	fail(errTru, t)
	fail(errFal, t)

	if nTru != 1 || nFal != 1 {
		t.Errorf("used more than 1 byte decoding booleans")
	}

	if tru != true || fal != false {
		t.Errorf("did not decode booleans propperly")
	}
}

func TestDecodingNil(t *testing.T) {
	t.Parallel()

	b := true
	i := 123
	str := "hello, world"
	arr := []int{1, 2, 3}
	dict := map[string]int{"1": 1, "2": 2}

	null := []byte{chrNone}

	Decode(null, &b)
	Decode(null, &i)
	Decode(null, &str)
	Decode(null, &arr)
	Decode(null, &dict)

	if b != false {
		t.Error("bool is a non-zero value")
	}
	if i != 0 {
		t.Error("int is a non-zero value")
	}
	if arr != nil {
		t.Error("array is a non-zero value")
	}
	if dict != nil {
		t.Error("map is a non-zero value")
	}
}

func TestDecodingIntStr(t *testing.T) {
	t.Parallel()

	fail(testDecodeIntStr(intStr("0024"), 24), t)
	fail(testDecodeIntStr(intStr("12345"), 12345), t)
	fail(testDecodeIntStr(intStr("-1024"), -1024), t)
	fail(testDecodeIntStr(intStr("-123456789"), -123456789), t)
}

func TestDecodingInt8(t *testing.T) {
	t.Parallel()

	fail(testDecodeInt32([]byte{62, 0}, 0), t)
	fail(testDecodeInt32([]byte{62, 128}, math.MinInt8), t)
	fail(testDecodeInt32([]byte{62, 127}, math.MaxInt8), t)
}

func TestDecodingInt16(t *testing.T) {
	t.Parallel()

	fail(testDecodeInt32([]byte{63, 0, 0}, 0), t)
	fail(testDecodeInt32([]byte{63, 128, 0}, math.MinInt16), t)
	fail(testDecodeInt32([]byte{63, 127, 255}, math.MaxInt16), t)
}

func TestDecodingInt32(t *testing.T) {
	t.Parallel()

	fail(testDecodeInt32([]byte{64, 0, 0, 0, 0}, 0), t)
	fail(testDecodeInt32([]byte{64, 128, 0, 0, 0}, math.MinInt32), t)
	fail(testDecodeInt32([]byte{64, 127, 255, 255, 255}, math.MaxInt32), t)
}

func TestDecodingInt64(t *testing.T) {
	t.Parallel()

	fail(testDecodeInt64([]byte{65, 0, 0, 0, 0, 0, 0, 0, 0}, 0), t)
	fail(testDecodeInt64([]byte{65, 128, 0, 0, 0, 0, 0, 0, 0}, math.MinInt64), t)
	fail(testDecodeInt64([]byte{65, 127, 255, 255, 255, 255, 255, 255, 255}, math.MaxInt64), t)
}

func TestDecodingFloat32(t *testing.T) {
	t.Parallel()

	fail(testDecodeFloat32([]byte{66, 65, 120, 0, 0}, 15.5), t)
	fail(testDecodeFloat32([]byte{66, 65, 69, 112, 164}, 12.34), t)
	fail(testDecodeFloat32([]byte{66, 195, 128, 89, 154}, -256.7), t)
}

func TestDecodingFloat64(t *testing.T) {
	t.Parallel()

	fail(testDecodeFloat64([]byte{44, 64, 47, 0, 0, 0, 0, 0, 0}, 15.5), t)
	fail(testDecodeFloat64([]byte{44, 192, 112, 11, 51, 51, 51, 51, 51}, -256.7), t)
	fail(testDecodeFloat64([]byte{44, 64, 40, 174, 20, 122, 225, 71, 174}, 12.34), t)
}

func testDecodeIntStr(value []byte, target int64) error {
	var val int64
	_, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt8(value []byte, target int8) error {
	var val int8
	_, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt16(value []byte, target int16) error {
	var val int16
	_, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt32(value []byte, target int32) error {
	var val int32
	_, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt64(value []byte, target int64) error {
	var val int64
	_, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeFloat32(value []byte, target float32) error {
	var val float32
	_, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if target != val {
		return fmt.Errorf("expected %f, but got %f", target, val)
	}

	return nil
}

func testDecodeFloat64(value []byte, target float64) error {
	var val float64
	_, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if target != val {
		return fmt.Errorf("expected %f, but got %f", target, val)
	}

	return nil
}

func intStr(s string) []byte {
	bytes := []byte{chrInt}
	bytes = append(bytes, s...)
	bytes = append(bytes, chrTerm)

	return bytes
}
