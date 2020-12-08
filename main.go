package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

/*
	Структура, представляющая собой один элемент(значение или параметр)
	из файла TestcaseStructure.json
	id - id структуры из файла
	idStrNum level title value values params:= Номера строк соответствующих параметров
 */
type Element struct{
	id int
	idStrNum int
	level int
	title int
	value int
	values int
	params int
}

/*
	Структура для представления правил замены из файла Values.json
	Две идентичные к файлу переменные
 */
type Values struct {
	id int
	str string
}


func writeError() {
	writeResultFile("./error.json", []string{
		"{",
		"\t\"error\": {",
		"\t\t\"message\": \"Входные файлы некорректны\"",
		"\t}",
		"}" })

	panic("Error while parse file!")
}

func validateElement(elem Element) {

	if (elem.id < 0) || (elem.title < 0) ||
		(elem.params > 0) && (elem.value + elem.values != -2) ||
		(elem.value > 0) && (elem.params > 0) {

		writeError()
	}

}

/*
	Функция чтения файла TestcaseStructure.json. Заполняет массив структур
	Element.
	path - путь к файлу
*/
func readElementsFile(path string) ([]string, []Element, int, int){
	file, _ := os.Open(path)
	var fileStrings []string
	var structs []Element

	strct:=false
	var tmp Element
	i:=0
	idsCnt:=0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str:=scanner.Text()
		fileStrings = append(fileStrings, str)

		if strings.Contains(str, "\"id\":") && strct == false{
			s := str[strings.Index(str, ":")+2:strings.Index(str, ",")]
			id, _ := strconv.Atoi(s)
			tmp = Element{ id, i, strings.Count(str, "\t"), -1, -1, -1, -1}
			strct = true
		} else if strings.Contains(str, "\"id\":") && strct == true {
			validateElement(tmp)
			structs=append(structs, tmp)
			s := str[strings.Index(str, ":")+2:strings.Index(str, ",")]
			id, _ := strconv.Atoi(s)
			tmp = Element{ id, i, strings.Count(str, "\t"), -1, -1, -1, -1}

			idsCnt++
		}


		if strct == true {
			if strings.Contains(str, "\"title\":") {
				tmp.title = i
			} else if strings.Contains(str, "\"value\":") {
				tmp.value = i
			} else if strings.Contains(str, "\"values\":") {
				tmp.values = i
			} else if strings.Contains(str, "\"params\":") {
				tmp.params = i
			} else {
				writeError()
			}
		}

		i++
	}
	structs=append(structs, tmp)
	idsCnt++

	return fileStrings, structs, i, idsCnt
}

/*
	Функция чтения файла Values.json. Заполняет массив структур
	Values.
	path - путь к файлу
*/
func readValuesFile(path string) ([]Values, int) {
	var toChange []Values
	changeFile, _ := os.Open(path)
	scanner := bufio.NewScanner(changeFile)

	var tcTmp Values
	j:=0
	strNum:=0
	idStrNum:=0

	for scanner.Scan() {
		str:=scanner.Text()

		if strings.Contains(str, "\"id\":") {
			s := str[strings.Index(str, ":")+2 : strings.Index(str, ",")]
			id, _ := strconv.Atoi(s)
			tcTmp.id = id
			idStrNum = strNum
		}else if strings.Contains(str, "\"value\":"){
			if idStrNum+1 != strNum {
				writeError()
			}
			s := str[strings.Index(str, ":")+2: len(str)]
			tcTmp.str = strings.TrimSpace(s)
			toChange=append(toChange, tcTmp)
			j++
		}else if idStrNum+1 != strNum {
			writeError()
		}
		strNum++
	}

	return toChange, j
}

/*
	Запись выходного файла StructureWithValues.json
	path - путь к выходному файлу
	fileStrings - построчное представление записываемого файла
*/
func writeResultFile(path string, fileStrings []string) {
	output, _ := os.Create(path)
	w := bufio.NewWriter(output)

	for _, line := range fileStrings {
		w.WriteString(line+"\n")
	}
	w.Flush()
}

/*
	Функция формирования измененной строки.
	(Замена value или title в строке)
	in - строка, подлежащая перезаписи
	concat - добавляемое значение
 */
func formNewStr(in string, concat string) string {
	return in[0:strings.Index(in, ":")+2]+concat
}

/*
	Функция замены элемента title или value
	a - значение заменяемого элемента
	b - значение замещающего элемента
	fileStrings - строки выходного файла
 */
func setString(a Element, b Element, fileStrings []string) {
	str:= ""
	if b.value > 0 {
		str = fileStrings[b.title][strings.Index(fileStrings[b.value], ":")+2:len(fileStrings[b.value])]
	} else if b.title > 0 {
		str = fileStrings[b.title][strings.Index(fileStrings[b.title], ":")+2:len(fileStrings[b.title])]
	}

	if a.value > 0 {
		fileStrings[a.value]=formNewStr(fileStrings[a.value], str)
	}else if a.title > 0 {
		fileStrings[a.title]=formNewStr(fileStrings[a.title], str)
	}
}

func main() {
	fileStrings, structs, _, elementsCount := readElementsFile("./TestcaseStructure.json")
	toChange, rulesCount:= readValuesFile("./values.json")

	for a := 0; a < elementsCount; a++ {
		for b := 0; b < rulesCount; b++ {
			if structs[a].id == toChange[b].id {
				if structs[a].values < 0 && structs[a].params < 0 && structs[a].value > 0 {
					fileStrings[structs[a].value]=formNewStr(fileStrings[structs[a].value], toChange[b].str)
				} else if structs[a].values > 0 || structs[a].params > 0 {
					for c:= a+1; (structs[a].level != (structs[c].level)) && c < elementsCount; c++ {
						num,_:=strconv.Atoi(toChange[b].str)
						if structs[c].id == num {
							setString(structs[a], structs[c], fileStrings)
							break
						}
					}
				}
			}
		}
	}

	writeResultFile("./StructureWithValues.json", fileStrings)
}