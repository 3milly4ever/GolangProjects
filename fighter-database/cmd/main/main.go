package main

import (
	"fighter-database/pkg/config"
	"fighter-database/pkg/models"
	"fighter-database/pkg/routes"
	f "fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//var db *sql.DB

func main() {
	dataSourceName := "root:Em50goats@tcp(localhost:3306)/user_info"
	//we create the table and check for error
	var err error
	err = config.Connect(dataSourceName)
	if err != nil {
		f.Println("Error connecting to the database", err)
		return
	}

	defer config.CloseDB()

	tableName := "user"
	err = models.CreateTable(tableName)
	if err != nil {
		f.Println("Error creating table:", err)
		return
	}
	f.Println("Table creation successful!")

	routes.RegisterUserRoutes()

	f.Println("Server listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
