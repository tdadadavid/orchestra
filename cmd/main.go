package cmd

import (
	"flag"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"orchestra/task"
	"orchestra/worker"
	"time"
)

var (
	Host string
	Port int
)

func setupFlags() {
	flag.StringVar(&Host, "host", "localhost", "Host on which orchestra runs")
	flag.IntVar(&Port, "port", 7777, "port to run orchestra")
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

// Execute TODO: Check the goprocinfo library to update `stats.go` ioutil.ReadFile(path) code.
func Execute() {
	setupFlags()

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	api := worker.API{Address: Host, Port: Port, Worker: &w}
	go runTasks(&w)
	go w.CollectStats()
	api.Start()
}
