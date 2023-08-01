package main

//we will create server which will also define our local host
//

import (
	"book-management-app/pkg/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	//r variable initializes new mux router
	r := mux.NewRouter()
	//calling the function from routes which we will pass our r variable to.
	routes.RegisterBookStoreRoutes(r)
	//registers handler for this specific path. has to be registered before it can be used to handle requests.
	http.Handle("/", r)
	//listen and serve starts the server on the specified url, will be handled by r at the specified address localhost:9090
	//if there is an error log.Fatal will log the error message and terminate the program
	//combining these two statements an http server will start and immediately checks for errors.
	//convenient way to handle start up errors without explicitly handling them in separate if statements
	log.Fatal(http.ListenAndServe("localhost:9090", r))
}
