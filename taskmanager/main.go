package main

import (
	"flag"
	"fmt"
	"log"
	"taskmanager/storage"
	"taskmanager/task"
)

func main() {

	flagAdd := flag.String("add", "", "-add \"Buy milk\"")
	flagList := flag.Bool("list", false, "-list, -list -all")
	flagAll := flag.Bool("all", false, "-list -all")
	flagDone := flag.Int("done", -1, "-done 1")

	flag.Parse()

	js := storage.NewJSONStorage("test.json")
	manager, err := task.NewManager(js)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	switch {
	case *flagAdd != "":
		err := manager.Add(*flagAdd)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	case *flagList:
		fmt.Println(manager.List(*flagAll))
	case *flagDone != -1:
		err := manager.MarkDone(*flagDone)
		if err != nil {
			log.Fatalf("%v\n", err)
		}

	default:
		flag.Usage()
	}
}
