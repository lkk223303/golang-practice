/*
* ----------------------------------------------------------------------------
* Copyright (c) 2022-present BigObject Inc.
* All Rights Reserved.
*
* Use of, copying, modifications to, and distribution of this software
* and its documentation without BigObject's written permission can
* result in the violation of U.S., Taiwan and China Copyright and Patent laws.
* Violators will be prosecuted to the highest extent of the applicable laws.
*
* BIGOBJECT MAKES NO REPRESENTATIONS OR WARRANTIES ABOUT THE SUITABILITY OF
* THE SOFTWARE, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
* TO THE IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
* PARTICULAR PURPOSE, OR NON-INFRINGEMENT.
*
*
* loadCSV.go
*
* @author:   Grace Chen, Kent Huang
* ----------------------------------------------------------------------------
*/
	
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var Pwd string
var FilePath string
var FileName = "sales.csv"

func Load() {
	var wrCsv [][]string
	// get pwd
	Pwd, _ = os.Getwd()

	FilePath = filepath.Join(Pwd, FileName)
	// fmt.Println(FilePath)
	file, err := os.OpenFile(FilePath, os.O_RDONLY, 0777) // os.O_RDONLY 表示只讀、0777 表示(owner/group/other)權限
	if err != nil {
		log.Fatalln("找不到CSV檔案路徑:", FilePath, err)
	}

	// read
	r := csv.NewReader(file)
	// r.Comma = ',' // 以何種字元作分隔，預設為`,`。所以這裡可拿掉這行
	// count := 1
	for {
		// read() 一次讀一行 readAll()一次讀全部
		wrCsv, err = r.ReadAll()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		for k, v := range wrCsv {
			if k == 100 {
				break
			}
			fmt.Println(k, " :  ", v)
		}
		break
		// count++
	}

	wFile, _ := os.Create("csvRewrite.csv")
	w := csv.NewWriter(wFile)
	if err = w.WriteAll(wrCsv); err != nil {
		log.Println(err)
	}
}
