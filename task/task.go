package task

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"log"
	"os"
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

// Config represents the configuration settings for a container.
type Config struct {
	Name          string   // Name denotes the name of the container.
	AttachStdin   bool     // AttachStdin specifies whether to attach the container's standard input.
	AttachStdout  bool     // AttachStdout specifies whether to attach the container's standard output.
	AttachStderr  bool     // AttachStderr specifies whether to attach the container's standard error.
	Cmd           []string // Cmd specifies the command to run in the container.
	Image         string   // Image denotes the container image to use.
	Memory        int64    // Memory specifies the memory limit (in bytes) for the container.
	Disk          int64    // Disk specifies the disk space limit (in bytes) for the container.
	Env           []string // Env lists the environment variables for the container.
	RestartPolicy string   // RestartPolicy defines the restart policy for the container.
}

type DockerAction string

const (

	// PULL represents the action to pull a Docker image from a repository.
	PULL DockerAction = "pull"

	// START represents the Docker action to start a container.
	START DockerAction = "start"

	// STOP represents the action to stop a Docker container.
	STOP DockerAction = "stop"

	// REMOVE represents the action of removing a Docker container or image.
	REMOVE DockerAction = "remove"

	// CREATE represents the action of creating a Docker container or image
	CREATE DockerAction = "create"
)

type DockerResultMessage string

const (
	SUCCESS DockerResultMessage = "Success"
	FAILURE                     = "Failure"
)

// DockerResult captures the outcome of a Docker operation, including error details, action type, container ID, and result message.
type DockerResult struct {
	Error       error
	Action      DockerAction
	ContainerId string
	Result      DockerResultMessage
}

type Docker struct {
	Client      *client.Client
	Config      Config
	ContainerId string
}

// Run this performs the same duty of 'docker run' on your command-line.
func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	rc, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Name, d.Config)
		return DockerResult{
			Error:  err,
			Action: PULL,
		}
	}
	// copy the reader to the standard output of the container.
	io.Copy(os.Stdout, rc)

	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}
	r := container.Resources{
		Memory: d.Config.Memory,
	}
	cc := container.Config{
		Image: d.Config.Image,
		Env:   d.Config.Env,
	}
	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container %s: %v\n", d.Config.Name, err)
		return DockerResult{
			Error:  err,
			Action: CREATE,
		}
	}

	if err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Printf("Error starting container %s: %v\n", d.Config.Name, err)
		return DockerResult{
			Error:  err,
			Action: START,
		}
	}

	// track the containerID.
	d.ContainerId = resp.ID

	out, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Printf("Error getting container logs: %v\n", err)
		return DockerResult{
			Error: err,
		}
	}

	// copy the logs of the container to the host stand output/error
	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, out); err != nil {
		log.Printf("Error copying container logs: %v\n", err)
		return DockerResult{}
	}

	return DockerResult{
		ContainerId: resp.ID,
		Action:      START,
		Result:      SUCCESS,
	}
}

// Stop performs the same function as both 'docker stop' and 'docker rm' commands.
func (d *Docker) Stop() DockerResult {
	log.Printf("Stopping container %s\n", d.Config.Name)
	ctx := context.Background()
	if err := d.Client.ContainerStop(ctx, d.ContainerId, container.StopOptions{}); err != nil {
		log.Printf("Error stopping container %s: %v\n", d.Config.Name, err)
		panic(err)
	}

	if err := d.Client.ContainerRemove(ctx, d.ContainerId, container.RemoveOptions{}); err != nil {
		log.Printf("Error removing container %s: %v\n", d.Config.Name, err)
		panic(err)
	}

	return DockerResult{
		ContainerId: d.ContainerId,
		Action:      REMOVE,
		Result:      SUCCESS,
	}
}
