package api

import "net/http"

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) Start() error {
	http.HandleFunc("/user/id", s.handleGetUserByID)
	return http.ListenAndServe(s.listenAddr, nil)
}
