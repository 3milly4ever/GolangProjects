package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Connect() {
	//opening connection with my sql
	d, err := gorm.Open("mysql", "root:Em50goats@/fighter_management?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	db = d
}

// since access to the db variable is not available through other packages, we need a getter method
// to allow different packages to access the db variable which contains the gorm object to the database connection
func GetDB() *gorm.DB {
	return db
}
