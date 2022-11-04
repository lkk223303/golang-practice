package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	xj "github.com/basgys/goxml2json"
)

// GOMAXPROCS 根據 CPU 有多少顆，做多少平行處理，但是可以透過 GOMAXPROCS 設定使用的 CPU 數量
// func init() {
// 	runtime.GOMAXPROCS(1)
// }

// Based on http://stackoverflow.com/questions/16466320/is-there-a-way-to-do-repetitive-tasks-at-intervals-in-golang
// Probably need to read more if some mutex if gotunite needs to access shared information.

// Based on http://stackoverflow.com/questions/16466320/is-there-a-way-to-do-repetitive-tasks-at-intervals-in-golang
// Probably need to read more if some mutex if gotunite needs to access shared information.

// const (
// 	duration      = 30 * time.Second
// 	interval_time = 5 * time.Second
// )

var (
	intflag    int
	boolflag   bool
	stringflag *string
)

type my struct {
	name string
	age  int
	man  bool `default:true`
}

func init() {
	flag.IntVar(&intflag, "intflag", 0, "int flag value")
	flag.BoolVar(&boolflag, "boolflag", false, "bool flag value")
	stringflag = flag.String("stringflag", "default", "string flag value")
}

func main() {
	parseJSONfile(getJSONbytes("users.json"))
}

func visit(friends []string, callback func(string)) {
	for _, n := range friends {
		callback(n)
	}
}

func writeCSVfile() {
	file, err := os.Create("new_users.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data := [][]string{
		{"id", "first_name", "last_name", "email"},
		{"1", "Sau Sheong", "Chang", "mailto:someemail@random.com"},
		{"2", "John", "Doe", "mailto:john@email.com"},
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		fmt.Println("write error:  ", err)
	}

	// one row at a time
	for _, row := range data {
		err = writer.Write(row)
		if err != nil {
			fmt.Println("write error:  ", err)
		}
	}
	writer.Flush()
}

type Recurluservers struct {
	XMLName     xml.Name `xml:"users"`
	Version     string   `xml:"version,attr"`
	Description string   `xml:",innerxml"`
	Users       []User   `xml:"user" json:"users"` //第一層nest要對到這
}
type userJson struct {
	User []User `json:"users"`
}
type User struct {
	XMLName  xml.Name `xml:"user"` //第一層nest要對到這
	UserName string   `xml:"name" json:"name"`
	Type     string   `xml:"type,attr" json:"type"`
	Age      int      `json:"age"`
	Social   Social   `xml:"social" json:"social"`
}
type Social struct {
	XMLName  xml.Name `xml:"social"`
	Facebook string   `xml:"facebook" json:"facebook"`
	Twitter  string   `xml:"twitter" json:"twitter"`
	Youtube  string   `xml:"youtube" json:"youtube"`
}
type Json struct {
	jKey string // json key

	jElemant interface{} //json element
}

func getJSONbytes(path string) []byte {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	return jsonFile
}

func parseJSONfile(jsonFile []byte) {
	// var users userJson

	var jVector Json
	json.Unmarshal(jsonFile, &jVector.jElemant)

	/*
		bool, for JSON booleans
		float64, for JSON numbers
		string, for JSON strings
		[]interface{}, for JSON arrays
		map[string]interface{}, for JSON objects
		nil for JSON null

	*/
	cnt := 0
	var recrusJson func(Json, int)
	recrusJson = func(jVector Json, cnt int) {
		fmt.Println("counter", cnt)
		switch vv := jVector.jElemant.(type) {
		case string: // for JSON strings
			fmt.Println(jVector.jKey, "  is string  ", vv)
		case float64: // for JSON numbers
			fmt.Println(jVector.jKey, "  is float64  ", vv)
		case []interface{}:
			fmt.Println(jVector.jKey, "  is an array:  ", vv)
			for i, v := range vv {
				fmt.Println("i: ", i, " v: ", v)
				var jj Json
				jj.jElemant = v
				jj.jKey = strconv.Itoa(i)
				recrusJson(jj, cnt+1)
			}
		case map[string]interface{}:
			fmt.Println(jVector.jKey, "  is a map:  ", vv)
			for i, v := range vv {
				fmt.Println("i: ", i, " v: ", v)
				var jj Json
				jj.jElemant = v
				jj.jKey = i
				recrusJson(jj, cnt+1)
			}
		case nil:
			return
		default:
			fmt.Println(vv, "is of a type I don't know how to handle")
		}
	}

	recrusJson(jVector, cnt)

	// for i, u := range jVector.(map[string]interface{}) {
	// 	fmt.Println("i: ", i)

	// 	// assert type
	// 	switch vv := u.(type) {
	// 	case string:
	// 		fmt.Println(u, "is string", vv)
	// 	case int:
	// 		fmt.Println(u, "is int", vv)
	// 	case float64:
	// 		fmt.Println(u, "is float64", vv)
	// 	case []interface{}:
	// 		fmt.Println(u, "is an array:")
	// 		for i, v := range vv {
	// 			fmt.Println(i, v)
	// 		}
	// 	case map[string]interface{}:

	// 		for i2, v := range vv {
	// 			fmt.Println("i2: ", i)
	// 			fmt.Println("v2: ", v)

	// 			// assert type
	// 			switch vv := v.(type) {
	// 			case string:
	// 				fmt.Println(i2, "is string", vv)
	// 			case int:
	// 				fmt.Println(i2, "is int", vv)
	// 			case float64:
	// 				fmt.Println(i2, "is float64", vv)
	// 			case []interface{}:
	// 				fmt.Println(i2, "is an array:")
	// 				for i, u := range vv {
	// 					fmt.Println(i, u)
	// 				}
	// 			case map[string]interface{}:

	// 				fmt.Println(v, "is an map:")
	// 				for i3, u := range vv {
	// 					fmt.Println(i3, u)

	// 					switch vv := v.(type) {
	// 					case string:
	// 						fmt.Println(i3, "is string", vv)
	// 					case int:
	// 						fmt.Println(i3, "is int", vv)
	// 					case float64:
	// 						fmt.Println(i3, "is float64", vv)
	// 					case []interface{}:
	// 						fmt.Println(i3, "is an array:")
	// 						for i, u := range vv {
	// 							fmt.Println(i, u)
	// 						}
	// 					case map[string]interface{}:

	// 						fmt.Println(i3, "is an map:")
	// 						for i4, u := range vv {
	// 							fmt.Println(i4, u)

	// 							switch vv := v.(type) {
	// 							case string:
	// 								fmt.Println(i4, "is string", vv)
	// 							case int:
	// 								fmt.Println(i4, "is int", vv)
	// 							case float64:
	// 								fmt.Println(i4, "is float64", vv)
	// 							case []interface{}:
	// 								fmt.Println(i4, "is an array:")
	// 								for i, u := range vv {
	// 									fmt.Println(i, u)
	// 								}
	// 							case map[string]interface{}:

	// 								fmt.Println(i4, "is an map:")
	// 								for i5, u := range vv {
	// 									fmt.Println(i, u)

	// 									switch vv := v.(type) {
	// 									case string:
	// 										fmt.Println(i5, "is string", vv)
	// 									case int:
	// 										fmt.Println(i5, "is int", vv)
	// 									case float64:
	// 										fmt.Println(i5, "is float64", vv)
	// 									case []interface{}:
	// 										fmt.Println(i5, "is an array:")
	// 										for i, u := range vv {
	// 											fmt.Println(i, u)
	// 										}
	// 									case map[string]interface{}:

	// 										fmt.Println(i5, "is an map:")
	// 										for i, u := range vv {
	// 											fmt.Println(i, u)
	// 										}
	// 									default:
	// 										fmt.Println(i5, "is of a type I don't know how to handle")

	// 									}
	// 								}
	// 							default:
	// 								fmt.Println(i4, "is of a type I don't know how to handle")

	// 							}
	// 						}
	// 					default:
	// 						fmt.Println(i3, "is of a type I don't know how to handle")
	// 					}
	// 				}
	// 			default:
	// 				fmt.Println(i2, "is of a type I don't know how to handle")
	// 			}
	// 		}
	// 	default:
	// 		fmt.Println(i, "is of a type I don't know how to handle")

	// 	}
	// }

	// csv.NewWriter(os.Stdout)
}

func transXML() []byte {
	// upload xml
	xml, err := os.Open("users.xml")
	if err != nil {
		log.Println(err)
	}
	defer xml.Close()
	json, err := xj.Convert(xml)
	if err != nil {
		panic("convert problem")
	}

	return json.Bytes()
}

func getXMLfile() {

	data, err := os.ReadFile("users.xml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// v := Recurluservers{}
	xmap := make(map[string]interface{})
	err = xml.Unmarshal(data, &xmap)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, v := range xmap {
		fmt.Println("i: ", i)
		fmt.Println("v: ", v)
	}

	// fmt.Println("v des: ", v.Description)
	// fmt.Println("v user: ", v.Users)
	// fmt.Println("v ver: ", v.Version)
	// fmt.Println("v xml: ", v.XMLName)
	// fmt.Println("v xmlName:  ", v.XMLName)
	// for _, i := range v.Users {
	// 	fmt.Println(i.UserName)
	// 	fmt.Println(i.Type)
	// 	fmt.Println(i.Social)
	// }

}

func goFuncExample() {
	done := make(chan bool)

	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func(i string) {
			fmt.Println(i)
			done <- true
		}(v)
	}

	// wait for all goroutines to complete before exiting
	for _ = range values {

		fmt.Println(<-done)
	}
}

func flagPrac() {
	flag.Parse()

	for i := range os.Args {
		fmt.Printf("Args %d: %v\n", i, os.Args[i])
	}
	fmt.Println("int flag : ", intflag)
	fmt.Println("bool flag : ", boolflag)
	fmt.Println("string flag : ", *stringflag)
}

func sliceSetBymap() {
	s := "this_is_a_test_slice_a_test_slice_is_kent_slice_handsome"
	sli := strings.Split(s, "_")
	m := make(map[string]string)
	for _, k := range sli {
		if _, ok := m[k]; !ok {
			log.Println("value: ", k)
			m[k] = k
		}
	}
}

func channelPractice() {

	outchan := make(chan int)
	errChan := make(chan error)
	finishChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(outchan chan<- int, errChan chan<- error, val int, wg *sync.WaitGroup) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			fmt.Println("finished job id : ", val)
			outchan <- val
			if val == 60 {
				errChan <- errors.New("error in 60")
			}

		}(outchan, errChan, i, &wg)
	}

	go func() {
		wg.Wait()
		fmt.Println("finish all job")
		close(finishChan)
	}()

Loop:
	for {
		select {
		case val := <-outchan:
			fmt.Println("finished: ", val)
		case err := <-errChan:
			fmt.Println("error: ", err)
			break Loop
		case <-finishChan:
			break Loop
			// fmt.Println(finishChan)
		case <-time.After(10000 * time.Millisecond):
			break Loop
		}
	}
}

func getPipeState() {
	for {
		log.Println("go")
		time.Sleep(2 * time.Second)

	}
}

// func main() {
// 	// name := "bigobject_license_ae751d28-8f8b-4353-80b2-55c751d56dc8_m67"
// 	// name_slice := strings.Split(name, "_")
// 	// last_file_name_token := name_slice[len(name_slice)-1]

// 	// var s []string
// 	// str := "1"
// 	// fmt.Println(len(str))

// 	// // s = append(s, "")
// 	// fmt.Println(s == nil)

// 	// // name_slice2 := append(name_slice[:1], name_slice[2:]...)
// 	// fmt.Println("before append: ", name_slice)
// 	// name_slice = append(name_slice[:1], name_slice[2:]...)
// 	// fmt.Println("after append: ", name_slice)
// 	// fmt.Println(last_file_name_token)
// 	// show_name := strings.Join(name_slice, "_")
// 	// // fmt.Println(name_slice2)
// 	// fmt.Println(show_name)

// 	// fmt.Println(interfaceFunc())

// 	for i := 10; i > 0; i-- {
// 		foo := addByShareMemory(10)
// 		fmt.Println(len(foo))
// 		fmt.Println(foo)
// 	}
// }

func addByShareMemory(n int) []int {
	var ints []int
	// var wg sync.WaitGroup
	// var mux sync.Mutex
	channel := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(channel chan<- int, order int) {
			channel <- order
		}(channel, i)
	}

	// wg.Add(n)
	for i := range channel {
		ints = append(ints, i)
		if len(ints) == n {
			break
		}
	}
	close(channel)

	// wg.Wait()
	return ints
}

func interfaceFunc() (map[string]interface{}, error) {

	content := map[string]interface{}{
		"red":   map[string]interface{}{"red1": "1", "red2": "2"},
		"green": map[string]interface{}{"green1": "1", "green2": "2"},
	}

	var v map[string]interface{}
	for _, val := range content {

		v, _ = val.(map[string]interface{})
		// if !ok {
		// 	return nil, fmt.Errorf("unable get content")
		// }
		return v, nil
	}

	return nil, fmt.Errorf("unable get content")
}

type worker interface {
	work()
}

type person struct {
	name string
	worker
}

func (p person) work() {
	fmt.Println("name: ", p.name)
}
func workerInstant() {
	var w worker = person{name: "Kent"}
	w.work()
	fmt.Println("w: ", w)
}
