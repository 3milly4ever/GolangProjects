package routes

import (
	"fighter-database/pkg/controllers"
	"net/http"
)

var RegisterUserRoutes = func() {
	//we set up the http server and handle requests
	http.HandleFunc("/users/", controllers.HandleCreateUserRequest)
	http.HandleFunc("/users/delete/{id}", controllers.HandleDeleteUserRequest)
	http.HandleFunc("/users/update/{id}", controllers.HandleUpdateUserRequest)
	http.HandleFunc("/users/all", controllers.HandleGetAllUsersRequest)
	http.HandleFunc("/users/{id}", controllers.HandleGetUserByIDRequest)
}
