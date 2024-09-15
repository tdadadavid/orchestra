package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"log"
	"orchestra/task"
	"os"
	"time"
)

func main() {
	fmt.Printf("Create test container \n")
	dTask, createResult := createContainer("ubuntu")
	if createResult.Error != nil {
		log.Printf("Error creating container: %v", createResult.Error)
		os.Exit(1)
	}

	time.Sleep(10 * time.Second)

	result := stopContainer(dTask)
	if result.Error != nil {
		log.Printf("Error stopping container: %v", result.Error)
	}

	fmt.Printf("Done testing \n")
}

func createContainer(image string) (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test_container_1",
		Image: image,
		Env:   []string{},
		Cmd: []string{
			"new Date()",
		},
	}
	dc, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	// run the docker image
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s is running with config ", result.ContainerId)
	return &d, &result
}

func stopContainer(d *task.Docker) *task.DockerResult {
	result := d.Stop()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}

	fmt.Printf("Container %s is stopped and removed.", result.ContainerId)
	return &result
}
