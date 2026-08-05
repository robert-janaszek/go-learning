package task

import (
	"testing"
)

type MockStorage struct {
	tasks []Task
}

func (m *MockStorage) Save(tasks []Task) error {
	m.tasks = tasks

	return nil
}

func (m *MockStorage) Load() ([]Task, error) {
	return m.tasks, nil
}

func TestTaskManager(t *testing.T) {
	m := MockStorage{}

	tm, err := NewManager(&m)

	if err != nil {
		t.Fatal("cannot create new manager")
	}

	err = tm.Add("test")

	if err != nil {
		t.Fatalf("failed add new task, found %v", err)
	}

	if len(m.tasks) != 1 {
		t.Fatalf("expected 1 task, found %d", len(m.tasks))
	}

	if m.tasks[0].Title != "test" {
		t.Fatalf("expected tasks title to be \"test\", found %s", m.tasks[0].Title)
	}
}
