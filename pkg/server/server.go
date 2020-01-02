package server

import "net/http"

// Server ...
type Server struct {
	address string
	mux     *http.ServeMux
}

// New ...
func New(address string, mux *http.ServeMux) *Server {
	return &Server{
		address: address,
		mux:     mux,
	}
}

// Start server
func (s *Server) Start() error {
	return http.ListenAndServe(":8080", s.mux)
}
