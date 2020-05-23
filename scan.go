package rencode

import (
	"fmt"
	"reflect"
)

// MapRef is used to reference a value in a map
// Key is the key to that value, and Ref is a pointer
// to where that value is (or should be) in memory
type MapRef struct {
	Key interface{}
	Ref interface{}
}

// ScanSlice iterates through a array and copies the values
// to each provided pointer.
// Nil references to be filled are skipped.
// If there are more values to be filled than the length of the array,
// they are not filled.
// ScanSlice returns the number of elements that have been filled
// (those skipped because of nil do not count)
func ScanSlice(data interface{}, args ...interface{}) (int, error) {

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	if dataType.Kind() != reflect.Array && dataType.Kind() != reflect.Slice {
		return 0, fmt.Errorf("can't use type \"%s\" as slice", dataType)
	}

	length := dataValue.Len()

	n := 0
	for i := 0; i < length && i < len(args); i++ {
		if args[i] == nil {
			continue
		}

		ref := reflect.ValueOf(args[i])
		val := dataValue.Index(i).Elem()

		err := scanSingle(val, ref)
		if err != nil {
			return n, err
		}

		n++
	}

	return n, nil
}

// ScanMap iterates through a map and copies the values
// of each key to the matching MapRef
// Returns the number of keys filled
func ScanMap(data interface{}, args ...MapRef) (int, error) {

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	if dataType.Kind() != reflect.Map {
		return 0, fmt.Errorf("can't use type \"%s\" as map", dataType)
	}

	n := 0
	iter := dataValue.MapRange()

	for iter.Next() {
		arg := mapRefKey(iter.Key().Interface(), args)

		// User didn't want this key
		if arg == nil {
			continue
		}

		ref := reflect.ValueOf(arg.Ref)
		val := iter.Value()

		err := scanSingle(val, ref)
		if err != nil {
			return n, err
		}

		n++
	}

	return n, nil
}

// Scans a single value
func scanSingle(val, ref reflect.Value) error {

	if val.Kind() == reflect.Interface {
		return scanSingle(val.Elem(), ref)
	}

	if ref.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("type \"%s\" is not a reference type", ref.Type())
	}

	ref = ref.Elem()

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		// Let's scan each value of the array into a new array,
		// and assign it to the reference we've received
		length := val.Len()
		values := reflect.MakeSlice(ref.Type(), length, length)

		// Pointer to each index of the values array
		refs := make([]interface{}, length)
		for j := 0; j < length; j++ {
			refs[j] = values.Index(j).Addr().Interface()
		}

		// Scan the values array using the references stored in refs
		ScanSlice(val.Interface(), refs...)

		// Finally set the array
		ref.Set(values)
	} else if val.Kind() == reflect.Map {

		// Allocate a new map
		result := reflect.MakeMapWithSize(ref.Type(), val.Len())

		refs := make([]MapRef, 0)
		valIter := val.MapRange()

		// Create a reference for each key
		for valIter.Next() {
			var data interface{}

			refs = append(refs, MapRef{
				Key: valIter.Key().Interface(),
				Ref: &data,
			})
		}

		// Fill the references
		_, err := ScanMap(val.Interface(), refs...)
		if err != nil {
			return err
		}

		// Store the decoded values onto the new map
		for _, index := range refs {
			key := reflect.ValueOf(index.Key)
			value := reflect.ValueOf(index.Ref).Elem()

			for key.Kind() == reflect.Interface {
				key = key.Elem()
			}
			for value.Kind() == reflect.Interface {
				value = value.Elem()
			}

			if key.IsValid() {
				key = key.Convert(ref.Type().Key())
			}
			if value.IsValid() {
				value = value.Convert(ref.Type().Elem())
			}

			result.SetMapIndex(key, value)
		}

		// Point the map to the newly created
		ref.Set(result)
	} else {
		if val.IsValid() {
			ref.Set(val.Convert(ref.Type()))
		}
	}

	return nil
}

// Finds the reference that has the supplied key
func mapRefKey(key interface{}, args []MapRef) *MapRef {

	for _, v := range args {
		if reflect.DeepEqual(key, v.Key) {
			return &v
		}
	}

	return nil
}
