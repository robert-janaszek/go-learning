package main

import (
	"fmt"
	"taskmanager/storage"
	"taskmanager/task"
)

func main() {
	js := storage.NewJSONStorage("test.json")

	manager, err := task.NewManager(js)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = manager.MarkDone(1)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(manager.List(false))
}
