// Package flatten take a list of embbedded array into a string or []byte
// from and return a flatten list based on content.
//
// if the list contains a single type of values (int, float64, or string),
// it returns an equivalent gl array ([]int, []float64, []string).
//
// if the lsit is made of mutliple type of values, it returns a
// Go interface array ([]interface{}).
package flatten

import (
	"strconv"
)

type flatList struct {
	inputList   interface{}
	returnList  interface{}
	trimedList  []interface{}
	countInt    int
	countString int
	countFloat  int
}

// FlattenList takes a list of embedded array in a string, and return
// the appropriate flatten list of values.
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
