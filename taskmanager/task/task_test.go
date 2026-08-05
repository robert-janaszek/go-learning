package task

import (
	"testing"
)

func TestNewTask(t *testing.T) {

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid title", "Buy milk", false},
		{"empty title", "", true},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task, err := NewTask(i, test.title)

			if test.wantErr {
				if err == nil {
					t.Errorf("%s/%s wanted error, found nil", test.name, test.title)
				}
			} else {
				if err != nil {
					t.Errorf("%s/%s unexpected error, found %v", test.name, test.title, err)
					return
				}

				if task.CreatedAt.IsZero() {
					t.Errorf("%s/%s expected createdAt to be non-zero", test.name, test.title)
				}
			}
		})
	}
}
