package repository

import "personal/go-proxy-service/pkg/models"

type TaskRepoInt interface {
	CreateTask(task *models.Task) (*models.Task, error)
	UpdateTask(task *models.Task, typeBody string) error
	GetTask(id int64) (*models.Task, error)
}
