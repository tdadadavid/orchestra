package worker

import (
	"errors"
	"fmt"
	"log"
	"orchestra/task"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

// @NOTES
// Responsibilities of a worker:
//1. Run tasks as Docker containers.
//2. Accept tasks to run from a manager.
//3. Provide relevant statistics to the manager for the purpose of scheduling
//tasks.
//4. Keep track of its tasks and their state.

//There are two possible scenarios for handling tasks:
//1. A task is being submitted for the first time, so the Worker will not know
//about it
//2. A task is being submitted for the Nth time, where the task submitted
//represents the desired state to which the current task should transition

// Worker represents a worker that processes tasks.
// we’re using the worker’s datastore (db) to
// represent the current state of tasks, while we’re using the worker’s queue to
// represent the desired state of task
type Worker struct {
	Name      string                   //
	Queue     queue.Queue              //
	Db        map[uuid.UUID]*task.Task // Db maps task identifiers (UUID) to their respective Task objects.
	TaskCount int                      //
	Stats *Stats
}

// RunTask starts or stop a task based on its current state
func (w *Worker) RunTask() task.DockerResult {
	//1. Pull a task of the queue.
	elem := w.Queue.Dequeue()
	if elem == nil {
		log.Println("No task in the queue")
		return task.DockerResult{Error: nil}
	}

	// 2. Convert it from an interface to a task.Task type.
	queuedTask, ok := elem.(task.Task) // convert the interface to of type task.
	if !ok {
		log.Println("Element is not of type task.Task{}")
		return task.DockerResult{Error: errors.New("element is not of type task.Task{}")}
	}

	// 3. Retrieve the task from the worker’s Db.
	persistedTask, found := w.Db[queuedTask.ID]
	if !found {
		log.Printf("Task %s is not in the worker's Db", queuedTask.ID)
	}

	if persistedTask == nil {
		persistedTask = &queuedTask
		w.Db[queuedTask.ID] = &queuedTask
	}

	// 4. Check if the state transition is valid.
	var result task.DockerResult
	if task.ValidateStateTransition(persistedTask.State, queuedTask.State) {
		switch queuedTask.State {
		case task.Scheduled: // 5. If the task from the queue is in a state Scheduled, call StartTask.
			result = w.StartTask(queuedTask)
		case task.Completed: // 6. If the task from the queue is in a state Completed, call StopTask.
			result = w.StopTask(queuedTask)
		default:
			result.Error = errors.New("invalid state transition")
		}
	} else {
		// 7. Else there is an invalid transition, so return an error.
		err := fmt.Errorf("invalid state transition from %v to %v", persistedTask.State, queuedTask.State)
		result.Error = err
	}

	return result
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
	w.UpdateTaskCount()
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()

	config := task.NewConfig(&t)
	d := task.NewDocker(config)

	result := d.Run()
	if result.Error != nil {
		log.Printf("Error running container: %v: %v", d.ContainerId, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}

	t.FinishTime = time.Now().UTC()
	t.Runtime.ContainerId = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = &t

	log.Printf("Started and ran container %v for task %v", d.ContainerId, t.Name)

	return result
}

// StopTask similar to 'docker stop <container_id>'
func (w *Worker) StopTask(t task.Task) task.DockerResult {
	//1. Create an instance of the Docker struct that allows us to talk to the Docker daemon using the Docker SDK.
	config := task.NewConfig(&t)
	d := task.NewDocker(config)

	// 2. Call the Stop() method on the Docker struct.
	result := d.Stop(t.Runtime.ContainerId)
	// 3. Check if there were any errors in stopping the task.
	if result.Error != nil {
		log.Printf("Error stopping container: %v: %v", d.ContainerId, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}

	// 4. Update the FinishTime field on the task t.
	t.FinishTime = time.Now().UTC()
	t.Runtime.ContainerId = result.ContainerId
	t.State = task.Completed
	// 5. Save the updated task t to the worker’s Db field.
	w.Db[t.ID] = &t

	// 6. Print an informative message and return the result of the operation
	log.Printf("Stopped and removed container %v for task %v", d.ContainerId, t.Name)

	return result
}

// GetTasks this fetches all the tasks in the workers store.
func (w *Worker) GetTasks() []*task.Task {
	tasks := make([]*task.Task, 0)
	for _, t := range w.Db {
		tasks = append(tasks, t)
	}
	return tasks
}


func (w *Worker) CollectStats() {
	for {
		log.Println("Collecting stats")
		w.Stats = GetStats()
		time.Sleep(15 * time.Second) // collect stats every 15 seconds.
	}
}

func (w *Worker) UpdateTaskCount() {
	w.TaskCount = w.Queue.Len()
	w.Stats.TaskCount = w.TaskCount
}
