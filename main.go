package main

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"orchestra/manager"
	"orchestra/node"
	"orchestra/task"
	"orchestra/worker"
	"time"
)

func main() {
	t := task.Task{
		ID:            uuid.New(),
		Name:          "test",
		State:         task.Pending,
		Image:         "ubuntu",
		Memory:        "90",
		Disk:          1,
		ExposedPorts:  make(nat.PortSet),
		PortBindings:  map[string]string{},
		RestartPolicy: "on-failure",
		StartTime:     time.Now(),
		FinishTime:    time.Now(),
	}

	tEvent := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		TimeStamp: time.Now(),
		Task:      t,
	}

	fmt.Printf("Task: %v\n, event: %v\n", t, tEvent)

	w := worker.Worker{
		Name:      "Worker-1",
		Queue:     *queue.New(),
		Db:        make(map[uuid.UUID]task.Task),
		TaskCount: 0,
	}
	fmt.Printf("Worker: %v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending:  *queue.New(),
		EventsDb: make(map[string][]task.TaskEvent),
		TasksDb:  make(map[string][]task.Task),
		Workers: []string{
			w.Name,
		},
		WorkerTaskMap: make(map[string][]uuid.UUID),
		TaskWorkerMap: make(map[uuid.UUID]string),
	}

	n := node.Node{
		Name:            "Daniels",
		IpAddr:          "192.0.0.1",
		Cores:           4,
		Memory:          1024,
		MemoryAllocated: 0,
		Disk:            25,
		DiskAllocated:   0,
		Role:            "Worker",
		TaskCount:       0,
	}
	fmt.Printf("Manager: %v\n, Node: %v\n", m, n)
}
