package repository

import "database/sql"

type StoreInt interface {
	TaskRepository() TaskRepoInt
	BeginTx() (*sql.Tx, error)
}
