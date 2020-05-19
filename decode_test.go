package rencode

import "testing"

func TestDecodingBool(t *testing.T) {
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
