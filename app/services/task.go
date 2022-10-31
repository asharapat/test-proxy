package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"personal/go-proxy-service/app/repository"
	"personal/go-proxy-service/pkg/models"
	"personal/go-proxy-service/pkg/utilities"
)

type TaskService struct {
	store repository.StoreInt
}

func NewService(store repository.StoreInt) *TaskService {
	return &TaskService{store: store}
}

func (svc *TaskService) CreateTask(task *models.Task) (*models.Task, error) {

	task.Status = utilities.NewJsonNullString(models.StatusNew)
	task, err := svc.store.TaskRepository().CreateTask(task)
	if err != nil {
		return nil, err
	}

	taskRes := new(models.Task)
	taskRes.Id = task.Id

	// start sending request in background
	go svc.DoRequest(task)

	return taskRes, nil
}

func (svc *TaskService) UpdateTask(task *models.Task, bodyType string) error {
	task.UpdatedOn = utilities.NewJsonNullTime(time.Now())
	return svc.store.TaskRepository().UpdateTask(task, bodyType)
}

func (svc *TaskService) GetTask(id int64) (*models.Task, error) {
	return svc.store.TaskRepository().GetTask(id)
}

func (svc *TaskService) DoRequest(task *models.Task) {

	task.Status = utilities.NewJsonNullString(models.StatusInProcess)
	err := svc.UpdateTask(task, models.Request)
	if err != nil {
		log.Println(err)
		return
	}

	var reqBody *bytes.Buffer

	if task.Body != nil {
		jsonBody, err := json.Marshal(&task.Body)
		if err != nil {
			log.Println(err)
			task.Status = utilities.NewJsonNullString(models.StatusError)
			if err = svc.UpdateTask(task, models.Request); err != nil {
				log.Println(err)
			}
			return
		}

		reqBody = bytes.NewBuffer(jsonBody)
	}

	var request *http.Request

	if reqBody == nil {
		request, err = http.NewRequest(task.Method.String, task.Url.String, nil)
	} else {
		request, err = http.NewRequest(task.Method.String, task.Url.String, reqBody)
	}
	if err != nil {
		log.Println(err)
		task.Status = utilities.NewJsonNullString(models.StatusError)
		if err = svc.UpdateTask(task, models.Request); err != nil {
			log.Println(err)
		}
		return
	}
	params := request.URL.Query()
	request.URL.RawQuery = params.Encode()

	if task.RequestHeaders != nil {
		if task.RequestHeaders.Authorization != nil {
			request.Header.Add("authorization", task.RequestHeaders.Authorization.String)
		}
		if task.RequestHeaders.ContentType != nil {
			request.Header.Add("content-type", task.RequestHeaders.ContentType.String)
		}
	}

	request.Close = true

	client := http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		log.Println(err)
		task.Status = utilities.NewJsonNullString(models.StatusError)
		if err = svc.UpdateTask(task, models.Request); err != nil {
			log.Println(err)
		}
		return
	}

	defer resp.Body.Close()

	// want to check response status code for error or success
	switch resp.StatusCode / 100 {
	case 2:
		task.Status = utilities.NewJsonNullString(models.StatusDone)
	case 4, 5:
		task.Status = utilities.NewJsonNullString(models.StatusError)
	}

	task.HttpStatusCode = utilities.NewJsonNullInt32(int32(resp.StatusCode))
	task.ResponseHeaders = models.ResHeaders(resp.Header)
	task.Length = utilities.NewJsonNullInt64(resp.ContentLength)

	if err = svc.UpdateTask(task, models.Response); err != nil {
		log.Println(err)
	} else {
		log.Println("Request made successfully")
	}

}
