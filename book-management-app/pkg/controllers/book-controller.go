package controllers

//handles http requests
//and based on the requests calls the corresponding function from the models pacakges
//Pointers when:
//When we want to update state
//When we want to optimizie memory for large objects that are getting called a lot

import (
	"book-management-app/pkg/models"
	"book-management-app/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// the models refers to the models folder(package) and the Book refers to the book struct, we are creating a new book and storing it in this NewBook reference
var NewBook models.Book

// when you are building a controller function you have two things, a response and a request.
// for the request you put a pointer for the request that you receive from the user.
func GetBook(w http.ResponseWriter, r *http.Request) {
	//the GetAllBooks function from the models package will get a list of the books for the user
	//we will have a json response to convert it into json using the marshal function.
	newBooks := models.GetAllBooks()

	//we are discarding the error value, so it is not being handled or utilized in this context
	res, _ := json.Marshal(newBooks) //this is a write only variable. the _ is where an error would be. the _ is declared because a Marshal function returns two values, the marshalled json data and an error.

	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK) //will write 200 if successful
	//the response will be the json version of the new books data that will be found in the mysql database
	w.Write(res)
}

// the request will be to find the book by bookid
func GetBookById(w http.ResponseWriter, r *http.Request) {
	//the max Vars function extracts the url parameter from /books/{bookId}
	vars := mux.Vars(r)      //using the mux router, we retrieve the value of the bookId variable from the URL
	bookId := vars["bookId"] //we assign the value we retrieved in the above like to the bookId variable here.
	//the 0, 0 represents the base and bit size of the resulting integer value.
	ID, err := strconv.ParseInt(bookId, 0, 0) //we convert the bookId string to an integer using strconv.ParseInt. it returns the passed integer Id if correct and an error if there is one.
	if err != nil {
		fmt.Println("error while parsing") //if error while parsing we return this message
	}
	bookDetails, _ := models.GetBookById(ID)           //blank variable because we don't need the db returned right now.
	res, _ := json.Marshal(bookDetails)                //sending a JSON response to the user, storing it in res
	w.Header().Set("Content-Type", "pkglication/json") //Informs the client that the response body will be in JSON format
	w.WriteHeader(http.StatusOK)                       //Sets the status code of the response to 200 (OK) to indicate a successful response
	w.Write(res)                                       //writes the JSON response to the writer, which sends it back to the client(us).
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	//new instance of the models book struct the ampersand(pointer) is used to get the memory address of the newly created CreateBook object
	CreateBook := &models.Book{} //1. first we receive json.
	//the parse body function is called to request body and populate the CreateBook object with data.
	//CreateBook is the target object to unmarshall data into. and the r parameter is the http request that we are making.
	utils.ParseBody(r, CreateBook) //2. then we parse into something golang will understand
	//first CreateBook is the variable we declared above and the second CreateBook is the function from the models package
	b := CreateBook.CreateBook() //3. we save the record to the database
	res, _ := json.Marshal(b)    //4. we are converting the record to json to send it to postman/user
	w.WriteHeader(http.StatusOK) //send status 200
	w.Write(res)                 //sends the JSON data to the client/postman/user
}

func CreateBookWithIDHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	utils.ParseBody(r, &book) //parse into something the database will understand

	//Get ID from request or query parameter
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	//If no errors are found this code will execute, creates the book with the specified ID
	createdBook, err := book.CreateBookWithID(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Convert the created book to JSON
	responseData, err := json.Marshal(createdBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set the response content type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

// This function receives the request from the client, extracts the book ID, calls the DeleteBooks function from the models package
// then constructs the response
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing") //if error while parsing we return this message
	}
	book := models.DeleteBook(ID)
	res, _ := json.Marshal(book)
	w.Header().Set("Content-Type", "pkglication/json")
	w.Write(res)
}

func DeleteAllBooks(w http.ResponseWriter, r *http.Request) {
	err := models.DeleteAllBooks()
	if err != nil {
		//error handling
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All books have been deleted."))
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook = &models.Book{}
	utils.ParseBody(r, updateBook)
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}

	bookDetails, db := models.GetBookById(ID) //we find the book in the database using our GetBookById function from the models package
	//if the name, author, or publication are not empty it will update them
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	//once we update the book we want to save the update
	db.Save(&bookDetails)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
