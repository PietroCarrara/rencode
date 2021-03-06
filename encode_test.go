package rencode

import (
	"fmt"
	"math"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	t.Parallel()

	fail(testEncodingInt(0, []byte{0}), t)
	fail(testEncodingInt(1, []byte{1}), t)
	fail(testEncodingInt(2, []byte{2}), t)
	fail(testEncodingInt(-1, []byte{70}), t)
	fail(testEncodingInt(-10, []byte{79}), t)
	fail(testEncodingInt(256, []byte{63, 1, 0}), t)
	fail(testEncodingInt(255, []byte{63, 0, 255}), t)
	fail(testEncodingInt(-256, []byte{63, 255, 0}), t)
	fail(testEncodingInt(math.MaxInt64, []byte{65, 127, 255, 255, 255, 255, 255, 255, 255}), t)
}

func TestEncodeFloat(t *testing.T) {
	t.Parallel()

	fail(testEncodingFloat(0, []byte{66, 0, 0, 0, 0}), t)
	fail(testEncodingFloat(-0.25, []byte{66, 190, 128, 0, 0}), t)
	fail(testEncodingFloat(256.5, []byte{66, 67, 128, 64, 0}), t)
	fail(testEncodingFloat(-256.5, []byte{66, 195, 128, 64, 0}), t)
	fail(testEncodingFloat(-256.5, []byte{66, 195, 128, 64, 0}), t)
	fail(testEncodingFloat(1024.333, []byte{66, 68, 128, 10, 168}), t)
	fail(testEncodingFloat(1234567890, []byte{66, 78, 147, 44, 6}), t)
	fail(testEncodingFloat(-1234567890, []byte{66, 206, 147, 44, 6}), t)
}

func TestEncodeString(t *testing.T) {
	t.Parallel()

	fail(testEncodingString("hello, world", []byte{140, 104, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100}), t)
	fail(testEncodingString("你好世界", []byte{140, 228, 189, 160, 229, 165, 189, 228, 184, 150, 231, 149, 140}), t)
	fail(testEncodingString("こんにちは世界", []byte{149, 227, 129, 147, 227, 130, 147, 227, 129, 171, 227, 129, 161, 227, 129, 175, 228, 184, 150, 231, 149, 140}), t)
	fail(testEncodingString("Здравствуй, мир", []byte{156, 208, 151, 208, 180, 209, 128, 208, 176, 208, 178, 209, 129, 209, 130, 208, 178, 209, 131, 208, 185, 44, 32, 208, 188, 208, 184, 209, 128}), t)
	fail(testEncodingString("γειά σου κόσμος", []byte{156, 206, 179, 206, 181, 206, 185, 206, 172, 32, 207, 131, 206, 191, 207, 133, 32, 206, 186, 207, 140, 207, 131, 206, 188, 206, 191, 207, 130}), t)
	fail(testEncodingString("여보세요 세계", []byte{147, 236, 151, 172, 235, 179, 180, 236, 132, 184, 236, 154, 148, 32, 236, 132, 184, 234, 179, 132}), t)
	fail(testEncodingString("The quick brown fox jumps over the lazy dog.", []byte{172, 84, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103, 46}), t)
	fail(testEncodingString("中国很早就出现了专门用于启蒙的识字课本，秦代出现的有《倉颉篇》、《爰历篇》，汉代则有司马相如的《凡将篇》、贾鲂的《滂喜篇》、蔡邕的《劝学篇》、史游的《急就章》，三国时代有《埤苍》、《广苍》、《始学篇》等，这些作品中只有《急就章》对后世产生了影响，其余影响不大，《急就章》虽然是《苍颉篇》之后较突出的小学之书，但由于流传中出现了种种问题，其权威性到南北朝时已大不如前，而这一时期出现的一些启蒙读物如《庭诰》、《诂幼》之类，可读性有限。就是在这样的背景下，《千字文》问世了", []byte{55, 48, 50, 58, 228, 184, 173, 229, 155, 189, 229, 190, 136, 230, 151, 169, 229, 176, 177, 229, 135, 186, 231, 142, 176, 228, 186, 134, 228, 184, 147, 233, 151, 168, 231, 148, 168, 228, 186, 142, 229, 144, 175, 232, 146, 153, 231, 154, 132, 232, 175, 134, 229, 173, 151, 232, 175, 190, 230, 156, 172, 239, 188, 140, 231, 167, 166, 228, 187, 163, 229, 135, 186, 231, 142, 176, 231, 154, 132, 230, 156, 137, 227, 128, 138, 229, 128, 137, 233, 162, 137, 231, 175, 135, 227, 128, 139, 227, 128, 129, 227, 128, 138, 231, 136, 176, 229, 142, 134, 231, 175, 135, 227, 128, 139, 239, 188, 140, 230, 177, 137, 228, 187, 163, 229, 136, 153, 230, 156, 137, 229, 143, 184, 233, 169, 172, 231, 155, 184, 229, 166, 130, 231, 154, 132, 227, 128, 138, 229, 135, 161, 229, 176, 134, 231, 175, 135, 227, 128, 139, 227, 128, 129, 232, 180, 190, 233, 178, 130, 231, 154, 132, 227, 128, 138, 230, 187, 130, 229, 150, 156, 231, 175, 135, 227, 128, 139, 227, 128, 129, 232, 148, 161, 233, 130, 149, 231, 154, 132, 227, 128, 138, 229, 138, 157, 229, 173, 166, 231, 175, 135, 227, 128, 139, 227, 128, 129, 229, 143, 178, 230, 184, 184, 231, 154, 132, 227, 128, 138, 230, 128, 165, 229, 176, 177, 231, 171, 160, 227, 128, 139, 239, 188, 140, 228, 184, 137, 229, 155, 189, 230, 151, 182, 228, 187, 163, 230, 156, 137, 227, 128, 138, 229, 159, 164, 232, 139, 141, 227, 128, 139, 227, 128, 129, 227, 128, 138, 229, 185, 191, 232, 139, 141, 227, 128, 139, 227, 128, 129, 227, 128, 138, 229, 167, 139, 229, 173, 166, 231, 175, 135, 227, 128, 139, 231, 173, 137, 239, 188, 140, 232, 191, 153, 228, 186, 155, 228, 189, 156, 229, 147, 129, 228, 184, 173, 229, 143, 170, 230, 156, 137, 227, 128, 138, 230, 128, 165, 229, 176, 177, 231, 171, 160, 227, 128, 139, 229, 175, 185, 229, 144, 142, 228, 184, 150, 228, 186, 167, 231, 148, 159, 228, 186, 134, 229, 189, 177, 229, 147, 141, 239, 188, 140, 229, 133, 182, 228, 189, 153, 229, 189, 177, 229, 147, 141, 228, 184, 141, 229, 164, 167, 239, 188, 140, 227, 128, 138, 230, 128, 165, 229, 176, 177, 231, 171, 160, 227, 128, 139, 232, 153, 189, 231, 132, 182, 230, 152, 175, 227, 128, 138, 232, 139, 141, 233, 162, 137, 231, 175, 135, 227, 128, 139, 228, 185, 139, 229, 144, 142, 232, 190, 131, 231, 170, 129, 229, 135, 186, 231, 154, 132, 229, 176, 143, 229, 173, 166, 228, 185, 139, 228, 185, 166, 239, 188, 140, 228, 189, 134, 231, 148, 177, 228, 186, 142, 230, 181, 129, 228, 188, 160, 228, 184, 173, 229, 135, 186, 231, 142, 176, 228, 186, 134, 231, 167, 141, 231, 167, 141, 233, 151, 174, 233, 162, 152, 239, 188, 140, 229, 133, 182, 230, 157, 131, 229, 168, 129, 230, 128, 167, 229, 136, 176, 229, 141, 151, 229, 140, 151, 230, 156, 157, 230, 151, 182, 229, 183, 178, 229, 164, 167, 228, 184, 141, 229, 166, 130, 229, 137, 141, 239, 188, 140, 232, 128, 140, 232, 191, 153, 228, 184, 128, 230, 151, 182, 230, 156, 159, 229, 135, 186, 231, 142, 176, 231, 154, 132, 228, 184, 128, 228, 186, 155, 229, 144, 175, 232, 146, 153, 232, 175, 187, 231, 137, 169, 229, 166, 130, 227, 128, 138, 229, 186, 173, 232, 175, 176, 227, 128, 139, 227, 128, 129, 227, 128, 138, 232, 175, 130, 229, 185, 188, 227, 128, 139, 228, 185, 139, 231, 177, 187, 239, 188, 140, 229, 143, 175, 232, 175, 187, 230, 128, 167, 230, 156, 137, 233, 153, 144, 227, 128, 130, 229, 176, 177, 230, 152, 175, 229, 156, 168, 232, 191, 153, 230, 160, 183, 231, 154, 132, 232, 131, 140, 230, 153, 175, 228, 184, 139, 239, 188, 140, 227, 128, 138, 229, 141, 131, 229, 173, 151, 230, 150, 135, 227, 128, 139, 233, 151, 174, 228, 184, 150, 228, 186, 134}), t)
}

func TestEncodeSlice(t *testing.T) {
	t.Parallel()

	fail(testEncodingSlice([]int{1, 2, 3, 4}, []byte{196, 1, 2, 3, 4}), t)
	fail(testEncodingSlice([]int{1, 2, -256, 255}, []byte{196, 1, 2, 63, 255, 0, 63, 0, 255}), t)
	fail(testEncodingSlice(List{1, 2, 3, List{4.0, 5.0, 6.0, []string{"7", "8", "9"}}}, []byte{196, 1, 2, 3, 196, 66, 64, 128, 0, 0, 66, 64, 160, 0, 0, 66, 64, 192, 0, 0, 195, 129, 55, 129, 56, 129, 57}), t)
}

func TestEncodeMap(t *testing.T) {
	t.Parallel()

	fail(testEncodingMap(Dict{1: Dict{2: 3}}, []byte{103, 1, 103, 2, 3}), t)
	fail(testEncodingMap(Dict{1: 2, "abc": "def", 1.5: 3}, []byte{105, 1, 2, 66, 63, 192, 0, 0, 3, 131, 97, 98, 99, 131, 100, 101, 102}), t)
	fail(testEncodingMap(Dict{1: Dict{2: List{3, 4, 5, 6, Dict{7: "Eight!"}}}}, []byte{103, 1, 103, 2, 197, 3, 4, 5, 6, 103, 7, 134, 69, 105, 103, 104, 116, 33}), t)
}

func testEncodingInt(value int64, target []byte) error {
	encoded, err := encodeInt(value)
	if err != nil {
		return err
	}
	if !bytesEqual(encoded, target) {
		return fmt.Errorf("Integer %d was not encoded propperly: expected \"%v\", but got \"%v\"", value, target, encoded)
	}
	return nil
}

func testEncodingFloat(value float64, target []byte) error {
	encoded, err := encodeFloat(value)
	if err != nil {
		return err
	}
	if !bytesEqual(encoded, target) {
		return fmt.Errorf("Float %f was not encoded propperly: expected \"%v\", but got \"%v\"", value, target, encoded)
	}
	return nil
}

func testEncodingSlice(value interface{}, target []byte) error {
	encoded, err := encodeSlice(value)
	if err != nil {
		return err
	}
	if !bytesEqual(encoded, target) {
		return fmt.Errorf("Slice %v was not encoded propperly: expected \"%v\", but got \"%v\"", value, target, encoded)
	}
	return nil
}

func testEncodingString(value string, target []byte) error {
	encoded, err := encodeString(value)
	if err != nil {
		return err
	}
	if !bytesEqual(encoded, target) {
		return fmt.Errorf("String \"%s\" was not encoded propperly: expected \"%v\", but got \"%v\"", value, target, encoded)
	}
	return nil

}

func testEncodingMap(value interface{}, target []byte) error {
	encoded, err := encodeMap(value)
	if err != nil {
		return err
	}
	if !bytesEqual(encoded, target) {
		return fmt.Errorf("Map %v was not encoded propperly: expected \"%v\", but got \"%v\"", value, target, encoded)
	}
	return nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, val := range a {
		if val != b[i] {
			return false
		}
	}

	return true
}

func fail(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
