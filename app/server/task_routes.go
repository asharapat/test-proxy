package server

import "net/http"

func TaskRoutes(srv *Server) {
	srv.Router.HandleFunc("/api/task/create", srv.MakeCreateTaskRequest()).
		Methods(http.MethodPost, http.MethodOptions)
	srv.Router.HandleFunc("/api/task/get/{id:[0-9]+}", srv.MakeGetTaskRequest()).
		Methods(http.MethodGet, http.MethodOptions)
}
