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
	db *sql.DB
}

// this makes sense if we have more fields in the MySqlStore struct.
func NewMySqlStore(db *sql.DB) *MySqlStore {
	return &MySqlStore{
		db,
	}
}

type Store interface {
	HandleInsertUser(w http.ResponseWriter, r *http.Request, sl *SqlLogger)
	HandleDeleteUser(w http.ResponseWriter, r *http.Request, sl *SqlLogger)
	HandleUpdateUser(w http.ResponseWriter, r *http.Request, sl *SqlLogger)
	HandleGetAllUsers(w http.ResponseWriter, r *http.Request, sl *SqlLogger)
	HandleGetUserById(w http.ResponseWriter, r *http.Request, sl *SqlLogger)
	CreateTableWithFields(tableName string, fields string) error
	CreateTables(tables []TableDefinition, sl *SqlLogger) error
	NewMySqlLogger() SqlLogger
}

type SqlLogger struct {
	Logger logger.Logger
}

func NewMySqlLogger(logger logger.Logger) *SqlLogger {
	return &SqlLogger{
		logger,
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

func (s *MySqlStore) CreateTables(tables []TableDefinition, sl *SqlLogger) error {
	for _, table := range tables {
		err := s.CreateTableWithFields(table.Name, table.Fields)
		if err != nil {
			sl.Logger.Error("Error creating table %s: %v", table.Name, err)
			return err
		} else {
			sl.Logger.Info("Table %s created successfully", table.Name)
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
func (s *MySqlStore) GetUserById(id int, sl *SqlLogger) (types.User, error) {
	getUserByIdQuery := `
	SELECT * FROM Users
	WHERE id = ?
	`
	//create new instance of a type that satisfies the Logger interface.
	row := s.db.QueryRow(getUserByIdQuery, id)

	var user types.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Weight, &user.Goal, &user.Regimen, &user.DateJoined)

	if err != nil {
		if err == sql.ErrNoRows {
			sl.Logger.Error("user with id %d not found", id)
			return types.User{}, err //returns empty user
		}
		sl.Logger.Error("database error: %v", err)
		return types.User{}, err //returns empty user
	}

	links := types.CreateUserHypermediaLinks(user.ID)
	user.Links = links
	return user, nil
}

// handles the GET http request to get user by id
func (s *MySqlStore) HandleGetUserById(w http.ResponseWriter, r *http.Request, sl *SqlLogger) {
	if r.Method != http.MethodGet {
		sl.Logger.HttpError(w, http.StatusMethodNotAllowed, "Invalid method request, should be GET")
		return
	}

	userIDStr := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusBadRequest, "Invalid user id")
		return //return so the function ends after an error
	}

	user, err := s.GetUserById(userID, sl)
	if err != nil {
		if err == sql.ErrNoRows {
			sl.Logger.HttpError(w, http.StatusNotFound, err.Error())
			return
		}
		sl.Logger.HttpError(w, http.StatusNotFound, "Error retrieving user by id")
		return
	}

	//we parse(convert) the dateJoined string into a time value
	dateJoined, err := time.Parse("2006-01-02", user.DateJoined)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error parsing into time object")
		return
	}

	user.DateJoined = dateJoined.Format("2006-01-02") //the format has a time receiver and a return of type string
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error creating JSON response for user by id")
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(jsonResponse)
	WriteJSONResponse(w, http.StatusOK, jsonResponse)
}

// will create/insert a new user into the database
func (s *MySqlStore) InsertUser(user *types.User, sl *SqlLogger) (int64, error) {
	insertQuery := `
	INSERT INTO Users (name, email, weight, goal, regimen, date_joined)
	VALUES(?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(insertQuery, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined)
	//error check and provide user name
	if err != nil {
		sl.Logger.Error("Error inserting user %s: %v", user.Name, err)
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		sl.Logger.Error(" \" Error getting last inserted ID: %v", err)
		return 0, err
	}

	sl.Logger.Info("Successfully inserted! User details: Name=%s, Email=%s, Weight=%d, Goal=%s, Regimen=%s, DateJoined=%s", user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined)
	return lastInsertID, nil
}

// will handle the http request to create/insert a new user
func (s *MySqlStore) HandleInsertUser(w http.ResponseWriter, r *http.Request, sl *SqlLogger) {
	if r.Method != http.MethodPost {
		sl.Logger.HttpError(w, http.StatusMethodNotAllowed, "Invalid request method, should be post")
		return
	}
	//
	var user types.User                          //we will decode the json data into this
	err := json.NewDecoder(r.Body).Decode(&user) //we parse the data the client sent in json to create new user
	if err != nil {
		sl.Logger.HttpError(w, http.StatusBadRequest, "Invalid JSON data during post")
		return
	}

	links := types.CreateUserHypermediaLinks(0)
	//now that we decoded the body into user we can create a new user with the credentials
	newUser := types.NewUser(user.ID, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined, links)
	insertedID, err := s.InsertUser(newUser, sl)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error inserting user")
		return
	}
	links = types.CreateUserHypermediaLinks(insertedID)
	userWithLinks := types.NewUser(insertedID, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined, links)
	w.WriteHeader(http.StatusCreated)
	jsonResponse, _ := json.Marshal(userWithLinks)
	w.Write(jsonResponse)
	fmt.Fprintln(w, "Data posted successfuly!")
}

// delete a user from the database
func (s *MySqlStore) DeleteUser(id int, sl *SqlLogger) error {
	deleteQuery := `
	DELETE FROM Users
	WHERE id = ?
	`

	_, err := s.db.Exec(deleteQuery, id)
	if err != nil {
		sl.Logger.Error("Error deleting user with id: %d", id)
		return err
	}
	return nil
}

// handles http request for deleting user from database
func (s *MySqlStore) HandleDeleteUser(w http.ResponseWriter, r *http.Request, sl *SqlLogger) {
	if r.Method != http.MethodDelete {
		sl.Logger.HttpError(w, http.StatusBadRequest, "Error, wrong request should be DELETE")
		return
	}

	userIDStr := strings.TrimPrefix(r.URL.Path, "/users/") //get id from url path
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = s.DeleteUser(userID, sl)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error deleting user data")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User data deleted successfuly!")
}

// updates user in our database
func (s *MySqlStore) UpdateUser(id int, user *types.User, sl *SqlLogger) error {
	updateQuery := `
	UPDATE Users
	SET name = ?, email = ?, weight = ?, goal = ?, regimen = ?, date_joined = ?
	WHERE id = ?
	`

	_, err := s.db.Exec(updateQuery, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined, id)
	if err != nil {
		sl.Logger.Error("Error updating user with id: %d", id)
		return err
	}

	return nil
}

// handles the http request to update the user
func (s *MySqlStore) HandleUpdateUser(w http.ResponseWriter, r *http.Request, sl *SqlLogger) {
	if r.Method != http.MethodPut {
		sl.Logger.HttpError(w, http.StatusMethodNotAllowed, "Error, wrong request should be put")
		return
	}

	userIdStr := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusBadRequest, "Error converting user id string to integer")
		return
	}

	var user types.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusUnprocessableEntity, "Error decoding the JSON data")
		return
	}
	links := types.CreateUserHypermediaLinks(user.ID)
	updatedUser := types.NewUser(user.ID, user.Name, user.Email, user.Weight, user.Goal, user.Regimen, user.DateJoined, links)

	s.UpdateUser(userID, updatedUser, sl)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User data updated successfully!")
}

// gets all user rows from database
func (s *MySqlStore) GetAllUsers(sl *SqlLogger) ([]types.User, error) {

	getAllQuery := `
	SELECT * FROM Users
	`
	var userData []types.User

	rows, err := s.db.Query(getAllQuery)
	if err != nil {
		sl.Logger.Error("Error querying all users, %v", err)
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

		links := types.CreateUserHypermediaLinks(user.ID)
		user.Links = links
		userData = append(userData, user)
	}

	return userData, nil
}

// handles the http request to get all users
func (s *MySqlStore) HandleGetAllUsers(w http.ResponseWriter, r *http.Request, sl *SqlLogger) {
	if r.Method != http.MethodGet {
		sl.Logger.HttpError(w, http.StatusMethodNotAllowed, "Error, invalid request to get all fighters")
		return
	}

	users, err := s.GetAllUsers(sl)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, fmt.Sprintf("error retreiving all users: %v", err))
		return
	}
	//convert the date strings to time objects
	for i := range users {
		dateJoined, err := time.Parse("2006-01-02", users[i].DateJoined)
		if err != nil {
			sl.Logger.HttpError(w, http.StatusInternalServerError, "Error converting date strings to time objects")
			return
		}
		users[i].DateJoined = dateJoined.Format("2006-01-02")
	}

	jsonResponse, err := json.Marshal(users)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error marshalling all users")
		return
	}

	WriteJSONResponse(w, http.StatusOK, jsonResponse)
}
