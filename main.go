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
/*
func readData(in *Element, str string, arr []string, len int) *Element{
	for i:=0; i < len; i++ {
		if strings.Contains(str, "\"title\":") {
			tmp.title = i
		}
	}
	if strings.Contains(str, "\"title\":") {
		tmp.title = i
	}
	return in
}
*/
func main() {
	file, _ := os.Open("D:/go.json")
	var fileStrings []string
	var structs []Element
	//structs := list.New()

	//toChange := list.New()
	var toChange []Values
	changeFile, _ := os.Open("D:/values.json")
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

	/*
	for e := toChange.Front(); e != nil; e = e.Next() {
		tmp:= Values(e.Value.(Values))
		fmt.Println(tmp)
	}
*/
	//el := Values{34, "", 298}
	//toChange.PushBack(el)
	//el := Values{34, "", 298}

	strct := false
	var tmp Element
	i:=0
	idsCnt:=0

	scanner = bufio.NewScanner(file)
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

	for a := 0; a < idsCnt-1; a++ {
		for b := 0; b < j; b++ {
			if structs[a].id == toChange[b].id {
				if structs[a].values < 0 && structs[a].params < 0 && structs[a].value > 0 {
					fileStrings[structs[a].value]=fileStrings[structs[a].value][0:strings.Index(fileStrings[structs[a].value], ":")+2]+toChange[b].str
				} else if structs[a].values > 0 || structs[a].params > 0 {
					for c:= a+1; (structs[a].level != (structs[c].level)) && c < idsCnt; c++ {
						num,_:=strconv.Atoi(toChange[b].str)
						if structs[c].id == num {
							str:= ""
							if structs[c].title > 0 {
								str = fileStrings[structs[c].title][strings.Index(fileStrings[structs[a].value], ":")+2:len(fileStrings[structs[c].title])]
							}
							if structs[a].value > 0 {
								fileStrings[structs[a].value]=fileStrings[structs[a].value][0:strings.Index(fileStrings[structs[a].value], ":")+2]+str
								break
							}else if structs[a].title > 0 {
								fileStrings[structs[a].title]=fileStrings[structs[a].title][0:strings.Index(fileStrings[structs[a].title], ":")+2]+str
								break
							}
						}
					}
				}
			}
		}
		//fmt.Println(tmp)
	}

	output, _ := os.Create("D:/StructureWithValues.json")
	w := bufio.NewWriter(output)

	for _, line := range fileStrings {
		w.WriteString(line+"\n")
	}
	w.Flush()
/*
	for e := structs.Front(); e != nil; e = e.Next() {
		tmp:= Element(e.Value.(Element))
		fmt.Println(tmp)
	}
*/


	fmt.Print("Hi")


	//Simple Employee JSON object which we will parse
	/*

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	coronaVirusJSON := `{
        "name" : "covid-11",
        "country" : "China",
        "city" : "Wuhan",
        "reason" : "Non vedge Food",
		"arr" : [
					{"hi": "hi1"},
					{"hi": "hi2"}
				]
    }`
	jq := New().FromString(coronaVirusJSON)
	res:=jq.From("arr").Where("hi", "=", "hi2").Get()

	fmt.Println(res)

	res=jq.From("arr").Where("hi", "=", "1337")

	fmt.Println(res)
	*/
	/*
		// Declared an empty map interface
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(coronaVirusJSON), &result)

		// Print the data type of result variable
		fmt.Println(reflect.TypeOf(result))

		// Reading each value by its key
		fmt.Println("Name :", result["name"],
			"\nCountry :", result["country"],
			"\nCity :", result["city"],
			"\nReason :", result["reason"],
			"\nReason :", result["arr"]["hi"])

	*/
}