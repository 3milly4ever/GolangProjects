package config

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func Connect(dataSourceName string) error {
	//dataSourceName := "root:Em50goats@tcp(localhost:3306)/user_info"

	//we open the connection and check for errors
	var err error
	db = GetDB()
	//open connection
	db, err = sql.Open("mysql", dataSourceName)
	//check for error connecting
	if err != nil {
		fmt.Println("Error connecting to the database", err)
		return err
	}
	//ping sends a small query or network packet, if the network responds the error will be nil
	err = db.Ping()
	if err != nil {
		panic("Error pinging the database:" + err.Error())
	}

	fmt.Println("Connection to the database successful")
	return nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing the database:", err)
		}
	}
}
