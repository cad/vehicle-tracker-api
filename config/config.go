package config

import (
	"os"
	"encoding/json"
	"fmt"
)

const VERSION = "0.1.1"

type Configuration struct {
	DB DBParams `json:"db`
	Server ServerParams `json:"server"`
}

type DBParams struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ServerParams struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

var C Configuration


func LoadConfigFile(filePath string) (err error) {
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		return err
	}
	if err = json.NewDecoder(file).Decode(&C); err != nil {
		fmt.Println(C.DB.Type)
		return err
	}
	return nil

}
