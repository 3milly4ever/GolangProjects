package models

import (
	"book-management-app/pkg/config"
	"errors"

	"github.com/jinzhu/gorm"
)

// db variable will represent a database connection in the gorm library
// we use it to talk to the database
var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:""json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	//calling the connect function from the config package which opens/establishes a connection with mysql using gorm
	config.Connect()
	//puts the data base in a variable for easy access, represents connection to the database. can be used for executing queries and database operations.
	db = config.GetDB()
	//the automigrate method automatically creates or modifies database tables based on the provided model structure.
	db.AutoMigrate(&Book{})
}

//create all functions that we need to communicate with the database

func (b *Book) CreateBook() *Book {
	//NewRecord is used to check if the book object b is a new record in the database.
	//the method returns true if the record does not exist in the database yet
	db.NewRecord(b)
	//Create creates a new record in the database using the b object
	db.Create(&b)
	//returns the modified "Book" object after it has been created in the database
	return b
}

// dont need just what I added
func (b *Book) CreateBookWithID(ID int64) (*Book, error) {

	// Check if a book with the specified ID already exists
	var existingBook Book
	result := db.First(&existingBook, ID)
	if result.RowsAffected != 0 {
		return nil, errors.New("book with specified ID already exists")
	}

	// Set the specified ID to the book object
	b.ID = uint(ID)

	// Create a new record in the database using the b object
	err := db.Create(&b).Error
	if err != nil {
		return nil, err
	}

	return b, nil
}

// function that will return all book models from the database and place it in the Go slice we create in this function
func GetAllBooks() []Book {
	var Books []Book //create a slice of book structs
	db.Find(&Books)  //fetch all books from the database and store them in the Books slice. the &Books passes the referene to the slice and the Find method populates it
	return Books     //returns the slice of the populated books.
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	//this variable will store the book with the particular id that we want to return
	var getBook Book
	//the db variable is initialized to find the book by the id
	db := db.Where("ID=?", Id).Find(&getBook) //the where and the find methods from the GORM library communicate with the SQL database
	//we return the particular book and the database connection(*gorm.DB)
	return &getBook, db
}

func DeleteAllBooks() error {

	if err := db.Delete(&Book{}).Error; err != nil {
		return err
	}

	return nil
}

func DeleteBook(ID int64) Book {
	//variable for the particular book that we want to delete
	var book Book
	//database connection communicates with mysql to delete the book with the id that is passed into this function
	db.Where("ID=?", ID).Delete(&book)
	//
	return book
}

//second version of this function
// func DeleteBook(ID int64) error {
// 	var book Book

// 	// Find the book by ID
// 	if err := db.First(&book, ID).Error; err != nil {
// 		return err
// 	}

// 	// Delete the book from the database
// 	if err := db.Delete(&book).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }
