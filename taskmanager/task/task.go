package task

import (
	"errors"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewTask(id int, title string) (*Task, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	task := Task{
		ID:        id,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}

	return &task, nil
}
