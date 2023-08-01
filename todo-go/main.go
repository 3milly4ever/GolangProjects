package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// define struct
type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

// create a slice of structs and give each field a value
var todos = []todo{
	{ID: "1", Item: "Clean Room", Completed: false},
	{ID: "2", Item: "Read Book", Completed: false},
	{ID: "3", Item: "Record Video", Completed: false},
}

func addTodo(context *gin.Context) {

	var newTodo todo

	//error check and return if there is an error
	if err := context.BindJSON(&newTodo); err != nil {
		return
	}
	//executes if there is no error and adds the new to do to the todos slice of structs
	todos = append(todos, newTodo)
	//sends an HTTP response of a JSON representation of the todo item, the IndentedJSON serializes the newTodo data into JSON and sends it
	context.IndentedJSON(http.StatusCreated, newTodo)
}

func toggleTodoStatus(context *gin.Context) {

	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	//toggles
	todo.Completed = !todo.Completed

	context.IndentedJSON(http.StatusOK, todo)
}

func main() {

	//creates a new connection
	router := gin.Default() //check
	//gets all of the todos
	router.GET("/todos", getTodos)
	//gets a specific todo by id
	router.GET("/todos/:id", getTodo)
	//selects the todo by id, and toggles the completion status
	router.PATCH("/todos/:id", toggleTodoStatus)
	//adds a todo to the slice
	router.POST("/todos", addTodo)

	router.Run("localhost:9090") //check
}

//route handlers

func getTodo(context *gin.Context) {

	//extracts the value associated with the id parameter(for me 1, 2, or 3) and stores it in the id variable
	id := context.Param("id")
	//we call our getToDoById which gets all the values in the todo with the ID we choose, assigned to todo and to err
	todo, err := getTodoById(id)
	//this next line handles the error, if the todo with the id we provided does not exist a json message is sent, gin.H allows us to write the custom message, then program exits
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	//if the todo is found by the id the status
	context.IndentedJSON(http.StatusOK, todo)
}

// sends an http response to get all todo structs in the todos slice in JSON format
// handler function for a route that retrieves all the todos
func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

// takes in a string as parameter and returns a todo struct as a pointer that is found by the id and nil error or no pointer and error message if ID is wrong.
func getTodoById(id string) (*todo, error) {
	//goes through every todo struct in the todos slice and checks if the id passed to the function is a match, if it is it returns the todo at the id and returns the error as nil
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}

	//if the correct id is not found above, the todo is nil and an error is returned
	return nil, errors.New("todo not found")
}
