package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Element struct{
	id int
	idStrNum int
	level int
	title int
	value int
	values int
	params int
}

type Values struct {
	id int
	str string
}

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
			if (tmp.params > 0 && tmp.title < 0) && tmp.value < 0{
				fmt.Print("Error\n")
			} else {
				structs=append(structs, tmp)
				s := str[strings.Index(str, ":")+2:strings.Index(str, ",")]
				id, _ := strconv.Atoi(s)
				tmp = Element{ id, i, strings.Count(str, "\t"), -1, -1, -1, -1}

				idsCnt++
			}
		}


		if strct == true {
			if strings.Contains(str, "\"title\":") {
				tmp.title = i
			}

			if strings.Contains(str, "\"value\":") {
				tmp.value = i
			}

			if strings.Contains(str, "\"values\":") {
				tmp.values = i
			}

			if strings.Contains(str, "\"params\":") {
				tmp.params = i
			}
		}

		i++
	}
	structs=append(structs, tmp)
	idsCnt++
	i++

	return fileStrings, structs, i, idsCnt
}

func readValuesFile(path string) ([]Values, int) {
	var toChange []Values
	changeFile, _ := os.Open(path)
	scanner := bufio.NewScanner(changeFile)

	var tcTmp Values
	j:=0

	for scanner.Scan() {
		str:=scanner.Text()

		if strings.Contains(str, "\"id\":") {
			s := str[strings.Index(str, ":")+2:strings.Index(str, ",")]
			id, _ := strconv.Atoi(s)
			tcTmp.id = id
		} else if strings.Contains(str, "\"value\":") {
			s := str[strings.Index(str, ":")+2: len(str)]
			tcTmp.str = s
			toChange=append(toChange, tcTmp)
			j++
		}
	}

	return toChange, j
}

func writeResultFile(path string, fileStrings []string) {
	output, _ := os.Create(path)
	w := bufio.NewWriter(output)

	for _, line := range fileStrings {
		w.WriteString(line+"\n")
	}
	w.Flush()
}

func formNewStr(in string, concat string) string {
	return in[0:strings.Index(in, ":")+2]+concat
}

func setString(a Element, b Element, fileStrings []string) {
	str:= ""
	if b.title > 0 {
		str = fileStrings[b.title][strings.Index(fileStrings[a.value], ":")+2:len(fileStrings[b.title])]
	}
	if a.value > 0 {
		fileStrings[a.value]=formNewStr(fileStrings[a.value], str)
	}else if a.title > 0 {
		fileStrings[a.title]=formNewStr(fileStrings[a.title], str)
	}
}

func main() {
	fileStrings, structs, _, idsCnt := readElementsFile("D:/go.json")
	toChange, j:= readValuesFile("D:/values.json")

	for a := 0; a < idsCnt-1; a++ {
		for b := 0; b < j; b++ {
			if structs[a].id == toChange[b].id {
				if structs[a].values < 0 && structs[a].params < 0 && structs[a].value > 0 {
					fileStrings[structs[a].value]=formNewStr(fileStrings[structs[a].value], toChange[b].str)
				} else if structs[a].values > 0 || structs[a].params > 0 {
					for c:= a+1; (structs[a].level != (structs[c].level)) && c < idsCnt; c++ {
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

	writeResultFile("D:/StructureWithValues.json", fileStrings)

	fmt.Print("Hi")

}