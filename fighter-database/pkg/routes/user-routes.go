package routes

import (
	"fighter-database/pkg/controllers"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// we set up the http server and handle requests
// http.HandleFunc("/users/delete/{id}", controllers.HandleDeleteUserRequest)
// http.HandleFunc("/users/update/{id}", controllers.HandleUpdateUserRequest)
// http.HandleFunc("/users/all", controllers.HandleGetAllUsersRequest)
// http.HandleFunc("/users/id/{id}", controllers.HandleGetUserByIDRequest)
// http.HandleFunc("/users/", controllers.HandleCreateUserRequest)
var RegisterUserRoutes = func() {
	// Set up the http server and handle requests
	http.HandleFunc("/users/delete/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			vars := mux.Vars(r)
			userIDStr := vars["id"]                // get id from url path
			userID, err := strconv.Atoi(userIDStr) // convert id to int for our Go functions
			if err != nil {
				http.Error(w, "Invalid userID", http.StatusBadRequest)
				return
			}

			controllers.HandleDeleteUserRequest(w, r, userID) // Call the corresponding function
		} else {
			http.Error(w, "Invalid request method, should be delete", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/update/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			vars := mux.Vars(r)
			userIDStr := vars["id"]                // get id from url path
			userID, err := strconv.Atoi(userIDStr) // convert id to int for our Go functions
			if err != nil {
				http.Error(w, "Invalid userID", http.StatusBadRequest)
				return
			}

			controllers.HandleUpdateUserRequest(w, r, userID) // Call the corresponding function
		} else {
			http.Error(w, "Invalid request method, should be put", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.HandleGetAllUsersRequest(w, r) // Call the corresponding function
		} else {
			http.Error(w, "Invalid request method, should be get", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/id/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			vars := mux.Vars(r)
			userIDStr := vars["id"]                // get id from url path
			userID, err := strconv.Atoi(userIDStr) // convert id to int for our Go functions
			if err != nil {
				http.Error(w, "Invalid userID", http.StatusBadRequest)
				return
			}

			controllers.HandleGetUserByIDRequest(w, r, userID) // Call the corresponding function
		} else {
			http.Error(w, "Invalid request method, should be get", http.StatusMethodNotAllowed)
		}
	})
}
