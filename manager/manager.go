package manager

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"orchestra/task"
)

// Responsibilities of a Manager.
//1. Accept requests from users to start and stop tasks.
//2. Scheduled tasks onto worker machines.
//3. Keep track of tasks, their states, and the machine on which they run.

// Manager is responsible for managing tasks and workers within the system.
type Manager struct {
	Pending       queue.Queue                 // Pending is a queue that holds tasks waiting to be processed.
	EventsDb      map[string][]task.TaskEvent // EventsDb maps task names to slices of TaskEvent, representing the history of events associated with each task.
	TasksDb       map[string][]task.Task      // TasksDb maps task names to lists of Task objects, categorizing tasks within the system.
	Workers       []string                    // Workers is a list of worker identifiers (name) assigned to manage tasks within the system.
	WorkerTaskMap map[string][]uuid.UUID      // WorkerTaskMap maps worker identifiers to lists of UUIDs representing the tasks they are responsible for.
	TaskWorkerMap map[uuid.UUID]string        // TaskWorkerMap maps task UUIDs to worker identifiers, indicating which worker is responsible for each task.
}

// SelectWorker is responsible for checking the needs of the tasks and check which worker should(is capable) of handling this.
func (m *Manager) SelectWorker() {
	fmt.Println("Selecting worker")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("Updating tasks")
}

func (m *Manager) SendWork() {
	fmt.Println("Sending work to workers")
}
