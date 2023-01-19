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
* toml.go
*
* @author:   Grace Chen, Kent Huang
* ----------------------------------------------------------------------------
*/
	
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Env struct {
	// SMTP
	Email    string `toml:"SENDER_EMAIL"`
	Name     string `toml:"SENDER_NAME"`
	Host     string `toml:"SMTP_HOST"`
	Port     string `toml:"SMTP_PORT"`
	User     string `toml:"SMTP_USER"`
	Password string `toml:"SMTP_PSWD"`
	// SLACK
	SlackUrl string `toml:"SLACK_URL"`
}

func setEnvToFile(env Env) error {
	file, err := os.OpenFile("EnvPath", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	bts, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	defer file.Close()
	ee := make(map[string]interface{})
	if _, err := toml.Decode(string(bts), &ee); err != nil {
		return err
	}

	ee["SENDER_EMAIL"] = env.Email
	ee["SENDER_NAME"] = env.Name
	ee["SMTP_HOST"] = env.Host
	ee["SMTP_PORT"] = env.Port
	ee["SMTP_USER"] = env.User
	ee["SMTP_PSWD"] = env.Password
	ee["SLACK_URL"] = env.SlackUrl

	buf := new(bytes.Buffer)
	tomlEncoder := toml.NewEncoder(buf)
	tomlEncoder.Indent = ""
	if err := tomlEncoder.Encode(ee); err != nil {
		return err
	}

	// 這裡是一個大坑，用osfile writeString 會寫出錯誤文件
	// if _, err := file.WriteString(buf.String()); err != nil {
	// 	log.Println("write err ", err)
	// }

	err = ioutil.WriteFile(file.Name(), buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
func fileEncode() error {

	file, err := os.OpenFile("EnvPath", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	s := `SENDER_EMAIL = "bigobject.iae2@gmail.com"
SENDER_NAME = "bigobject"
SMTP_HOST = "smtp.gmailx.com"
SMTP_PORT = "587"
SMTP_USER = "bigobject.iae2@gmail.com"
SMTP_PSWD = "123"
SLACK_URL = "iaesla"
FERNET_KEY="Eq7zT_fBIxwqlWMffGVEnj64GYv8UJhusYraFbm6E9Q="
SENDER_EMAIL = "gracechen@bigobject.io"
SENDER_NAME = "IAE_52"
SMTP_HOST = "msa.hinet.net"
SMTP_PORT = "25"
SMTP_USER = ""
SMTP_PSWD = ""
SLACK_URL = "iaesla"aaa
`
	// buf := new(bytes.Buffer)
	// tomlEncoder := toml.NewEncoder(buf)
	// tomlEncoder.Indent = ""
	// if err := tomlEncoder.Encode(env); err != nil {
	// 	return fmt.Errorf("env encode: %s", err)
	// }

	x, err := file.WriteString(s)
	if err != nil {
		log.Println("err", err)
	}
	fmt.Println(x)
	return nil

}
