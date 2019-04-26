package flatten

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_isintorfloat(t *testing.T) {
	t.Log("Start tests function isintorfloat")
	input := []string{"12.46", "12", ".6f", "12.04.2019", "$12", "hello", "hello.com"}
	resint := []bool{true, true, false, false, false, false, false}
	resfloat := []bool{true, false, false, false, false, false, false}

	for i, v := range input {
		t.Logf("test input %s", v)
		isint, isfloat := isintorfloat([]byte(v))
		if isint != resint[i] && isfloat != resfloat[i] {
			t.Logf("    - ERROR\n")
			t.Errorf("\nisint = %t; want %t - isfloat = %t; want %t - ", isint, resint[i], isfloat, resfloat[i])
		}
		t.Logf("    - OK\n")
	}
	t.Logf("tests function isintorfloat completed\n")
}

func Test_parseList(t *testing.T) {
	var fl flatList
	fl.trimedList = []interface{}{"hello", "12", "1.0"}
	fl.parseList()

	for i, v := range fl.trimedList {
		x := reflect.TypeOf(v).Kind()
		if i == 0 && x != reflect.String {
			t.Errorf("%s not correctly evaluated as string", v)
			t.Logf("%T", fl.trimedList[0])
		} else if i == 1 && x != reflect.Int {
			t.Errorf("%s not correctly evaluated as int", v)
			t.Logf("%T", fl.trimedList[0])
		} else if i == 2 && x != reflect.Float64 {
			t.Errorf("%s not correctly evaluated as float64", v)
			t.Logf("%T", fl.trimedList[0])
		}

	}
}

func Test_convertInt(t *testing.T) {
	var fl flatList
	fl.trimedList = []interface{}{int(4)}
	l := fl.convertInt()
	r := fmt.Sprintf("%T", l)
	if r != "[]int" {
		t.Errorf("convertion failure, type is not []int")
	}
}

func Test_convertFloat64(t *testing.T) {
	var fl flatList
	fl.trimedList = []interface{}{float64(4.1)}
	l := fl.convertFloat64()
	r := fmt.Sprintf("%T", l)
	if r != "[]float64" {
		t.Errorf("convertion failure, type is not []float64")
	}
}

func Test_convertString(t *testing.T) {
	var fl flatList
	fl.trimedList = []interface{}{"hello"}
	l := fl.convertString()
	r := fmt.Sprintf("%T", l)
	t.Log(r)
	if r != "[]string" {
		t.Errorf("convertion failure, type is not []string")
	}
}

func Test_trimList(t *testing.T) {
	expected := []string{"1", "hello"}
	unexpected := []string{"1h", "ello"}
	var fl flatList
	fl.inputList = "[1, [hello]]"
	fl.trimList()
	if fl.trimedList == nil {
		t.Error("trimedList is empty")
	} else if len(fl.trimedList) != 2 {
		t.Error("trimedList contain an incorrect number of values")
	}
	for i, v := range fl.trimedList {
		if v != expected[i] {
			t.Error("trimedList incorrect values")
		}
	}
	for i, v := range fl.trimedList {
		if v == unexpected[i] {
			t.Error("trimedList incorrect values")
		}
	}
}

func Test_parseInterfaceArray_mixed(t *testing.T) {
	var fl flatList
	expectedinterface := []interface{}{5, 7.1, 8, "hello", 11, 7, 1}
	got := fl.parseInterfaceArray(&[]interface{}{5, []interface{}{7.1, 8}, []interface{}{"hello", 11, []interface{}{7}}, 1})

	for i := range got {
		if got[i] != expectedinterface[i] {
			t.Error("FlattenList incorrect returned list")
		}
	}
}

func Test_FlattenList_int(t *testing.T) {
	expected := []int{5, 7, 8, 0, 11, 7, 1}

	t.Log("Test FlattenList() - []interface, all int")
	input1 := []interface{}{5, []interface{}{7, 8}, []interface{}{0, 11, []interface{}{7}}, 1}
	got1 := FlattenList(input1)
	t.Log("Test FlattenList() ")
	switch v := got1.(type) {
	case []int:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %d:received %d", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []int")
	}

	t.Log("Test FlattenList() - string, all int")
	input2 := "[5, [7, [8], 0, 11], 7, [1]]"
	got2 := FlattenList(input2)
	switch v := got2.(type) {
	case []int:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %d:received %d", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []int")
	}

	t.Log("Test FlattenList() - []byte, all int")
	input3 := []byte("[5, [7, [8], 0, 11], 7, [1]]")
	got3 := FlattenList(input3)
	switch v := got3.(type) {
	case []int:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %d:received %d", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []int")
	}
}

func Test_FlattenList_float64(t *testing.T) {
	expected := []float64{1.6, 22.789, 3.1, 4.541, 8.96}

	t.Log("Test FlattenList() - []interface, all float")
	input1 := []interface{}{1.6, []interface{}{22.789, 3.1}, []interface{}{4.541, 8.96}}
	got1 := FlattenList(input1)
	switch v := got1.(type) {
	case []float64:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %f:received %f", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []float64")
	}

	t.Log("Test FlattenList() - string, all float")
	input2 := []byte("[1.6, 22.789, [3.1, [4.541], []], 8.96]")
	got2 := FlattenList(input2)
	switch v := got2.(type) {
	case []float64:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %f:received %f", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []float64")
	}

	t.Log("Test FlattenList() - []byte, all float")
	input3 := []byte("[1.6, 22.789, [3.1, [4.541], []], 8.96]")
	got3 := FlattenList(input3)
	switch v := got3.(type) {
	case []float64:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %f:received %f", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []float64")
	}
}

func Test_FlattenList_string(t *testing.T) {
	expected := []string{"hello", "thanks", "2018.7.5"}
	input1 := []interface{}{"hello", []interface{}{"thanks", "2018.7.5"}}
	got1 := FlattenList(input1)

	t.Log("Test FlattenList() - []interface, all string")
	switch v := got1.(type) {
	case []string:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %s:received %s", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []string")
	}

	t.Log("Test FlattenList() - string, all string")
	input2 := "[hello, [thanks, [2018.7.5]]]"
	got2 := FlattenList(input2)

	switch v := got2.(type) {
	case []string:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %s:received %s", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []string")
	}

	t.Log("Test FlattenList() - []byte, all string")
	input3 := []byte("[hello, [thanks, [2018.7.5]]]")
	got3 := FlattenList(input3)

	switch v := got3.(type) {
	case []string:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %s:received %s", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []string")
	}
}

func Test_FlattenList_interface(t *testing.T) {
	var expected []interface{}
	expected = append(expected, 5.1)
	expected = append(expected, "hello")
	expected = append(expected, "bye")
	expected = append(expected, 3)
	expected = append(expected, 6)

	t.Log("Test FlattenList() - []interface, mixed")
	input1 := []interface{}{5.1, []interface{}{"hello", "bye"}, []interface{}{3, 6}}
	got1 := FlattenList(input1)
	switch v := got1.(type) {
	case []interface{}:
		t.Log("good")
		t.Log(v)
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %f:received %f", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []interface{}")
	}

	t.Log("Test FlattenList() - string, mixed")
	input2 := "[5.1, [hello, [bye]], [3], 6]"
	got2 := FlattenList(input2)
	switch v := got2.(type) {
	case []interface{}:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %f:received %f", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []float64")
	}

	t.Log("Test FlattenList() - []byte, mixed")
	input3 := []byte("[5.1, [hello, [bye]], [3], 6]")
	got3 := FlattenList(input3)
	switch v := got3.(type) {
	case []interface{}:
		for i := range v {
			if v[i] != expected[i] {
				t.Errorf("FlattenList incorrect returned list. Expected %f:received %f", expected[i], v[i])
			}
		}
	default:
		t.Error("FlattenList is not of type []float64")
	}
}

func Benchmark_parseInterfaceArray(b *testing.B) {
	var fl flatList

	for i := 0; i < b.N; i++ {
		fl.parseInterfaceArray(&[]interface{}{5, []interface{}{7.1, 8}, []interface{}{"hello", 11, []interface{}{7}}, 1})
	}
}
