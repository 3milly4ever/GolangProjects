package main

import (
	"database/sql"
	"log"
	"net/http"
	"user-fitness/api"
	"user-fitness/caching"
	"user-fitness/logger"
	"user-fitness/store"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dataSourceName := "root:Em50goats@tcp(localhost:3306)/user_fitness"
	sl := store.NewMySqlLogger(logger.NewLogger())

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		sl.Logger.Error("Error connecting to the database: %v", err)
		return
	}
	defer db.Close()
	err = db.Ping()

	if err != nil {
		sl.Logger.Error("Error pinging the database: %v", err)
		return
	}
	sl.Logger.Info("Connection to the database successful")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cachForMyStore := caching.NewRedisCache(redisClient)
	myStore := store.NewMySqlStore(db, cachForMyStore)
	server := api.NewServer("localhost:9090", db, redisClient)
	http.HandleFunc("/users/", server.HandleUserRequests)

	err = CreateAllTables(myStore, db, sl)
	if err != nil {
		sl.Logger.Error("Error creating tables", err)
		return
	}

	sl.Logger.Info("Server listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))

}
