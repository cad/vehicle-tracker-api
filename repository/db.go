package repository

import (
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
	db.AutoMigrate(
		&Vehicle{},
		&Agent{},
		&Group{},
	)
}

func CloseDB() {
	db.Close()
}
