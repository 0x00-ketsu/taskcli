package utils

import "reflect"

// RemoveEmptyString removes empty string from string slice
func RemoveEmptyString(s []string) []string {
	var target []string
	for _, str := range s {
		if str != "" {
			target = append(target, str)
		}
	}

	return target
}

// InArray checks is val exists in an Array(Slice)
func InArray(val interface{}, array interface{}) bool {
	return AtArrayPosition(val, array) != -1
}

// AtArrayPosition finds the position(int) of val in an Array(Slice)
func AtArrayPosition(val interface{}, array interface{}) (index int) {
	index = -1

	if reflect.TypeOf(array).Kind() == reflect.Slice {
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i

				return
			}
		}
	}

	return
}
