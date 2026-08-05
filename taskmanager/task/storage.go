package task

type Storage interface {
	Save(tasks []Task) error
	Load() ([]Task, error)
}
