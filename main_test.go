package main

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	main()
	m.Run()

	strsTest1, _, size1, _:=readElementsFile("./test_json/StructureWithValuesTest.json")
	strsTest2, _, size2, _:=readElementsFile("./StructureWithValues.json")
	if size1 != size2 {
		panic("Files not equals")
	}

	for i:=0;i<size1;i++{
		if strings.Compare(strsTest1[i], strsTest2[i]) != 0 {
			panic("Files not equals")
		}
	}

}

func TestValidateElement(t *testing.T) {
	elem:= Element{1, 1, 1, 1, -1, -1, -1}
	validateElement(elem)

	elem= Element{1, 1, 1, 1, 1, 1, -1}
	validateElement(elem)

	elem= Element{1, 1, 1, 1, -1, -1, 1}
	validateElement(elem)
}

func TestReadElementsFile(t *testing.T) {
	_, elems, _, _ := readElementsFile("./test_json/test2.json")

	if elems[0].id != 73 || elems[0].title != 3 ||
		elems[0].value != -1 || elems[0].idStrNum != 2 ||
		elems[0].params != 4 || elems[0].values != -1 ||
		elems[0].level != 2{
		t.Errorf("Error while read fist element from test2.json")
	}

	if elems[1].id != 345 || elems[1].title != 6 ||
		elems[1].value != -1 || elems[1].idStrNum != 5 ||
		elems[1].params != -1 || elems[1].values != -1 ||
		elems[1].level != 3{
		t.Errorf("Error while read second element from test2.json")
	}

}

func TestReadValuesFile(t *testing.T) {
	vals, _:=readValuesFile("./test_json/testValues.json")

	if vals[0].id != 34 || vals[0].str != "298" {
		t.Errorf("Error while read first value")
	}

	if vals[1].id != 146 || vals[1].str != "\"Валидация\"" {
		t.Errorf("Error while read second value")
	}
}

func TestWriteResultFile(t *testing.T) {
	original, _, cnt1, _ := readElementsFile("./test_json/test1.json")

	writeResultFile("./test_json/tmp.json", original)

	strs, _, cnt2, _ := readElementsFile("./test_json/test1.json")

	if cnt1 != cnt2 {
		t.Errorf("String numbers not equals: %d != %d", cnt1, cnt2)
	}

	for i:=0; i < cnt1; i++ {
		if strings.Compare(original[i], strs[i]) != 0 {
			t.Errorf("Write func works incorrect: %s != %s", original[i], strs[i])
		}

	}

	os.Remove("./test_json/tmp.json")
}

func TestFormNewStr(t *testing.T) {
	str:="\"value\": \"old\""
	str=formNewStr(str, "\"new\"", 0, []string{"\"value\": \"old\"", "}"})
	if strings.Contains(str, "\"value\": \"new\"") == false {
		t.Errorf("test_formNewStr failed! %s != %s", str, "\"value\": \"new\"")
	}
}

func TestSetString(t *testing.T) {
	strs, elems, _, _ := readElementsFile("./test_json/test1.json")

	setString(elems[0], elems[1], strs)
	if strings.Contains(strs[elems[0].value], "SellerX") == false {
		t.Errorf("test_setString failed! (set value)")
	}

	strs, elems, _, _ = readElementsFile("./test_json/test2.json")

	setString(elems[0], elems[1], strs)
	if strings.Contains(strs[elems[0].title], "SellerX") == false {
		t.Errorf("test_setString failed! (set title)")
	}
}