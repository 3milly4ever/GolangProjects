// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"user-fitness/logger"
// 	"user-fitness/store"

// 	_ "github.com/go-sql-driver/mysql"
// )

package main

import (
	"database/sql"
	"user-fitness/logger"
	"user-fitness/store"
)

func CreateAllTables(db *sql.DB) error {
	log := logger.NewLogger()
	myStore := store.NewMySqlStore(log, db)

	tableData := []store.TableDefinition{
		{
			Name: "Users",
			Fields: `
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255),
            email VARCHAR(255),
            weight INT,
            goal VARCHAR(255),
            regimen VARCHAR(255),
            date_joined DATE
            `,
		},
	}
	err := myStore.CreateTables(tableData)
	if err != nil {
		return err
	}
	return nil
}
