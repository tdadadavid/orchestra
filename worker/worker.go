package worker

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"orchestra/task"
)

// Responsibilities of a worker:
//1. Run tasks as Docker containers.
//2. Accept tasks to run from a manager.
//3. Provide relevant statistics to the manager for the purpose of scheduling
//tasks.
//4. Keep track of its tasks and their state.

// Worker represents a worker that processes tasks.
type Worker struct {
	Name      string                  //
	Queue     queue.Queue             //
	Db        map[uuid.UUID]task.Task // Db maps task identifiers (UUID) to their respective Task objects.
	TaskCount int                     //
}

func (w Worker) RunTask() {
	fmt.Println("Task Running")
}

func (w Worker) StartTask() {
	fmt.Println("Task Started")
}

func (w Worker) StopTask() {
	fmt.Println("Task Stopped")
}

func (w Worker) CollectStats() {
	fmt.Println("Collecting Stats")
}
