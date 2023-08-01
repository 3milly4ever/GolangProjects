package models

import (
	"database/sql"
	"fighter-database/pkg/config"
	f "fmt"
)

type UserRecord struct {
	ID         int    `json:"id"`
	Weight     int    `json:"weight"`
	Goal       string `json:"goal"`
	Regimen    string `json:"regimen"`
	DateJoined string `json:"date_joined"`
}

func CreateTable(tableName string) error {
	createTableQuery := f.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
		id INT AUTO_INCREMENT PRIMARY KEY,
		weight INT,
		goal VARCHAR(255),
		regimen VARCHAR(255),
		date_joined DATE
	)`, tableName)
	//if no error db.Exec executes otherwise an error is returned which is not nil or nil if no error.
	_, err := config.GetDB().Exec(createTableQuery)
	return err
}

func NewUserRecord(id, weight int, goal string, regimen string, dateJoined string) UserRecord {
	return UserRecord{
		ID:         id,
		Weight:     weight,
		Goal:       goal,
		Regimen:    regimen,
		DateJoined: dateJoined,
	}
}

func InsertUserData(user UserRecord) error {
	insertQuery := `
		INSERT INTO user (weight, goal, regimen, date_joined)
		VALUES (?, ?, ?, ?)
	`
	_, err := config.GetDB().Exec(insertQuery, user.Weight, user.Goal, user.Regimen, user.DateJoined)
	return err
}

func DeleteUserData(userID int) error {
	deleteQuery := `
	DELETE FROM user
	WHERE id = ?
	`
	_, err := config.GetDB().Exec(deleteQuery, userID)
	return err
}

func UpdateUserData(userID int, user UserRecord) error {
	updateQuery := `
	UPDATE user
	SET weight = ?, goal = ?, regimen = ?, date_joined = ?
	WHERE id = ?
	`

	_, err := config.GetDB().Exec(updateQuery, user.Weight, user.Goal, user.Regimen, user.DateJoined, userID)
	return err

}

func GetAllUserData() ([]UserRecord, error) {
	getAllQuery := `
	SELECT * FROM user
	`
	rows, err := config.GetDB().Query(getAllQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRecords []UserRecord

	for rows.Next() {
		var user UserRecord

		err := rows.Scan(&user.ID, &user.Weight, &user.Goal, &user.Regimen, &user.DateJoined)
		if err != nil {
			return nil, err
		}

		userRecords = append(userRecords, user)
	}
	//check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userRecords, nil

}

// get user by id
func GetUserById(userID int) (UserRecord, error) {
	getUserByIdQuery := `
	SELECT * FROM user
	WHERE id = ?
	`
	row := config.GetDB().QueryRow(getUserByIdQuery, userID)

	var user UserRecord

	err := row.Scan(&user.ID, &user.Weight, &user.Goal, &user.Regimen, &user.DateJoined)
	//check for errors with scan
	if err != nil {
		if err == sql.ErrNoRows {
			return UserRecord{}, f.Errorf("user with ID %d not found", userID)
		}
		//for other errors, return the error as is
		return UserRecord{}, err
	}
	//if no errors
	return user, nil
}
