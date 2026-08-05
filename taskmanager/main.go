package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"taskmanager/storage"
	"taskmanager/task"
	"text/tabwriter"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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

	go func() {
		<-sigs
		fmt.Println("Received SIGNAL, exiting os")
		err := manager.Flush()
		fmt.Println(err)

		os.Exit(0)
	}()

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	switch {
	case *flagAdd != "":
		err := manager.Add(*flagAdd)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	case *flagList:
		tasks := manager.List(*flagAll)

		fmt.Fprint(writer, "ID\tTitle\tDone\tDate\n")
		for _, task := range tasks {
			mark := " "
			if task.Done {
				mark = "x"
			}
			fmt.Fprintf(writer, "%d\t%s\t%s\t%s\n", task.ID, task.Title, mark, task.CreatedAt.Format("2006-01-02 15:04"))
		}
	case *flagDone != -1:
		err := manager.MarkDone(*flagDone)
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		msg := <-manager.Channel
		fmt.Println(msg)

	default:
		flag.Usage()
	}

	err = writer.Flush()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
