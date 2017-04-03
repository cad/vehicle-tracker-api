package repository

import (
//	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB



func ConnectDB(dbType string, dbURL string) {
	var err error
	db, err = gorm.Open(dbType, dbURL)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Vehicle{})
}

func CloseDB() {
	db.Close()
}
