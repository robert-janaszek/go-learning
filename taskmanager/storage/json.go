package storage

import (
	"encoding/json"
	"os"
	"sync"
	"taskmanager/task"
)

type JSONStorage struct {
	filename string
	mutex    sync.Mutex
}

func NewJSONStorage(filename string) *JSONStorage {
	return &JSONStorage{
		filename: filename,
	}
}

func (js *JSONStorage) Save(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	js.mutex.Lock()
	defer js.mutex.Unlock()

	err = os.WriteFile(js.filename, data, 0o644)

	if err != nil {
		return err
	}

	return nil
}

func (js *JSONStorage) Load() ([]task.Task, error) {
	js.mutex.Lock()
	val, err := os.ReadFile(js.filename)
	js.mutex.Unlock()

	if err != nil {
		return nil, err
	}

	var tasks []task.Task

	err = json.Unmarshal(val, &tasks)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
