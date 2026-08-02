package main

import (
	"fmt"
	"taskmanager/storage"
	"taskmanager/task"
)

func main() {
	js := storage.NewJsonStorage("test.json")

	tasks, err := js.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tasks)

	learnGo, err := task.NewTask(3, "Learn Go")

	if err != nil {
		fmt.Println(err)
		return
	}

	moreTasks := append(tasks, *learnGo)

	err = js.Save(moreTasks)

	if err != nil {
		fmt.Println(err)
		return
	}
}
