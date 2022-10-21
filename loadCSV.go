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
		record, err := r.ReadAll()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		for k, v := range record {
			if k == 100 {
				break
			}
			fmt.Println(k, " :  ", v)
		}
		break
		// count++
	}
}
