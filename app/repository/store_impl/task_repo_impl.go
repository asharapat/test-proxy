package store_impl

import (
	"personal/go-proxy-service/pkg/models"
	"personal/go-proxy-service/pkg/utilities"
)

type TaskRepository struct {
	store *Store
}

const (
	queryInsertTask = `insert into tasks (url, request_headers, method, status, body)
						values ($1, $2, $3, $4, $5) returning id`
	queryUpdateTaskResponse = `update tasks
								set updated_on = $1,
								status = $2,
								http_status_code = $3,
								response_headers = $4,
								length = $5
								where id = $6
								returning id`
	queryUpdateTaskRequest = `update tasks
								set updated_on = $1,
								status = $2
								where id = $3
								returning id`
	querySelectTask = `select id, status, http_status_code, response_headers, length 
						from tasks where id = $1`
)

func (tr *TaskRepository) CreateTask(task *models.Task) (*models.Task, error) {

	tx, err := tr.store.BeginTx()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(queryInsertTask, task.Url, task.RequestHeaders, task.Method, task.Status, task.Body)
	id := utilities.NewJsonNullInt64(0)
	if err = row.Scan(&id); err != nil {
		tx.Rollback()
		return nil, err
	}
	task.Id = id

	tx.Commit()

	return task, nil

}

func (tr *TaskRepository) UpdateTask(task *models.Task, typeBody string)  error {

	tx, err := tr.store.BeginTx()
	if err != nil {
		return err
	}

	id := utilities.NewJsonNullInt64(0)
	if typeBody == models.Request {
		err = tx.QueryRow(queryUpdateTaskRequest, task.UpdatedOn, task.Status, task.Id).Scan(&id)
	} else if typeBody == models.Response {
		err = tx.QueryRow(queryUpdateTaskResponse, task.UpdatedOn, task.Status, task.HttpStatusCode,
			task.ResponseHeaders, task.Length, task.Id).Scan(&id)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (tr *TaskRepository) GetTask(id int64) (*models.Task, error) {

	task := new(models.Task)
	row := tr.store.db.QueryRow(querySelectTask, id)

	err := row.Scan(&task.Id, &task.Status, &task.HttpStatusCode, &task.ResponseHeaders, &task.Length)
	if err != nil {
		return nil, err
	}

	return task, nil
}

