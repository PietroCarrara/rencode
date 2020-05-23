package rencode

import "testing"

func TestScanSlice(t *testing.T) {
	t.Parallel()

	data := List{
		List{1, 2, 3},
		List{"1", 2, 3.0},
		Dict{"hello": "world", "good": "morning"},
		Dict{1: "1", "2": 1},
		[]interface{}{1, 2, 3},
		1,
		nil,
	}

	var arr1 []int
	var arr2 List
	var map1 map[string]string
	var map2 Dict
	var arr3 []int
	var one int
	var last rune

	_, err := ScanSlice(data, &arr1, &arr2, &map1, &map2, &arr3, &one, &last)
	if err != nil {
		t.Error(err)
	}

	assertEquals(data[0], arr1, t)
	assertEquals(data[1], arr2, t)
	assertEquals(data[2], map1, t)
	assertEquals(data[3], map2, t)
	assertEquals(data[4], arr3, t)
	assertEquals(data[5], one, t)
	assertEquals(rune(0), last, t)
}

func TestScanMap(t *testing.T) {
	t.Parallel()

	data := Dict{
		"123":         List{1, 2, 3},
		"2nd element": List{"1", 2, 3.0},
		123:           Dict{"hello": "world", "good": "morning"},
		"array":       []interface{}{1, 2, 3},
		1.0:           1,
		"nil":         nil,
	}

	var elem1 []int
	var elem2 []interface{}
	var elem3 map[string]string
	var elem4 List
	var elem5 int
	var elem6 rune

	n, err := ScanMap(
		data,
		MapRef{
			Key: "123",
			Ref: &elem1,
		},
		MapRef{
			Key: "2nd element",
			Ref: &elem2,
		},
		MapRef{
			Key: 123,
			Ref: &elem3,
		},
		MapRef{
			Key: "array",
			Ref: &elem4,
		},
		MapRef{
			Key: 1.0,
			Ref: &elem5,
		},
		MapRef{
			Key: "nil",
			Ref: &elem6,
		},
	)

	if err != nil {
		t.Error(err)
	}

	if n != 6 {
		t.Errorf("expected 6 keys to be scanned, found %d", n)
	}

	assertEquals(elem1, data["123"], t)
	assertEquals(elem2, data["2nd element"], t)
	assertEquals(elem3, data[123], t)
	assertEquals(elem4, data["array"], t)
	assertEquals(elem5, data[1.0], t)
	assertEquals(elem6, rune(0), t)
}
