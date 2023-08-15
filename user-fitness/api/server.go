package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"user-fitness/logger"
	"user-fitness/store"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	listenAddr string
	//we inject the logger interface
	logger logger.Logger
	store  store.Store
	db     *sql.DB
}

var Logger = logger.NewLogger()

// var Store = store.NewMySqlStore(Logger)

func NewServer(listenAddr string, logger logger.Logger, store store.Store) *Server {
	return &Server{
		listenAddr: listenAddr,
		logger:     logger,
		store:      store,
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

func (s *Server) RegisterUserRoutes() {
	http.HandleFunc("/users/", s.store.HandleInsertUser)
	http.HandleFunc("/users/delete/{id:[0-9]+}", s.store.HandleDeleteUser)
	http.HandleFunc("/users/update/{id:[0-9]+}", s.store.HandleUpdateUser)
	http.HandleFunc("/users/all", s.store.HandleGetAllUsers)
	http.HandleFunc("/users/{id:[0-9]+}", s.store.HandleGetUserById)
}
