package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const VERSION = "1.1.10"

type Configuration struct {
	DB     DBParams     `json:"db`
	Server ServerParams `json:"server"`
}

type DBParams struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ServerParams struct {
	Port string `json:"port"`
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
