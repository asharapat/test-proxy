package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"personal/go-proxy-service/app/services"
)

type Server struct {
	Router      *mux.Router
	TaskService services.TaskServiceInt
}

func NewServer(service services.TaskServiceInt) *Server {
	srv := &Server{
		Router:      mux.NewRouter(),
		TaskService: service,
	}

	srv.configureRouter()

	return srv
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// allow for browsers to do CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, "+
		"Accept-Encoding, X-CSRF-Token, Authorization")
	srv.Router.ServeHTTP(w, r)
}

func (srv *Server) configureRouter() {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	cors := handlers.CORS(headers, methods, origins)

	srv.Router.Use(cors)

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	srv.Router.PathPrefix("/static/").Handler(s)

	srv.Router.Use(loggingMiddleware)

	TaskRoutes(srv)
}

type SimpleResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (srv *Server) error(w http.ResponseWriter, status int, err error) {
	srv.respond(w, status, &SimpleResponse{Code:"ERROR", Message: err.Error()})
}

func (srv *Server) respond(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
