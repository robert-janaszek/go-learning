package storage

import "taskmanager/task"

type Storage interface {
	Save(tasks []task.Task) error
	Load() ([]task.Task, error)
}
