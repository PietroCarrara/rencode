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

	fail(testDecodeString([]byte{140, 104, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100}, "hello, world"), t)
	fail(testDecodeString([]byte{140, 228, 189, 160, 229, 165, 189, 228, 184, 150, 231, 149, 140}, "你好世界"), t)
	fail(testDecodeString([]byte{149, 227, 129, 147, 227, 130, 147, 227, 129, 171, 227, 129, 161, 227, 129, 175, 228, 184, 150, 231, 149, 140}, "こんにちは世界"), t)
	fail(testDecodeString([]byte{156, 208, 151, 208, 180, 209, 128, 208, 176, 208, 178, 209, 129, 209, 130, 208, 178, 209, 131, 208, 185, 44, 32, 208, 188, 208, 184, 209, 128}, "Здравствуй, мир"), t)
	fail(testDecodeString([]byte{171, 116, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103}, "the quick brown fox jumps over the lazy dog"), t)
}

func TestDecodingStringVariable(t *testing.T) {
	t.Parallel()

	fail(testDecodeString([]byte{55, 48, 50, 58, 228, 184, 173, 229, 155, 189, 229, 190, 136, 230, 151, 169, 229, 176, 177, 229, 135, 186, 231, 142, 176, 228, 186, 134, 228, 184, 147, 233, 151, 168, 231, 148, 168, 228, 186, 142, 229, 144, 175, 232, 146, 153, 231, 154, 132, 232, 175, 134, 229, 173, 151, 232, 175, 190, 230, 156, 172, 239, 188, 140, 231, 167, 166, 228, 187, 163, 229, 135, 186, 231, 142, 176, 231, 154, 132, 230, 156, 137, 227, 128, 138, 229, 128, 137, 233, 162, 137, 231, 175, 135, 227, 128, 139, 227, 128, 129, 227, 128, 138, 231, 136, 176, 229, 142, 134, 231, 175, 135, 227, 128, 139, 239, 188, 140, 230, 177, 137, 228, 187, 163, 229, 136, 153, 230, 156, 137, 229, 143, 184, 233, 169, 172, 231, 155, 184, 229, 166, 130, 231, 154, 132, 227, 128, 138, 229, 135, 161, 229, 176, 134, 231, 175, 135, 227, 128, 139, 227, 128, 129, 232, 180, 190, 233, 178, 130, 231, 154, 132, 227, 128, 138, 230, 187, 130, 229, 150, 156, 231, 175, 135, 227, 128, 139, 227, 128, 129, 232, 148, 161, 233, 130, 149, 231, 154, 132, 227, 128, 138, 229, 138, 157, 229, 173, 166, 231, 175, 135, 227, 128, 139, 227, 128, 129, 229, 143, 178, 230, 184, 184, 231, 154, 132, 227, 128, 138, 230, 128, 165, 229, 176, 177, 231, 171, 160, 227, 128, 139, 239, 188, 140, 228, 184, 137, 229, 155, 189, 230, 151, 182, 228, 187, 163, 230, 156, 137, 227, 128, 138, 229, 159, 164, 232, 139, 141, 227, 128, 139, 227, 128, 129, 227, 128, 138, 229, 185, 191, 232, 139, 141, 227, 128, 139, 227, 128, 129, 227, 128, 138, 229, 167, 139, 229, 173, 166, 231, 175, 135, 227, 128, 139, 231, 173, 137, 239, 188, 140, 232, 191, 153, 228, 186, 155, 228, 189, 156, 229, 147, 129, 228, 184, 173, 229, 143, 170, 230, 156, 137, 227, 128, 138, 230, 128, 165, 229, 176, 177, 231, 171, 160, 227, 128, 139, 229, 175, 185, 229, 144, 142, 228, 184, 150, 228, 186, 167, 231, 148, 159, 228, 186, 134, 229, 189, 177, 229, 147, 141, 239, 188, 140, 229, 133, 182, 228, 189, 153, 229, 189, 177, 229, 147, 141, 228, 184, 141, 229, 164, 167, 239, 188, 140, 227, 128, 138, 230, 128, 165, 229, 176, 177, 231, 171, 160, 227, 128, 139, 232, 153, 189, 231, 132, 182, 230, 152, 175, 227, 128, 138, 232, 139, 141, 233, 162, 137, 231, 175, 135, 227, 128, 139, 228, 185, 139, 229, 144, 142, 232, 190, 131, 231, 170, 129, 229, 135, 186, 231, 154, 132, 229, 176, 143, 229, 173, 166, 228, 185, 139, 228, 185, 166, 239, 188, 140, 228, 189, 134, 231, 148, 177, 228, 186, 142, 230, 181, 129, 228, 188, 160, 228, 184, 173, 229, 135, 186, 231, 142, 176, 228, 186, 134, 231, 167, 141, 231, 167, 141, 233, 151, 174, 233, 162, 152, 239, 188, 140, 229, 133, 182, 230, 157, 131, 229, 168, 129, 230, 128, 167, 229, 136, 176, 229, 141, 151, 229, 140, 151, 230, 156, 157, 230, 151, 182, 229, 183, 178, 229, 164, 167, 228, 184, 141, 229, 166, 130, 229, 137, 141, 239, 188, 140, 232, 128, 140, 232, 191, 153, 228, 184, 128, 230, 151, 182, 230, 156, 159, 229, 135, 186, 231, 142, 176, 231, 154, 132, 228, 184, 128, 228, 186, 155, 229, 144, 175, 232, 146, 153, 232, 175, 187, 231, 137, 169, 229, 166, 130, 227, 128, 138, 229, 186, 173, 232, 175, 176, 227, 128, 139, 227, 128, 129, 227, 128, 138, 232, 175, 130, 229, 185, 188, 227, 128, 139, 228, 185, 139, 231, 177, 187, 239, 188, 140, 229, 143, 175, 232, 175, 187, 230, 128, 167, 230, 156, 137, 233, 153, 144, 227, 128, 130, 229, 176, 177, 230, 152, 175, 229, 156, 168, 232, 191, 153, 230, 160, 183, 231, 154, 132, 232, 131, 140, 230, 153, 175, 228, 184, 139, 239, 188, 140, 227, 128, 138, 229, 141, 131, 229, 173, 151, 230, 150, 135, 227, 128, 139, 233, 151, 174, 228, 184, 150, 228, 186, 134}, "中国很早就出现了专门用于启蒙的识字课本，秦代出现的有《倉颉篇》、《爰历篇》，汉代则有司马相如的《凡将篇》、贾鲂的《滂喜篇》、蔡邕的《劝学篇》、史游的《急就章》，三国时代有《埤苍》、《广苍》、《始学篇》等，这些作品中只有《急就章》对后世产生了影响，其余影响不大，《急就章》虽然是《苍颉篇》之后较突出的小学之书，但由于流传中出现了种种问题，其权威性到南北朝时已大不如前，而这一时期出现的一些启蒙读物如《庭诰》、《诂幼》之类，可读性有限。就是在这样的背景下，《千字文》问世了"), t)
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

func testDecodeString(value []byte, target string) error {
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
