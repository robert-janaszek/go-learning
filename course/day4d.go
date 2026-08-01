package course

import (
	"fmt"

	"github.com/robert-janaszek/go-learning/course/store"
)

// ex 16
type reader interface {
	Read() string
}

type writer interface {
	Write(data string)
}

type readWriter interface {
	reader
	writer
}

// ex 17
type fileHandler struct {
	content string
}

func (f *fileHandler) Read() string {
	return f.content
}

func (f *fileHandler) Write(data string) {
	f.content = data
}

type userGetter interface {
	GetUsers() []string
}

type userService struct {
	// store userGetter
	userGetter
}

type mockStore struct{}

func (m mockStore) GetUsers() []string {
	return []string{"a", "b"}
}

type customError struct {
	Code    int
	Message string
}

func (e customError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func throws() error {
	return customError{
		Code:    123,
		Message: "custom message",
	}
}

func Day4d() {
	var file readWriter = &fileHandler{}
	file.Write("test")
	fmt.Println(file.Read())

	// ex 18
	svc := userService{
		store.PostgresStore{},
	}
	fmt.Println(svc.GetUsers())

	// ex 19
	svc2 := userService{
		mockStore{},
	}

	fmt.Println(svc2.GetUsers())

	// ex 20
	err := throws()
	fmt.Println(err)
}
