package main

import (
	"encoding/json"
	f "fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

//models & controllers packages

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies) //encoding the response into json
}

func deleteMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {

		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) //so the movies[:index] is the current one we receive through the client's request,
			// and it is being replaced by the rest of the movies slice which basically deletes it
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //the r is what is specified by the request through postman(or user), and we want the request to match up the json id!
	for _, item := range movies {
		if item.ID == params["id"] { //!which is why we put id here
			json.NewEncoder(w).Encode(item)
			f.Printf("The selected movie is: ", item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //we set the type of the request
	newMovie := Movie{}                                //create a new moivie struct
	_ = json.NewDecoder(r.Body).Decode(&newMovie)      //we are decoding the json which the request is sent in to convert it to golang and store the data in our program
	newMovie.ID = strconv.Itoa(rand.Intn(100000000))   //gives us a random number between 1 and 10000000 for the ID field of the new movie
	movies = append(movies, newMovie)                  //we store the newMovie in the movies slice
	json.NewEncoder(w).Encode(newMovie)                //we respond to the user with the details of the newly created movie
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	updatedMovie := Movie{}
	params := mux.Vars(r)
	_ = json.NewDecoder(r.Body).Decode(&updatedMovie)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			if movie.Director != nil {
				movie.Director = updatedMovie.Director
			}
			if movie.Isbn != "" {
				movie.Isbn = updatedMovie.Isbn
			}
			if movie.Title != "" {
				movie.Title = updatedMovie.Title
			}
			json.NewEncoder(w).Encode(movie)
		}
	}

}

func main() {
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie 1", Director: &Director{Firstname: "Jon", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "234551", Title: "Movie 2", Director: &Director{Firstname: "Don", Lastname: "Tarantino"}})
	//would be in router package
	r := mux.NewRouter()
	r.HandleFunc("/movies", getAllMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovieById).Methods("DELETE")

	f.Printf("Starting server at port 9090\n")
	log.Fatal(http.ListenAndServe(":9090", r))
}
