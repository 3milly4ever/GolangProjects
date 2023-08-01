package routes

import (
	"book-management-app/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	//creates a new book posts it
	router.HandleFunc("/book/", controllers.CreateBook).Methods("POST")
	//creates book by id
	router.HandleFunc("/book/{bookId}", controllers.CreateBookWithIDHandler).Methods("POST")
	//gets all books
	router.HandleFunc("/book/", controllers.GetBook).Methods("GET")
	//we get a specific book by ID
	router.HandleFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
	//we update an existing books fields
	router.HandleFunc("/book/{bookId}", controllers.UpdateBook).Methods("PUT")
	//we delete a book
	router.HandleFunc("/book/{bookId}", controllers.DeleteBook).Methods("DELETE")
	//delete all books
	router.HandleFunc("/book/", controllers.DeleteAllBooks).Methods("DELETE")

}
