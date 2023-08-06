package api

import (
	"encoding/json"
	"net/http"
	"structure/store"
)

type Server struct {
	listenAddr string
	store      store.Store
}

func NewServer(listenAddr string, store store.Store) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

// you can make a separate file in this api package for all handlers if there's a lot, or for easy access in general. but don't make a bunch of folders.
func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	user := s.store.Get(10)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) Start() error {
	http.HandleFunc("/user/id", s.handleGetUserByID)
	return http.ListenAndServe(s.listenAddr, nil)
}
