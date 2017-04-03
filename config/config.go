package config

import (
	"os"
	"encoding/json"
	"fmt"
)

type Configuration struct {
	DB DBParams `json:"db`
}

type DBParams struct {

	Type string `json:"type"`
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Name string `json:"name"`
	Pass string `json:"pass"`

}

var Config Configuration


func LoadConfigFile(filePath string) (err error) {
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		return err
	}
	if err = json.NewDecoder(file).Decode(&Config); err != nil {
		fmt.Println(Config.DB.Type)
		return err
	}
	return nil

}

func (db *DBParams) ConnectionString() (string, error) {
	if db.Type == "postgres" {
		return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable", db.Name, db.User, db.Pass, db.Host, db.Port), nil
	} else if db.Type == "mysql" {
		return fmt.Sprintf("%s:%s@/%s", db.User, db.Pass, db.Name), nil
	} else {
		return "", fmt.Errorf("Unknown database type")
	}

}
