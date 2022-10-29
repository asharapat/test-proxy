package store_impl

import (
	"database/sql"
	"personal/go-proxy-service/app/repository"
)

type Store struct {
	db *sql.DB
	taskRepository *TaskRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db:             db,
	}
}

func (st *Store) TaskRepository() repository.TaskRepoInt {
	if st.taskRepository != nil {
		return st.taskRepository
	} else {
		st.taskRepository = &TaskRepository{store: st}
	}
	return st.taskRepository
}

func (st *Store) BeginTx() (*sql.Tx, error) {
	return st.db.Begin()
}