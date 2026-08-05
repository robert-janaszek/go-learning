package task

import "slices"

type TaskManager struct {
	tasks   []Task
	storage Storage
}

func NewManager(s Storage) (*TaskManager, error) {
	ts := TaskManager{
		storage: s,
	}

	tasks, err := ts.storage.Load()
	if err != nil {
		return nil, err
	}

	ts.tasks = tasks

	return &ts, nil
}

func (ts *TaskManager) Add(title string) error {
	task, err := NewTask(len(ts.tasks)+1, title)

	if err != nil {
		return err
	}

	ts.tasks = append(ts.tasks, *task)

	err = ts.storage.Save(ts.tasks)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskManager) MarkDone(id int) error {
	found := false
	for i := range ts.tasks {
		task := &ts.tasks[i]

		if task.ID == id {
			task.Done = true
			found = true
			break
		}
	}

	if !found {
		return ErrTaskNotFound
	}

	err := ts.storage.Save(ts.tasks)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskManager) List(showAll bool) []Task {
	if showAll {
		return slices.Clone(ts.tasks)
	}

	tasks := []Task{}

	for _, task := range ts.tasks {
		if !task.Done {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (ts *TaskManager) Flush() error {
	return ts.storage.Save(ts.tasks)
}
