package main

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"orchestra/task"
	"orchestra/worker"
	"time"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "local_dice",
		State: task.Scheduled,
		Image: "dicedb/dicedb",
	}

	// first time the worker will see the task.
	fmt.Println("starting task.")
	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.Runtime.ContainerId = result.ContainerId
	fmt.Printf("task is runnig in container %s\n", t.Runtime.ContainerId)
	fmt.Println("sleepy time")
	time.Sleep(time.Second * 30)

	fmt.Printf("stopping task %s\n", t.Runtime.ContainerId)
	t.State = task.Completed
	w.AddTask(t)
	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
}
