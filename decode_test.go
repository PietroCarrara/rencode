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

func TestDecodingIntSmall(t *testing.T) {
	t.Parallel()

	fail(testDecodeIntSmall([]byte{1}, 1), t)
	fail(testDecodeIntSmall([]byte{10}, 10), t)
	fail(testDecodeIntSmall([]byte{84}, -15), t)
	fail(testDecodeIntSmall([]byte{101}, -32), t)
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

func TestDecodingStringFixed(t *testing.T) {
	t.Parallel()

	fail(testDecodeStringFixed([]byte{140, 104, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100}, "hello, world"), t)
	fail(testDecodeStringFixed([]byte{140, 228, 189, 160, 229, 165, 189, 228, 184, 150, 231, 149, 140}, "你好世界"), t)
	fail(testDecodeStringFixed([]byte{149, 227, 129, 147, 227, 130, 147, 227, 129, 171, 227, 129, 161, 227, 129, 175, 228, 184, 150, 231, 149, 140}, "こんにちは世界"), t)
	fail(testDecodeStringFixed([]byte{156, 208, 151, 208, 180, 209, 128, 208, 176, 208, 178, 209, 129, 209, 130, 208, 178, 209, 131, 208, 185, 44, 32, 208, 188, 208, 184, 209, 128}, "Здравствуй, мир"), t)
	fail(testDecodeStringFixed([]byte{171, 116, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103}, "the quick brown fox jumps over the lazy dog"), t)
}

func TestDecodingSliceVarLength(t *testing.T) {
	t.Parallel()

	fail(testDecodeSliceVarLength([]byte{chrList, chrTerm}, list{}), t)
	fail(testDecodeSliceVarLength([]byte{chrList, 1, 2, 3, chrTerm}, list{1, 2, 3}), t)
	fail(testDecodeSliceVarLength([]byte{chrList, chrList, chrTerm, chrTerm}, list{list{}}), t)
}

func testDecodeIntStr(value []byte, target int64) error {
	var val int64
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt8(value []byte, target int8) error {
	var val int8
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt16(value []byte, target int16) error {
	var val int16
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt32(value []byte, target int32) error {
	var val int32
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeInt64(value []byte, target int64) error {
	var val int64
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeIntSmall(value []byte, target int) error {
	var val int
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %d, but got %d", target, val)
	}

	return nil
}

func testDecodeFloat32(value []byte, target float32) error {
	var val float32
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %f, but got %f", target, val)
	}

	return nil
}

func testDecodeFloat64(value []byte, target float64) error {
	var val float64
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected %f, but got %f", target, val)
	}

	return nil
}

func testDecodeStringFixed(value []byte, target string) error {
	var val string
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if target != val {
		return fmt.Errorf("expected \"%s\", but got \"%s\"", target, val)
	}

	return nil
}

func testDecodeSliceVarLength(value []byte, target list) error {
	var val list
	bytes, err := Decode(value, &val)

	if err != nil {
		return err
	}
	if bytes != len(value) {
		return fmt.Errorf("had %d bytes, but only %d were consumed", len(value), bytes)
	}
	if !listEquals(val, target) {
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

// TODO: Be smart
// Really poor list comparison, but gets the job done
func listEquals(a, b list) bool {
	str1 := fmt.Sprint(a)
	str2 := fmt.Sprint(b)

	return str1 == str2
}
