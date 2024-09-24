package main

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"orchestra/task"
	"orchestra/worker"
	"os"
	"strconv"
	"time"
)

func main() {
	host := os.Getenv("ORCHESTRA_HOST")
	port, _ := strconv.Atoi(os.Getenv("ORCHESTRA_PORT"))

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "local_dice",
		State: task.Scheduled,
		Image: "dicedb/dicedb",
	}

	w.AddTask(t)

	api := worker.API{Address: host, Port: port, Worker: &w}
	go runTasks(&w)
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %s", result.Error)
			}
		} else {
			log.Printf("No task found to be processed in the queue.")
		}
		log.Println("sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}
