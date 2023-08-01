package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// helps other files interact with the data base
var (
	db *gorm.DB
)

func Connect() {

	//open connection with the mySQL database using the GORM library
	d, err := gorm.Open("mysql", "root:Em50goats@/book_management?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	db = d //the db from above var
}

func GetDB() *gorm.DB {
	return db
}
