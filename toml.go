package main

import (
	"bytes"
	"io/ioutil"
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
