package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (

	// Pending indicates that a task is waiting to be scheduled or processed.
	Pending = iota

	// Schedule indicates that a task is scheduled to be run but has not started yet.
	Schedule

	// Running represents the state in which a process or task is currently in execution.
	Running

	// Completed indicates that the associated task or process has been successfully finished.
	Completed

	// Failed indicates that a process or operation has been unsuccessful in completing its intended tasks.
	Failed
)

// Task struct represents the metadata and properties associated with a specific task.
type Task struct {
	ID            uuid.UUID         // ID represents the unique identifier for a Task.
	Name          string            // Name is the human-readable identifier for
	State         State             // State represents the current status of a Task within the system.
	Image         string            // Image specifies the Docker image to be used for the task's container.
	Memory        string            // Memory is the amount of memory allocated to the task's container.
	Disk          int               // Disk is the amount of disk space allocated to the task's container in gigabytes.
	ExposedPorts  nat.PortSet       // ExposedPorts is a set of ports that are exposed by the task's container.
	PortBindings  map[string]string // PortBindings maps container ports to host ports for network binding in the task's container.
	RestartPolicy string            // RestartPolicy specifies the restart policy for the task's container, e.g., "always", "on-failure", or "never".
	StartTime     time.Time         // StartTime is the timestamp indicating when the task started.
	FinishTime    time.Time         // FinishTime is the timestamp indicating when the task finished.
}

// TaskEvent represents an event that occurs within the lifecycle of a task.
type TaskEvent struct {
	ID        uuid.UUID // ID represents the unique identifier for a TaskEvent.
	State     State     // State represents the current status of the task in the TaskEvent struct.
	TimeStamp time.Time // TimeStamp is the time at which the TaskEvent occurred.
	Task      Task      // Task represents the metadata and properties associated with a specific task.
}
