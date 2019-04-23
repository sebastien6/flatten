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

func Test_FlattenList(t *testing.T) {
	expectedinterface := []interface{}{5, "hello", 7.6}
	got := FlattenList("[5, [hello, [7.6]]]")
	for i, v := range got.([]interface{}) {
		if v != expectedinterface[i] {
			t.Error("FlattenList incorrect returned list")
		}
	}

	expectedint := []int{5, 8, 2, 1, 3}
	got = FlattenList("[5, [8], [2, [1, 3]]]")
	for i, v := range got.([]int) {
		if v != expectedint[i] {
			t.Error("FlattenList incorrect returned list")
		}
	}

	expectedfloat := []float64{5.1, 8.6, 2.4, 1.0, 3.9}
	got = FlattenList("[5.1, [8.6], [2.4, [1.0, 3.9]]]")
	for i, v := range got.([]float64) {
		if v != expectedfloat[i] {
			t.Error("FlattenList incorrect returned list")
		}
	}

	expectedstring := []string{"hello", "goodbye"}
	got = FlattenList("[hello, [goodbye]]")
	for i, v := range got.([]string) {
		if v != expectedstring[i] {
			t.Error("FlattenList incorrect returned list")
		}
	}

	got = FlattenList([]byte("[5, [8], [2, [1, 3]]]"))
	for i, v := range got.([]int) {
		if v != expectedint[i] {
			t.Error("FlattenList incorrect returned list")
		}
	}
}

// func Benchmarkisintorfloat(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		isintorfloat([]byte("79458124545.1548545"))
// 	}
// }
