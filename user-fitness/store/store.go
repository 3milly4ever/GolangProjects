package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-fitness/logger"
	"user-fitness/types"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlStore struct {
	db     *sql.DB
	logger logger.Logger
}

// var (
// 	// ErrTableNotFound is a custom error indicating that a table does not exist
// 	ErrTableNotFound = errors.New("table not found")
// )

// var db *sql.DB

type Store interface {
	HandleInsertUser(w http.ResponseWriter, r *http.Request)
	HandleDeleteUser(w http.ResponseWriter, r *http.Request)
	HandleUpdateUser(w http.ResponseWriter, r *http.Request)
	HandleGetAllUsers(w http.ResponseWriter, r *http.Request)
	HandleGetUserById(w http.ResponseWriter, r *http.Request)
	CreateTableWithFields(tableName string, fields string) error
	CreateTables(tables []TableDefinition) error
}

// this makes sense if we have more fields in the MySqlStore struct.
func NewMySqlStore(logger logger.Logger, db *sql.DB) *MySqlStore {
	return &MySqlStore{
		db:     db,
		logger: logger,
	}
}

type TableDefinition struct {
	Name   string
	Fields string
}

// var Logger = logger.NewLogger()
// var store = NewMySqlStore(ogger, db)

// sets content-type header, writes the http status code, and writes the JSON response
func WriteJSONResponse(w http.ResponseWriter, statusCode int, jsonResponse []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func (s *MySqlStore) CreateTables(tables []TableDefinition) error {
	for _, table := range tables {
		err := s.CreateTableWithFields(table.Name, table.Fields)
		if err != nil {
			s.logger.Error("Error creating table %s: %v", table.Name, err)
			return err
		} else {
			s.logger.Info("Table %s created successfully", table.Name)
		}
	}
	return nil
}

// we create a table in our database
func (s *MySqlStore) CreateTableWithFields(tableName string, fields string) error {
	//	fieldStr := strings.Join(fields, ", ")

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
		%s
	)`, tableName, fields)

	//either one of these will happen if error or not error
	_, err := s.db.Exec(createTableQuery)
	return err
}

// will get a user based on ID
func (s *MySqlStore) GetUserById(id int) (types.User, error) {
	getUserByIdQuery := `
	SELECT * FROM user
	WHERE id = ?
	`
	//create new instance of a type that satisfies the Logger interface.
	row := s.db.QueryRow(getUserByIdQuery, id)

	var user types.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Weight, &user.Goal, &user.Regimen, &user.DateJoined)

	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("user with id %d not found", id)
			return types.User{}, err //returns empty user
		}
		s.logger.Error("database error: %v", err)
		return types.User{}, err //returns empty user
	}
	return user, nil
}

// handles the GET http request to get user by id
func (s *MySqlStore) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.HttpError(w, http.StatusMethodNotAllowed, "Invalid method request, should be GET")
		return
	}

	userIDStr := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		s.logger.HttpError(w, http.StatusBadRequest, "Invalid user id")
		return //return so the function ends after an error
	}

	user, err := s.GetUserById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.HttpError(w, http.StatusNotFound, err.Error())
			return
		}
		s.logger.HttpError(w, http.StatusNotFound, "Error retrieving user by id")
		return
	}

	//we parse(convert) the dateJoined string into a time value
	dateJoined, err := time.Parse("2006-01-02", user.DateJoined)
	if err != nil {
		s.logger.HttpError(w, http.StatusInternalServerError, "Error parsing into time object")
		return
	}

	user.DateJoined = dateJoined.Format("2006-01-02") //the format has a time receiver and a return of type string
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		s.logger.HttpError(w, http.StatusInternalServerError, "Error creating JSON response for user by id")
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(jsonResponse)
	WriteJSONResponse(w, http.StatusOK, jsonResponse)
}

// will create/insert a new user into the database
func (s *MySqlStore) InsertUser(user *types.User) error {
	insertQuery := `
	INSERT INTO user (name, email, weight, goal, regimen, date_joined)
	VALUES(?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(insertQuery, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined)
	//error check and provide user name
	if err != nil {
		s.logger.Error("Error inserting user %s: %v", user.Name, err)
		return err
	}

	return nil
}

// will handle the http request to create/insert a new user
func (s *MySqlStore) HandleInsertUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logger.HttpError(w, http.StatusMethodNotAllowed, "Invalid request method, should be post")
		return
	}
	//
	var user types.User                          //we will decode the json data into this
	err := json.NewDecoder(r.Body).Decode(&user) //we parse the data the client sent in json to create new user
	if err != nil {
		s.logger.HttpError(w, http.StatusBadRequest, "Invalid JSON data during post")
		return
	}
	//now that we decoded the body into user we can create a new user with the credentials
	newUser := types.NewUser(user.ID, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined)
	s.InsertUser(newUser)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Data posted successfuly!")
}

// delete a user from the database
func (s *MySqlStore) DeleteUser(id int) error {
	deleteQuery := `
	DELETE * FROM user
	WHERE id = ?
	`

	_, err := s.db.Exec(deleteQuery, id)
	if err != nil {
		s.logger.Error("Error deleting user with id: %d", id)
		return err
	}
	return nil
}

// handles http request for deleting user from database
func (s *MySqlStore) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		s.logger.HttpError(w, http.StatusBadRequest, "Error, wrong request should be DELETE")
		return
	}

	userIDStr := strings.TrimPrefix(r.URL.Path, "users/delete/") //get id from url path
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		s.logger.HttpError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = s.DeleteUser(userID)
	if err != nil {
		s.logger.HttpError(w, http.StatusInternalServerError, "Error deleting user data")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User data deleted successfuly!")
}

// updates user in our database
func (s *MySqlStore) UpdateUser(id int, user *types.User) error {
	updateQuery := `
	UPDATE user
	SET name = ?, email = ?, weight = ?, goal = ?, regimen = ?, date_joined = ?
	WHERE id = ?
	`

	_, err := s.db.Exec(updateQuery, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined, id)
	if err != nil {
		s.logger.Error("Error deleting user with id: %d", id)
		return err
	}

	return nil
}

// handles the http request to update the user
func (s *MySqlStore) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		s.logger.HttpError(w, http.StatusMethodNotAllowed, "Error, wrong request should be put")
		return
	}

	userIdStr := strings.TrimPrefix(r.URL.Path, "/users/update/")
	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		s.logger.HttpError(w, http.StatusBadRequest, "Error converting user id string to integer")
		return
	}

	var user *types.User
	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		s.logger.HttpError(w, http.StatusUnprocessableEntity, "Error decoding the JSON data")
		return
	}

	updatedUser := types.NewUser(user.ID, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined)

	s.UpdateUser(userID, updatedUser)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User data updated successfully!")
}

// gets all user rows from database
func (s *MySqlStore) GetAllUsers() ([]types.User, error) {

	getAllQuery := `
	SELECT * FROM user
	`
	var userData []types.User

	rows, err := s.db.Query(getAllQuery)
	if err != nil {
		s.logger.Error("Error querying all users, %v", err)
		return nil, err
	}
	//defer means it executes after the rest of the code finishes executing
	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Weight, &user.Goal, &user.Regimen, &user.DateJoined)
		if err != nil {
			return nil, err
		}

		userData = append(userData, user)
	}

	return userData, nil
}

// handles the http request to get all users
func (s *MySqlStore) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.HttpError(w, http.StatusMethodNotAllowed, "Error, invalid request to get all fighters")
		return
	}

	users, err := s.GetAllUsers()
	if err != nil {
		s.logger.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("error retreiving all users: %v", err))
		return
	}
	//convert the date strings to time objects
	for i := range users {
		dateJoined, err := time.Parse("2006-01-02", users[i].DateJoined)
		if err != nil {
			s.logger.HttpError(w, http.StatusInternalServerError, "Error converting date strings to time objects")
			return
		}
		users[i].DateJoined = dateJoined.Format("2006-01-02")
	}

	jsonResponse, err := json.Marshal(users)
	if err != nil {
		s.logger.HttpError(w, http.StatusInternalServerError, "Error marshalling all users")
		return
	}

	WriteJSONResponse(w, http.StatusOK, jsonResponse)
}
