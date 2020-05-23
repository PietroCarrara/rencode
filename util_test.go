package rencode

import (
	"fmt"
	"testing"
)

func encode(data interface{}, t *testing.T) []byte {
	bytes, err := Encode(data)

	if err != nil {
		t.Error(err)
	}

	return bytes
}

func decode(bytes []byte, t *testing.T) interface{} {
	var x interface{}
	n, err := Decode(bytes, &x)

	if err != nil {
		t.Error(err)
	}

	if n != len(bytes) {
		t.Errorf("used only %d bytes out of %d while decoding", n, len(bytes))
	}

	return x
}

func assertEquals(a, b interface{}, t *testing.T) {
	if !objEquals(a, b) {
		t.Errorf("expected \"%v\", but found \"%v\"", a, b)
	}
}

// TODO: Be smart
// Really poor object comparison, but gets the job done
func objEquals(a, b interface{}) bool {
	str1 := fmt.Sprint(a)
	str2 := fmt.Sprint(b)

	return str1 == str2
}
