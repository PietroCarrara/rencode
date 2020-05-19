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
	fail(testEncodingInt(255, []byte{63, 0, 255}), t)
	fail(testEncodingInt(256, []byte{63, 1, 0}), t)
	fail(testEncodingInt(-256, []byte{63, 255, 0}), t)
	fail(testEncodingInt(math.MaxInt64, []byte{65, 127, 255, 255, 255, 255, 255, 255, 255}), t)
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
