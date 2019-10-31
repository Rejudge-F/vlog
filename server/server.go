package server

import (
	"net/http"
	"vlog/controllers"
)

// iServer the interface for server
type iServer interface {
	Run(addr string)
}

// Server the instance for a server
type Server struct {
	Controller *controllers.MainController
}

// Run
func (server *Server) Run(addr string) {
	server.Controller.Initialize()
	err := http.ListenAndServe(addr, server.Controller.Router)
	if err != nil {
		panic("http listen failed")
	}
}
