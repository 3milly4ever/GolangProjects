//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"time"
	"user-fitness/logger"
	"user-fitness/store"
	"user-fitness/types"
)

func main() {
	sl := store.NewMySqlLogger(logger.NewLogger())

	// Create a new store instance
	dataSourceName := "root:Em50goats@tcp(localhost:3306)/user_fitness"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		sl.Logger.Error("Error connecting to the database: %v", err)
		return
	}
	defer db.Close()

	myStore := store.NewMySqlStore(db)

	for i := 0; i < 100; i++ {
		// Customize the user data as needed
		user := types.User{
			Name:       fmt.Sprintf("User %d", i+1),
			Email:      fmt.Sprintf("user%d@example.com", i+1),
			Weight:     150 + i,
			Goal:       "Maintain",
			Regimen:    "Regular workout",
			DateJoined: time.Now().Format("2006-01-02"),
		}

		// Insert the user into the database
		_, err := myStore.InsertUser(&user, sl)
		if err != nil {
			sl.Logger.Error("Error inserting user: %v", err)
		}
	}
}
