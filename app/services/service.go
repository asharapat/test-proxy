package services

import "personal/go-proxy-service/pkg/models"

type TaskServiceInt interface {
	CreateTask(task *models.Task) (*models.Task, error)
	UpdateTask(task *models.Task, bodyType string) error
	DoRequest(task *models.Task)
	GetTask(id int64) (*models.Task, error)
}
