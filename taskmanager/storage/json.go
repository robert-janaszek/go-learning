// Stwórz strukturę JSONStorage posiadającą pole filename string.
// Zaimplementuj metody Save oraz Load przy użyciu standardowych
// pakietów os oraz encoding/json (json.MarshalIndent / json.Unmarshal).

package storage

import (
	"encoding/json"
	"os"
	"taskmanager/task"
)

type JSONStorage struct {
	filename string
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

	err = os.WriteFile(js.filename, data, 0o644)

	if err != nil {
		return err
	}

	return nil
}

func (js *JSONStorage) Load() ([]task.Task, error) {
	val, err := os.ReadFile(js.filename)

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
