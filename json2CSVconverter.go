package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Json struct {
	Key   string      // json key
	Value interface{} //json element
}

func getJSONfile(path string) []byte {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	return jsonFile
}

type KeyValue map[string]interface{}

func parseJSONfile(jsonFile []byte) {

	var jVector Json
	json.Unmarshal(jsonFile, &jVector.Value)

	f := make(KeyValue, 0)

	cnt := 0

	err := recrusJson(jVector, cnt, f)
	if err != nil {
		fmt.Println("err")
	}

	keys := make([]string, 0, len(f))

	for k := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		// fmt.Println(k, f[k])
		fmt.Println("CSV Key: ", k)
		fmt.Println("    Value: ", f[k])
		fmt.Println()
	}

	// flatting extracted json data to csv data
	flatMap := make(map[string][]interface{})

	for _, k := range keys {

		fmt.Println("KEY is ", k)
		// "[-]?\d[\d,]*[\.]?[\d{2}]*" for any number
		re := regexp.MustCompile(`'[0-9]+'`)

		t := strings.Split(k, "\t") // "\t" as csv key seperater
		var pos []interface{}
		// find if there is array index num
		if re.MatchString(k) {

			for i := len(t) - 1; i >= 0; i-- {
				// if has array index
				if re.MatchString(t[i]) {
					// csv array position ----------- this must be first, cause key will change
					t[i] = strings.TrimPrefix(t[i], "'")
					t[i] = strings.TrimSuffix(t[i], "'")
					fmt.Println("you'r arr index :  ", t[i])
					position, err := strconv.Atoi(t[i])
					if err != nil {
						fmt.Println("error: ", err)
					}

					// original key is k

					// ---------changing key
					t = append(t[:i], t[i+1:]...) // take out t[i]
					// Create csv key and assign value, value must insert in the position index
					if _, ok := flatMap[strings.Join(t, "_")]; !ok {
						pos = make([]interface{}, position+1)
						pos[position] = f[k]
						flatMap[strings.Join(t, "_")] = pos
					} else {
						// element key exist
						l := position - len(flatMap[strings.Join(t, "_")])
						if l >= 0 {
							for l > 0 {
								flatMap[strings.Join(t, "_")] = append(flatMap[strings.Join(t, "_")], "")
								l--
							}
							flatMap[strings.Join(t, "_")] = append(flatMap[strings.Join(t, "_")], f[k])
						} else {
							flatMap[strings.Join(t, "_")][position] = f[k]
						}

					}
					break
				}
			}
		} else {
			// has no array num in the key, should have unique key/value
			flatMap[k] = append(flatMap[k], f[k])
		}

		// trmK := re.ReplaceAllString(k, "$1")
		// csvKey := strings.TrimSuffix(strings.Replace(trmK, "__", "_", -1), "_")

		// flatMap[csvKey] = pos

		// // fmt.Printf("Pattern: %v \n", re.String()) //print pattern

		// repl
		// fmt.Println("trimded k is ", trmK)

	}

	// print CSV
	for i, k := range flatMap {
		fmt.Println("csv key : ", i)
		fmt.Println("csv value : ", k)
	}

	WriteCSV(flatMap)

}

func WriteCSV(csvMap map[string][]interface{}) {
	// sort the map
	keys := make([]string, 0, len(csvMap))

	for k := range csvMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	/*
		map[col] = [value1,value2...]

			col1      	col2
			value1		value1
			value2		value2
	*/
	// data := [][]string{}
	// 	0 {"id", "first_name", "last_name", "email"},
	// 	1 {"1", "Sau Sheong", "Chang", "mailto:someemail@random.com"},
	// 	2 {"2", "John", "Doe", "mailto:john@email.com"},
	// 	3 {"",  "Kent", "Doe", "mailto:john@email.com"},

	file, err := os.Create("new1.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// allocate csv row numbers
	x := 0
	for _, v := range csvMap {
		if x < len(v) {
			x = len(v)
		}
	}
	data := make([][]string, x+1) // allocate real len +1 for precaution
	// allocate csv col numbers
	for _, key := range keys {
		values := csvMap[key]
		for i := 0; i <= len(values); i++ {
			data[i] = make([]string, len(keys))
		}
	}

	pivotY := 0
	for _, key := range keys {
		values := csvMap[key]

		data[0][pivotY] = key
		for i := 1; i <= len(values); i++ {
			if values[i-1] == nil {
				data[i][pivotY] = ""
			} else {

				data[i][pivotY] = fmt.Sprintf("%v", values[i-1])
			}

			// if len(data[i]) == pivot {
			// 	if values[i-1] == nil {
			// 		data[i] = append(data[i], "")
			// 	} else {
			// 		data[i] = append(data[i], fmt.Sprintf("%v", values[i-1]))
			// 	}
			// } else {
			// 	// pivot always >= len(data[i])
			// 	l := pivot - len(data[i])
			// 	for l > 0 {
			// 		data[i] = append(data[i], "")
			// 		l--
			// 	}
			// 	if values[i-1] == nil {
			// 		data[i] = append(data[i], "")
			// 	} else {
			// 		data[i] = append(data[i], fmt.Sprintf("%v", values[i-1]))
			// 	}
			// }
		}
		pivotY++
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)

	if err != nil {
		fmt.Println("write error:  ", err)
	}
}

func recrusJson(jVector Json, cnt int, out KeyValue) (err error) {
	/*
		bool, for JSON booleans
		float64, for JSON numbers
		string, for JSON strings
		[]interface{}, for JSON arrays
		map[string]interface{}, for JSON objects
		nil for JSON null

	*/

	fmt.Println("counter", cnt)
	if cnt == 0 {
		jVector.Key = "JSON"
	}

	switch vv := jVector.Value.(type) {
	case string: // for JSON strings
		fmt.Println(jVector.Key, "  is string value:  ", vv)

		out[jVector.Key] = vv

	case float64: // for JSON numbers

		fmt.Println(jVector.Key, "  is float64 value: ", vv)
		out[jVector.Key] = vv
	case bool: //for JSON bool
		fmt.Println(jVector.Key, "  is bool value: ", vv)
		out[jVector.Key] = vv
	case []interface{}:

		fmt.Println(jVector.Key, "  is an array:  ", vv)
		for i, v := range vv {
			fmt.Println("i: ", i, " v: ", v)

			var jj Json
			jj.Value = v
			jj.Key = jVector.Key + "\t" + "'" + strconv.Itoa(i) + "'" //  array key use special sign

			recrusJson(jj, cnt+1, out)
		}
	case map[string]interface{}:

		fmt.Println(jVector.Key, "  is a map:  ", vv)
		for i, v := range vv {
			fmt.Println("i: ", i, " v: ", v)

			var jj Json
			jj.Value = v
			jj.Key = jVector.Key + "\t" + i

			recrusJson(jj, cnt+1, out)
		}
	case nil:
		return
	default:
		fmt.Println(vv, "is of a type I don't know how to handle")
	}

	return
}
