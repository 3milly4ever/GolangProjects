package main

import (
	"database/sql"
	"user-fitness/store"
)

func CreateAllTables(myStore *store.MySqlStore, db *sql.DB, sl *store.SqlLogger) error {
	// log := logger.NewLogger()

	// cachForMyStore := caching.NewRedisCache(redisClient)
	// myStore := store.NewMySqlStore(db)
	// sl := store.NewMySqlLogger(logger.NewLogger())

	tableData := []store.TableDefinition{
		{
			//dont expose internal id primary key to apis for security, others can get user info, expose fake id for api
			//generate uuids, replace auto increment with uuid which will be a string. a complex number that is impossible to guess.
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
		}
		{
			Name: "UserCredentials",
			Fields:`
			id UUID PRIMARY KEY,
			email VARCHAR(255),
			password VARCHAR(255)
			`,
		}
		},
	}
	err := myStore.CreateTables(tableData, sl)
	if err != nil {
		return err
	}
	return nil
}
