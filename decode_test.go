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
