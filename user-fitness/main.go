package main

import (
	"database/sql"
	"log"
	"net/http"
	"user-fitness/api"
	"user-fitness/logger"
	"user-fitness/store"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dataSourceName := "root:Em50goats@tcp(localhost:3306)/user_fitness"

	logger := logger.NewLogger()
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		logger.Error("Error connecting to the database: %v", err)
		return
	}
	defer db.Close()
	err = db.Ping()

	if err != nil {
		logger.Error("Error pinging the database: %v", err)
		return
	}
	logger.Info("Connection to the database successful")

	myStore := store.NewMySqlStore(logger, db)
	server := api.NewServer("localhost:9090", logger, myStore)
	http.HandleFunc("/users/", server.HandleUserRequests)

	err = CreateAllTables(db)
	if err != nil {
		logger.Error("Error creating tables", err)
		return
	}

	logger.Info("Server listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))

}
