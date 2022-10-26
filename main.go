package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
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

func init() {
	flag.IntVar(&intflag, "intflag", 0, "int flag value")
	flag.BoolVar(&boolflag, "boolflag", false, "bool flag value")
	stringflag = flag.String("stringflag", "default", "string flag value")
}

func main() {

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

// func GetResponseFromStreamerXPipe(pipeID string) (oapi.ModelStreamerX, error) {
// 	/*
// 		type AddStreamerXJSONBody struct {
// 			BoID            ColumnUUID                       `json:"boID"`
// 			BoWorkspaceName string                           `json:"boWorkspaceName"`
// 			Freq            string                           `json:"freq"`
// 			Meta            ModelStreamerXMetaInfo           `json:"meta"`
// 			Path            string                           `json:"path"`
// 			Regex           string                           `json:"regex"`
// 			TableDef        []ModelStreamerXTableCreationDef `json:"tableDef"`
// 		}
// 	*/
// 	var reserv oapi.ModelStreamerX

// 	flat_cfg_str, err := GetPipeCfg(pipeID)
// 	if err != nil {
// 		return reserv, errors.New("CANNOT get streamerX configuration")
// 	}

// 	pipe_status, err := GetPipeStatus(pipeID)
// 	if err != nil {
// 		return reserv, errors.New("CANNOT get streamerX pipe status")
// 	}

// 	var obj_config = map[string]interface{}{}
// 	json.Unmarshal([]byte(flat_cfg_str), &obj_config)

// 	// globalLogger.Debugf("[getResponseFromStreamerXPipe]obj_config:%v", obj_config)

// 	//var main_config_asserted map[string]interface{}
// 	if main_config, exist := obj_config["main"]; !exist {
// 		return reserv, errors.New("CANNOT find streamerX MAIN configuration")
// 	} else {
// 		// globalLogger.Debugf("[getResponseFromStreamerXPipe]main_config:%v", main_config)
// 		jsonBytes_config, err := json.Marshal(main_config)
// 		if err != nil {
// 			return reserv, err
// 		}

// 		json.Unmarshal(jsonBytes_config, &reserv)
// 	}

// 	/*
// 		type ModelStreamerX struct {
// 			BoID ColumnUUID `json:"boID"`
// 			BoWorkspaceName ColumnName             `json:"boWorkspaceName"`
// 			PipeID ColumnUUID `json:"pipeID"`
// 			Freq            string                 `json:"freq"`
// 			Meta            ModelStreamerXMetaInfo `json:"meta"`
// 			Path            string                 `json:"path"`
// 			Regex  string     `json:"regex"`
// 		}
// 	*/
// 	reserv.PipeID = oapi.ColumnUUID(pipeID)
// 	reserv.Status = pipe_status

// 	// globalLogger.Debugf("[getResponseFromStreamerXPipe]ModelStreamerX:%v", reserv)

// 	return reserv, nil
// }

// func getBigobjectStatus(p pipe.Pipe, contents map[string]interface{}) (string, error) {
// 	// var err error
// 	var stateList []string
// 	loggerPrefix := "[PipeServBigobjectStatus]"

// 	globalLogger.Debug("PIPE TYPE: ", p.Type)
// 	if p.Type == pipe.STREAMER_X_TYPE {
// 		strXModel, err := GetResponseFromStreamerXPipe(p.ID)

// 		bigobject, err := boLite.GetBigobject(string(strXModel.BoID))

// 		if err == sql.ErrNoRows {
// 			// were there are no bigobject in strX
// 			return "missing", nil
// 		}
// 		if err != nil {
// 			return "", err
// 		}

// 		// if content has value from nagios service use it directly
// 		if contents != nil {
// 			globalLogger.Debug(loggerPrefix, "Get Content Directly from cache")
// 			// get selected bo content from contents
// 			boHostName := bigobject.MachineID + "_bigobject_" + bigobject.Name
// 			selectedContent, err := selectNagiosContent(boHostName, contents)
// 			if err != nil {
// 				return "", err
// 			}
// 			globalLogger.Debug(loggerPrefix, " streamerX BO: ", boHostName)

// 			/*
// 				get services from selected BO
// 				ex:
// 					ae751d28-8f8b-4353-80b2-55c751d56dc8_bigobject_BO0906:
// 						services:
// 			*/
// 			services, ok := selectedContent["services"].(map[string]interface{})
// 			if !ok {
// 				return "", fmt.Errorf(loggerPrefix, "unable get nagios Bo response services")
// 			}

// 			for name, service := range services {
// 				// get all bigobject service state
// 				if strings.Contains(name, "bigobject_") {
// 					globalLogger.Debug(loggerPrefix, "service name: ", name)
// 					currentState, ok := service.(map[string]interface{})["current_state"].(string)
// 					if !ok {
// 						return "", fmt.Errorf(loggerPrefix, "unable get nagios bigobject service current_state")
// 					}
// 					globalLogger.Debug(loggerPrefix, "service currentState : ", currentState)
// 					_iae_state := bowatch.GetIaeStateByServiceState(currentState)
// 					globalLogger.Debug(loggerPrefix, "service _iae_state : ", _iae_state)
// 					stateList = append(stateList, _iae_state)
// 				}
// 			}
// 		} else {
// 			// if content no value, get machine status from bowatch (slow approach)
// 			bigobjectStatus, err := bowatch.GetBigobjectStatus(bigobject.MachineID, bigobject.Name)
// 			if err != nil {
// 				return "", err
// 			}
// 			globalLogger.Debug("bigobject status: ", bigobjectStatus, " from pipe ", p.Name)
// 			stateList = append(stateList, bigobjectStatus)
// 		}
// 		// return bowatch.GetWorstState(stateList)
// 	} else {
// 		// get bigobject by streamers
// 		streamers, err := strmSrvs.GetStreamerByPipe(p.ID)
// 		if err != nil {
// 			return "", err
// 		}
// 		if streamers != nil {
// 			// get BO status
// 			for _, s := range streamers {
// 				// get boID from streamer
// 				bigobject, err := boLite.GetBigobject(s.BOID)
// 				if err != nil {
// 					return "", err
// 				}

// 				// if content has value from nagios service use it directly
// 				if contents != nil {
// 					globalLogger.Debug(loggerPrefix, "Get Content Directly from cache")
// 					// get selected bo content from contents
// 					boHostName := bigobject.MachineID + "_bigobject_" + bigobject.Name
// 					selectedContent, err := selectNagiosContent(boHostName, contents)
// 					if err != nil {
// 						return "", err
// 					}

// 					/*
// 						get services from selected BO
// 						ex:
// 							ae751d28-8f8b-4353-80b2-55c751d56dc8_bigobject_BO0906:
// 								services:
// 					*/
// 					services, ok := selectedContent["services"].(map[string]interface{})
// 					if !ok {
// 						return "", fmt.Errorf(loggerPrefix, "unable get nagios Bo response services")
// 					}

// 					for name, service := range services {
// 						// get all bigobject service state
// 						if strings.Contains(name, "bigobject_") {
// 							globalLogger.Debug(loggerPrefix, "service name: ", name)
// 							currentState, ok := service.(map[string]interface{})["current_state"].(string)
// 							if !ok {
// 								return "", fmt.Errorf(loggerPrefix, "unable get nagios bigobject service current_state")
// 							}
// 							globalLogger.Debug(loggerPrefix, "service currentState : ", currentState)
// 							_iae_state := bowatch.GetIaeStateByServiceState(currentState)
// 							globalLogger.Debug(loggerPrefix, "service _iae_state : ", _iae_state)
// 							stateList = append(stateList, _iae_state)
// 						}
// 					}
// 				} else {
// 					// if content no value, get machine status from bowatch (slow approach)
// 					bigobjectStatus, err := bowatch.GetBigobjectStatus(bigobject.MachineID, bigobject.Name)
// 					if err != nil {
// 						return "", err
// 					}
// 					globalLogger.Debug("bigobject status: ", bigobjectStatus, " from pipe ", p.Name)
// 					stateList = append(stateList, bigobjectStatus)
// 				}
// 			}
// 		} else {
// 			return GetPipeStatus(p.ID)
// 		}
// 	}

// 	return bowatch.GetWorstState(stateList)
// }
