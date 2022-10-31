package store_impl

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"personal/go-proxy-service/pkg/models"
	"personal/go-proxy-service/pkg/utilities"
)

func TestTaskRepository_CreateTask(t *testing.T) {
	db, mock := MockDb()

	defer db.Close()

	st := New(db)

	mockRows := mock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()

	task := &models.Task{
		Url: utilities.NewJsonNullString("url"),
		RequestHeaders: &models.Header{
			Authorization: utilities.NewJsonNullString("auth"),
			ContentType:   utilities.NewJsonNullString("content"),
		},
		Method: utilities.NewJsonNullString("method"),
		Body:   make(map[string]interface{}),
	}

	mock.ExpectQuery("insert").WithArgs(task.Url, task.RequestHeaders, task.Method, task.Status, task.Body).
		WillReturnRows(mockRows)
	mock.ExpectCommit()

	restTask, err := st.TaskRepository().CreateTask(task)
	assert.NotNil(t, restTask)
	assert.NoError(t, err)

}

func TestTaskRepository_UpdateTask(t *testing.T) {
	db, mock := MockDb()

	defer db.Close()

	st := New(db)

	task := &models.Task{
		Id:        utilities.NewJsonNullInt64(1),
		UpdatedOn: utilities.NewJsonNullTime(time.Now()),
		Url:       utilities.NewJsonNullString("url"),
		RequestHeaders: &models.Header{
			Authorization: utilities.NewJsonNullString("auth"),
			ContentType:   utilities.NewJsonNullString("content"),
		},
		Method:          utilities.NewJsonNullString("method"),
		Body:            make(map[string]interface{}),
		Status:          utilities.NewJsonNullString("status"),
		HttpStatusCode:  utilities.NewJsonNullInt32(200),
		ResponseHeaders: make(map[string][]string),
		Length:          utilities.NewJsonNullInt64(5),
	}

	mockRows1 := mock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery("update").WithArgs(task.UpdatedOn, task.Status, task.Id).WillReturnRows(mockRows1)
	mock.ExpectCommit()
	err := st.TaskRepository().UpdateTask(task, models.Request)
	assert.NoError(t, err)

	mockRows2 := mock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery("update").WithArgs(task.UpdatedOn, task.Status, task.HttpStatusCode,
		task.ResponseHeaders, task.Length, task.Id).WillReturnRows(mockRows2)
	mock.ExpectCommit()
	err = st.TaskRepository().UpdateTask(task, models.Response)
	assert.NoError(t, err)

}

func TestTaskRepository_GetTask(t *testing.T) {
	db, mock := MockDb()
	defer db.Close()
	st := New(db)

	mockRows := mock.NewRows([]string{"id", "status", "http_status_code", "response_headers", "length"}).
		AddRow(1, "ok", 200, []byte(`{"key":["value1", "value2"]}`), 1)
	mock.ExpectQuery("select").WillReturnRows(mockRows)

	res, err := st.TaskRepository().GetTask(1)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func MockDb() (*sql.DB, sqlmock.Sqlmock) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		log.Println(nil, "error while creating mock db: %s", err)
	}
	return sqlDb, mock
}
