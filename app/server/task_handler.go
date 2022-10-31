package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"personal/go-proxy-service/pkg/models"
	"strconv"
)

func (srv *Server) MakeCreateTaskRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := &models.Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			srv.error(w, http.StatusBadRequest, err)
			return
		}
		res, err := srv.TaskService.CreateTask(task)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err)
			return
		} else {
			srv.respond(w, http.StatusOK, res)
		}
	}
}

func (srv *Server) MakeGetTaskRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		res, err := srv.TaskService.GetTask(id)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err)
		} else {
			srv.respond(w, http.StatusOK, res)
		}
	}
}
