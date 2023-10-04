package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-fitness/caching"
	"user-fitness/logger"
	"user-fitness/types"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlStore struct {
	db    *sql.DB
	cache caching.Cache
}

// this makes sense if we have more fields in the MySqlStore struct.
func NewMySqlStore(db *sql.DB, cache caching.Cache) *MySqlStore {
	return &MySqlStore{
		db,
		cache,
	}
}

type StoreWithCache struct {
	*MySqlStore
	cache caching.Cache
}

func NewStoreWithCache(s *MySqlStore, cache caching.Cache) *StoreWithCache {
	return &StoreWithCache{
		MySqlStore: s,
		cache:      cache,
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
	GetUserByIdFromDB(id int, sl *SqlLogger) (types.User, error)
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

func (s *MySqlStore) GetUserByIdWithCache(id int, sl *SqlLogger) (types.User, error) {
	sc := NewStoreWithCache(s, s.cache)
	return sc.GetUserById(id, sl)
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
func (s *MySqlStore) GetUserByIdFromDB(id int, sl *SqlLogger) (types.User, error) {
	getUserByIdQuery := `
	SELECT * FROM Users
	WHERE id = ?
	`
	sl.Logger.Info("Getting user from database")
	//create new instance of a type that satisfies the Logger interface.
	row := s.db.QueryRow(getUserByIdQuery, id)

	var user types.User

	sl.Logger.Info("Fetching user with ID: %d", id)

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
	sl.Logger.Info("Fetched user successfully: %v", user)
	return user, nil
}

// this get user by ID incorporates caching
func (sc *StoreWithCache) GetUserById(id int, sl *SqlLogger) (types.User, error) {
	//recover from panic and log the error
	defer func() {
		if r := recover(); r != nil {
			sl.Logger.Error("Panic occurred: %v", r)
		}
	}()
	//we create a cache key that we will use to either Get or Set from our redis cache
	cacheKey := fmt.Sprintf("user:%d", id)

	//Try to get user from cache
	cachedUserBytes, err := sc.cache.Get(cacheKey)
	//if there is no error getting the cached data with the key, and the cachedUserBytes are empty
	if err == nil && cachedUserBytes != nil {
		//create a cachedUser type
		var cachedUser types.User
		//if there is an error unmarshaling the data from the cache in bytes, and assigning it to our cachedUser variable
		if err := json.Unmarshal(cachedUserBytes, &cachedUser);
		//if there is no error
		err == nil {
			//we return the newly created cachedUser data with the now bytes data converted(unmarshaled) from the cache.
			sl.Logger.Info("Found cached user data for ID %d", id)
			return cachedUser, nil
		} else {
			//if there is an error, and the user is not in the then we do this.
			sl.Logger.Info("User not in cache: %v, %d", cachedUserBytes, cachedUser.ID)
			sl.Logger.Error("Error unmarshaling cached user data:", err)
		}

	}

	//fetch from database if user is not in cache
	user, err := sc.MySqlStore.GetUserByIdFromDB(id, sl)
	if err != nil {
		//if there is an error fetching from database
		sl.Logger.Error("Error fetching user from the database:", err)
		return types.User{}, err
	}

	//cache the retrieved user so in the future its data can be accessed from the cache
	userJSON, err := json.Marshal(user)
	//the user is being marshaled into a JSON byte array and set to the userJSON variable
	if err == nil {
		//if there is no error we will also cache this user.
		sl.Logger.Info("User JSON data being cached: %s", userJSON)
		//we put the user into the cache with the cacheKey we created earlier, for an hour.
		err = sc.cache.Set(cacheKey, userJSON, time.Hour)
		if err != nil {
			//if there is an error setting the cache
			sl.Logger.Error("Error caching user data:", err)
			sl.Logger.Error("Failed to create user with ID %d: %v", id, err)
		} else {
			//no error
			sl.Logger.Info("Successfully cached user data")
		}
		//if there is an error marshaling the data
	} else {
		sl.Logger.Error("Error marshaling user data:", err)
	}
	//we return the user to the handler
	return user, nil
}

// handles the GET http request to get user by id
func (s *MySqlStore) HandleGetUserById(w http.ResponseWriter, r *http.Request, sl *SqlLogger) {
	//we get the id from the route path
	userIDStr := strings.TrimPrefix(r.URL.Path, "/users/")
	//we change the id to int
	userID, err := strconv.Atoi(userIDStr)
	//if the id is incorrect, or doesn't match the id we set in the route
	if err != nil {
		sl.Logger.HttpError(w, http.StatusBadRequest, "Invalid user id")
		return //return so the function ends after an error
	}

	//we attempt to get user from cache
	cachedUser, err := s.GetUserByIdWithCache(userID, sl)
	if err != nil {
		//if there is no
		sl.Logger.Error("Error getting user from cache:", err)
		return
	}

	//if cachedUser ID is not 0, this means the cashed user has an id
	if cachedUser.ID != 0 {
		user := cachedUser
		// User found in cache, use cached data

		// We parse(convert) the user.DateJoined string into a time value and assign it to the dateJoined reference
		dateJoined, err := time.Parse("2006-01-02", user.DateJoined)
		if err != nil {
			//if there is an error we just want to see what format it is
			//to have a better understanding
			sl.Logger.Info("date is in format: %v", user.DateJoined)
			//we write there was an error parsing into the time object
			sl.Logger.HttpError(w, http.StatusInternalServerError, "Error parsing into time object")
			return
		}

		//we format it in the right format and change it back to string to then marshal
		//the whole reason we converted to a time value was to be able to call .Format it right before we marshal it
		user.DateJoined = dateJoined.Format("2006-01-02") // the format has a time receiver and a return of type string
		jsonResponse, err := json.Marshal(user)           //string format can now be marshaled
		if err != nil {
			sl.Logger.HttpError(w, http.StatusInternalServerError, "Error creating JSON response for user by id")
			return
		}

		// Write the JSON response
		WriteJSONResponse(w, http.StatusOK, jsonResponse)
		return
	}

	//if the user is not in cache the following executes, we get the user from the database instead of the cache
	user, err := s.GetUserByIdFromDB(userID, sl)
	//we check what format the date is i
	sl.Logger.Info("user from database check date is in format: %v", user.DateJoined)
	if err != nil {
		//if the error is that no rows were found in the sql database, we return the following HttpError
		//meaning that the id was not found
		if err == sql.ErrNoRows {
			sl.Logger.HttpError(w, http.StatusNotFound, err.Error())
			return
		}
		//if there was a different error, not the sql.ErrNoRows then this will be thrown
		sl.Logger.HttpError(w, http.StatusNotFound, "Error retrieving user by id")
		return
	}

	//we parse(convert) the dateJoined string into a time value, to then format it correctly
	dateJoined, err := time.Parse("2006-01-02", user.DateJoined)
	//if there is an error during the parsing from string to time.Time
	if err != nil {
		sl.Logger.Info("date is in format: %v", user.DateJoined)
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error parsing into time object")
		return
	}
	//we parse the time.Time back to string with the correct format to then marshal it
	user.DateJoined = dateJoined.Format("2006-01-02") //the format has a time receiver and a return of type string
	jsonResponse, err := json.Marshal(user)           //we marshal it for the jsonResponse
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
	//getting query parameters for page and pageSize
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	//Parse page and pageSize string values to int or default values if parsing fails
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		//if parsing fails page is set to default value
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		//if parsing fails pageSize is set to default value
		pageSize = 10
	}

	//Calculate the offset based on page and pageSize
	//so if we are on page 4 and pageSize(how many records will be shown on a page) is 10, offset would be (4-1) * 10
	//which would equal to 30, which means the first 30 records will be skipped and we will begin seeing them from record 31
	offset := (page - 1) * pageSize

	//Call GetUsersWithPagination to retrieve users for the requested page
	users, totalUsers, err := s.GetUsersWithPagination(offset, pageSize, sl)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error retrieving users")
		return
	}

	//calculate total pages
	totalPages := int(math.Ceil(float64(totalUsers) / float64(pageSize)))

	//create pagination links
	baseURL := "/users/"
	links := types.CreatePaginationLinks(baseURL, page, pageSize, totalUsers)

	//construct the paginated response
	response := types.PaginatedUserResponse{
		Users:       users,
		TotalUsers:  totalUsers,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
		Links:       links,
	}
	//we range through the response, which has the 10 users based on our pageSize and do the time formatting here too
	//first from string to time.Time, then format correctly back into string
	for i := range response.Users {
		dateJoined, err := time.Parse("2006-01-02", response.Users[i].DateJoined)
		if err != nil {
			sl.Logger.HttpError(w, http.StatusInternalServerError, "Error converting date strings to time objects")
			return
		}
		response.Users[i].DateJoined = dateJoined.Format("2006-01-02")
	}
	//marshal and write the response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		sl.Logger.HttpError(w, http.StatusInternalServerError, "Error marshalling all users")
		return
	}

	WriteJSONResponse(w, http.StatusOK, jsonResponse)
}

func (s *MySqlStore) GetUsersWithPagination(offset, pageSize int, sl *SqlLogger) ([]types.User, int, error) {
	//get total number of users in the database
	totalUsers, err := s.GetTotalUsers()
	if err != nil {
		return nil, 0, err
	}

	//now we get the users for the specified page using LIMIT and OFFSET
	getUsersQuery := `
	SELECT * FROM Users
	LIMIT ? OFFSET ?
	`
	rows, err := s.db.Query(getUsersQuery, pageSize, offset)
	if err != nil {
		sl.Logger.Error("Error querying users: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	//we iterate through the rows and parse data
	var users []types.User
	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Weight, &user.Goal, &user.Regimen, &user.DateJoined)
		if err != nil {
			return nil, 0, err
		}

		links := types.CreateUserHypermediaLinks(user.ID)
		user.Links = links
		users = append(users, user)
	}

	return users, totalUsers, nil
}

func (sc *StoreWithCache) GetUsersWithPaginationCached(offset, pageSize int, sl *SqlLogger) ([]types.User, int, error) {
	//we create cache key based off of offset and pagSize
	cacheKey := fmt.Sprintf("users:%d:%d", offset, pageSize)

	//1.try to get paginated user data with cache
	cachedUsersJSON, err := sc.cache.Get(cacheKey)
	if err == nil && cachedUsersJSON != nil {
		//If cached data exists, we unmarshal and return it
		var cachedUsers []types.User
		if err := json.Unmarshal([]byte(cachedUsersJSON), &cachedUsers); err == nil {
			return cachedUsers, len(cachedUsers), nil
		}
	}
	//2.if user data not in cache then get it from database
	users, totalUsers, err := sc.MySqlStore.GetUsersWithPagination(offset, pageSize, sl)
	if err != nil {
		return nil, 0, err
	}
	//3. if 2 then cache retrieved users data for future access
	usersJSON, err := json.Marshal(users)
	if err == nil {
		err = sc.cache.Set(cacheKey, usersJSON, time.Hour)
		if err != nil {
			sl.Logger.Error("Failed to cache paginated users data: %v", err)
		}
	}

	return users, totalUsers, nil
}

func (s *MySqlStore) GetTotalUsers() (int, error) {
	//this will return the number of total users from our database
	query := "SELECT COUNT(*) FROM Users"
	var totalUsers int
	err := s.db.QueryRow(query).Scan(&totalUsers)
	if err != nil {
		return 0, err
	}
	return totalUsers, nil
}
