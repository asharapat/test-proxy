package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"personal/go-proxy-service/app/services"
	"personal/go-proxy-service/pkg/models"
	"personal/go-proxy-service/pkg/utilities"
	"testing"
)

type taskServiceTest struct {
	ans     *models.Task
	err     error
}

func (t taskServiceTest) CreateTask(task *models.Task) (*models.Task, error) {
	if t.err != nil {
		fmt.Println("ERROR", t.err)
		return nil, t.err
	}
	return t.ans, nil
}

func (t taskServiceTest) UpdateTask(task *models.Task, bodyType string) error {
	return nil
}

func (t taskServiceTest) DoRequest(task *models.Task) {
	return
}

func (t taskServiceTest) GetTask(id int64) (*models.Task, error) {
	if t.err != nil {
		return nil, t.err
	}
	return t.ans, nil
}

var service services.TaskServiceInt = &taskServiceTest{}

func TestServer_MakeCreateTaskRequest(t *testing.T) {
	r := mux.NewRouter()

	taskBody := &models.Task{
		Url: utilities.NewJsonNullString("url"),
		RequestHeaders: &models.Header{
			Authorization: utilities.NewJsonNullString("auth"),
			ContentType:   utilities.NewJsonNullString("content"),
		},
		Method: utilities.NewJsonNullString("method"),
		Body:   make(map[string]interface{}),
	}

	jsonBody, err := json.Marshal(taskBody)
	if err != nil {
		t.Error(err)
	}

	// we must have separate reqBodies for each request
	reqBody1 := bytes.NewBuffer(jsonBody)
	reqBody2 := bytes.NewBuffer(jsonBody)

	testServers := []Server{
		{
			Router: r,
			TaskService: taskServiceTest{
				ans: &models.Task{},
				err: nil,
			},
		},
		{
			Router: r,
			TaskService: taskServiceTest{
				ans: nil,
				err: errors.New("some internal error"),
			},
		},
		{
			Router: r,
			TaskService: taskServiceTest{
				ans: &models.Task{},
				err: nil,
			},
		},
	}

	for i, ts := range testServers {
		urlStr := "/api/task/create"
		ts.Router.HandleFunc(urlStr, ts.MakeCreateTaskRequest())

		var resp *http.Response
		var err error

		switch i {
		case 0:
			resp, err = MakeRequest(ts, urlStr, nil)
		case 1:
			resp, err = MakeRequest(ts, urlStr, reqBody1)
		case 2:
			resp, err = MakeRequest(ts, urlStr, reqBody2)
		}

		if err != nil {
			t.Error(err)
		}
		fmt.Println(resp.StatusCode)
	}
}

func TestServer_MakeGetTaskRequest(t *testing.T) {
	r := mux.NewRouter()
	testServers := []Server{
		{
			Router: r,
			TaskService: taskServiceTest{
				ans:     nil,
				err:     errors.New("some internal error"),
			},
		},
		{
			Router: r,
			TaskService: taskServiceTest{
				ans:     &models.Task{},
				err:     nil,
			},
		},
	}

	for _, ts := range testServers {
		urlStr := "/api/task/get/"
		ts.Router.HandleFunc(urlStr+"{id:[0-9]+}", ts.MakeGetTaskRequest())

		resp,  err := MakeRequest(ts, urlStr+"1", nil)

		if err != nil {
			t.Error(err)
		}

		fmt.Println(resp.StatusCode)

	}
}

func MakeRequest(ts Server, urlStr string, reqBody *bytes.Buffer) (*http.Response, error) {
	htts := httptest.NewServer(ts.Router)
	defer htts.Close()

	var request *http.Request
	var err error

	if reqBody == nil {
		request, err = http.NewRequest(http.MethodPost, htts.URL+urlStr, nil)
	} else {
		request, err = http.NewRequest(http.MethodPost, htts.URL+urlStr, reqBody)
	}

	if err != nil {
		return nil, err
	}

	request.Close = true

	client := http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
