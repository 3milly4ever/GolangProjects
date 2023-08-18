package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"user-fitness/logger"
	"user-fitness/store"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	listenAddr string
	//we inject the logger interface
	logger logger.Logger
	store  *store.MySqlStore
	db     *sql.DB
}

var Logger = logger.NewLogger()

// var Store = store.NewMySqlStore(Logger)

func NewServer(listenAddr string, store *store.MySqlStore) *Server {
	return &Server{
		listenAddr: listenAddr,
		// logger:     logger,
		store: store,
	}
}

func (s *Server) Connect(dataSourceName string) (*sql.DB, error) {
	var err error

	s.db, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println("Error connecting to the database", err)
		return nil, err
	}

	err = s.db.Ping()
	if err != nil {
		fmt.Println("error pinging the database: ", err)
		return nil, err
	}

	fmt.Println("Connection to the database successful")
	return s.db, nil
}

func (s *Server) GetDB() *sql.DB {
	return s.db
}

func (s *Server) CloseDB(logger logger.Logger) {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			s.logger.Error("Error closing the database", err)
		} else {
			fmt.Println("Database connection closed")
		}
	}
}

func (s *Server) HandleUserRequests(w http.ResponseWriter, r *http.Request) {
	sl := store.NewMySqlLogger(logger.NewLogger())
	switch r.Method {
	case "GET":
		if strings.HasPrefix(r.URL.Path, "/users/") {
			//if its only /users/ and nothing else then it will get all
			if r.URL.Path == "/users/" {
				s.store.HandleGetAllUsers(w, r, sl)
				//if the path doesnt exactly match /users/ but still starts with /users/
			} else if strings.HasPrefix(r.URL.Path, "/users/") {
				// Extract the user ID from the URL path
				parts := strings.Split(r.URL.Path, "/")
				//parts is a slice of strings that takes, for example, the localhost:9090/users/3 path.
				//then splits it according to the delimiter /. then forms an array such as parts = ["localhost:9090", "users", "3"].
				//the if statement checks if the length of the array is 3, and if index 2 of the parts array is not an empty string.
				//if thats true then our handlegetuserbyid is called.
				if len(parts) == 3 && parts[2] != "" {
					s.store.HandleGetUserById(w, r, sl)
				} else {
					http.Error(w, "Invalid user ID", http.StatusBadRequest)
				}
			} else {
				http.Error(w, "Invalid endpoint", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
		}
	case "POST":
		if r.URL.Path == "/users/" {
			s.store.HandleInsertUser(w, r, sl)
		} else {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
		}
	case "DELETE":
		if strings.HasPrefix(r.URL.Path, "/users/") {
			// Extract the user ID from the URL path
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) == 3 && parts[2] != "" {
				s.store.HandleDeleteUser(w, r, sl)
			} else {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
		}
	case "PUT":
		if strings.HasPrefix(r.URL.Path, "/users/") {
			// Extract the user ID from the URL path
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) == 3 && parts[2] != "" {
				s.store.HandleUpdateUser(w, r, sl)
			} else {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
		}
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
