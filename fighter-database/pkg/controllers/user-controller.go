package controllers

import (
	"encoding/json"
	"fighter-database/pkg/models"
	f "fmt"
	"net/http"
	"strings"
	"time"
)

// get user by id handler
func HandleGetUserByIDRequest(w http.ResponseWriter, r *http.Request, userID int) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method, should be get", http.StatusMethodNotAllowed)
	}
	var err error
	// userIDStr := strings.TrimPrefix(r.URL.Path, "/users/id/")
	// userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user id: ", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserById(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Error retrieving data", http.StatusInternalServerError)
		return
	}

	dateJoined, err := time.Parse("2006-01-02", user.DateJoined)
	if err != nil {
		http.Error(w, "Error parsing into time object", http.StatusInternalServerError)
	}

	user.DateJoined = dateJoined.Format("2006-01-02") //the format has a time receiver and a return of type string

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// get all users handler
func HandleGetAllUsersRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method, should be get", http.StatusMethodNotAllowed)
		return
	}

	userRecords, err := models.GetAllUserData()
	if err != nil {
		http.Error(w, f.Sprintf("Error retrieving data: %v", err), http.StatusInternalServerError)
		return
	}
	//convert the DateJoined strings to time.Time objects
	for i := range userRecords {
		dateJoined, err := time.Parse("2006-01-02", userRecords[i].DateJoined)
		if err != nil {
			http.Error(w, "Error converting Date Joined strings to time.Time objects", http.StatusInternalServerError)
		}
		userRecords[i].DateJoined = dateJoined.Format("2006-01-02")
	}

	jsonResponse, err := json.Marshal(userRecords)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// create user
func HandleCreateUserRequest(w http.ResponseWriter, r *http.Request, userID int) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method, should be post", http.StatusMethodNotAllowed)
		return
	}

	var user models.UserRecord
	err := json.NewDecoder(r.Body).Decode(&user) //parsing the data from json to go into a UserRecord instance
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	newUser := models.NewUserRecord(user.ID, user.Weight, user.Goal, user.Regimen, user.DateJoined)

	models.InsertUserData(newUser)
	w.WriteHeader(http.StatusCreated)
	f.Fprintln(w, "Data received successfully!")
}

// update user
func HandleUpdateUserRequest(w http.ResponseWriter, r *http.Request, userID int) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method, should be put", http.StatusMethodNotAllowed)
		return
	}
	var err error
	// userIDStr := strings.TrimPrefix(r.URL.Path, "/users/update/")
	// userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user id: ", userID)
	}

	var user models.UserRecord
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	newUser := models.NewUserRecord(user.ID, user.Weight, user.Goal, user.Regimen, user.DateJoined)

	models.UpdateUserData(userID, newUser)
	w.WriteHeader(http.StatusCreated)
	f.Fprintln(w, "Data received successfully!")
}

// delete user
func HandleDeleteUserRequest(w http.ResponseWriter, r *http.Request, userID int) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method, should delete", http.StatusMethodNotAllowed)
		return
	}
	var err error
	// vars := mux.Vars(r)
	// userIDStr := vars["id"]                // get id from url path
	// userID, err := strconv.Atoi(userIDStr) // convert id to int for our Go functions
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	err = models.DeleteUserData(userID)
	if err != nil {
		http.Error(w, "Error deleting user data", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	f.Fprintln(w, "User data deleted successfully!")
}
