// Package flatten convert embedded array to a flat array
//
// the package export a single function (FlattenList) that take
// an interface as argument and return the resulting flatten
// array through an interface.
package flatten

import (
	"strconv"
)

// flatlist internal struct used to keep track of some metrics
// during the computation. Those metric will be used to define
// if a resulting array contain elements of the same type or
// a mix of various type.
type flatList struct {
	inputList   interface{}
	returnList  interface{}
	trimedList  []interface{}
	countInt    int
	countString int
	countFloat  int
}

// FlattenList take as input embedded array in three formats:
//
// string: Embedded array inside a string "[1, [2, 3], 4]"
// []byte: Such as JSON embedded array []byte("[1, [2, 3], 4]")
// []interface{}: Native Golang embedded array format
//
// The embedded array is first flatten, and then if all the
// element of the flatten array are of the same type, the
// resulting array is converted in the appropriate array
// format.
//
// the function can return the flatten array in 4 differents
// format based on its content:
//
// []int: 			All elements are integer, not exceeding the limit
// 					relative to int32
// []float64: 		All elements are float value, not exceeding the limit
// 					relative to float32
// []string:    	All elements are of type string or interpreted as string
// []interface{}:	Element of the array are a mix of different types (integer,
//					float and string)
func FlattenList(input interface{}) interface{} {
	var fl flatList
	fl.inputList = input

	switch v := fl.inputList.(type) {
	case string, []byte:
		fl.trimList()
		fl.parseList()
	case []interface{}:
		fl.trimedList = fl.parseInterfaceArray(&v)
	}

	if fl.countString == 0 && fl.countFloat == 0 {
		l := fl.convertInt()
		return l
	} else if fl.countInt == 0 && fl.countFloat == 0 {
		l := fl.convertString()
		return l
	} else if fl.countInt == 0 && fl.countString == 0 {
		l := fl.convertFloat64()
		return l
	}

	return fl.trimedList
}

// convertInt convert an array of interfaces to an array
// of Int
func (f *flatList) convertInt() (convertedList []int) {
	for _, x := range f.trimedList {
		convertedList = append(convertedList, x.(int))
	}
	return
}

// convert64 convert an array of interfaces to an array
// of float64
func (f *flatList) convertFloat64() (convertedList []float64) {
	for _, x := range f.trimedList {
		convertedList = append(convertedList, x.(float64))
	}
	return
}

// convertString convert an array of interfaces to an array
// of string
func (f *flatList) convertString() (convertedList []string) {
	for _, x := range f.trimedList {
		convertedList = append(convertedList, x.(string))
	}
	return
}

// trimList analyse the string to extract the values contained
// in it as strings.
func (f *flatList) trimList() {
	t := ""

	switch v := f.inputList.(type) {
	case string:
		a := []byte(v)
		for i, r := range a {
			if (r >= 48 && r <= 57) || (r >= 65 && r <= 90) || (r >= 97 && r <= 122) || r == '.' || r == '\'' || r == '-' || r == '@' || r == ':' || r == '/' {
				t = t + string(v[i])
			} else {
				if len(t) > 0 {
					f.trimedList = append(f.trimedList, t)
					t = ""
				}
			}
		}
	case []byte:
		for i, r := range v {
			if (r >= 48 && r <= 57) || (r >= 65 && r <= 90) || (r >= 97 && r <= 122) || r == '.' || r == '\'' || r == '-' || r == '@' || r == ':' || r == '/' {
				t = t + string(v[i])
			} else {
				if len(t) > 0 {
					f.trimedList = append(f.trimedList, t)
					t = ""
				}
			}
		}
	}
}

// parseList loop through the list of values to check their type
// and convert them accordingly.
func (f *flatList) parseList() {
	for i, v := range f.trimedList {
		isint, isfloat := isintorfloat([]byte(v.(string)))
		if isint {
			if isfloat {
				x, _ := strconv.ParseFloat(v.(string), 64)
				f.countFloat++
				f.trimedList[i] = x
			} else {
				x, _ := strconv.Atoi(v.(string))
				f.countInt++
				f.trimedList[i] = x
			}
		} else {
			f.countString++
		}
	}
}

// parseInterfaceArray loop recusively through embedded interface array
// to return a flatten array of elements.
// The function accept elements of type int16, int32, float32, float64,
// and string. Element of other type will just be discard from the list.
func (f *flatList) parseInterfaceArray(input *[]interface{}) (result []interface{}) {
	for _, v := range *input {
		switch x := v.(type) {
		case []interface{}:
			l := f.parseInterfaceArray(&x)
			for _, t := range l {
				result = append(result, t)
			}
		case int, int16, int32:
			result = append(result, x)
			f.countInt++
		case float32, float64:
			result = append(result, x)
			f.countFloat++
		case string:
			result = append(result, x)
			f.countString++
		}
	}

	return
}

// isintorfloat test if a value is a numerical value, and if
// the numerical value is of type float. countdot is used to
// prevent false positive with float, such as with string
// dates (12.02.2019).
func isintorfloat(b []byte) (bool, bool) {
	numerical := 0
	isfloat := false
	isnumerical := false
	countdot := 0

	for _, r := range b {
		if (r >= 48 && r <= 57) || r == '.' {
			numerical++
			if r == '.' {
				countdot++
			}
		}
	}

	if len(b) == numerical {
		if countdot <= 1 {
			isnumerical = true
		}

		if countdot == 1 {
			isfloat = true
		}
	}

	return isnumerical, isfloat
}
