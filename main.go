package main

import (
	"flag"
	"log"
	"orchestra/task"
	"orchestra/worker"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

var (
	Host string  
	Port int 
)

func setupFlags() {
	flag.StringVar(&Host, "host", "localhost", "Host on which orchestra runs")
	flag.IntVar(&Port, "port", 7777, "port to run orchestra")
}


//TODO: Check the goprocinfo libarary to update `stats.go`
// ioutil.ReadFile(path) code.
func main() {
	setupFlags()

	// host := os.Getenv("ORCHESTRA_HOST")
	// port, _ := strconv.Atoi(os.Getenv("ORCHESTRA_PORT"))

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	api := worker.API{Address: Host, Port: Port, Worker: &w}
	go runTasks(&w)
	go w.CollectStats()
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
