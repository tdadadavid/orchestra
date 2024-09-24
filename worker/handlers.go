package worker

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	uuid2 "github.com/google/uuid"
	"log"
	"net/http"
	"orchestra/task"
)

func (a *API) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var taskEvent task.TaskEvent
	err := d.Decode(&taskEvent)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to decode task event: %v", err)
		a.APIError(w, http.StatusBadRequest, errMsg)
		return
	}

	a.Worker.AddTask(taskEvent.Task)
	log.Printf("Added task %v\n", taskEvent.Task.ID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskEvent.Task)
}

func (a *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	wTasks := a.Worker.GetTasks()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wTasks)
}

func (a *API) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	tID := chi.URLParam(r, "taskID")
	if tID == "" {
		errMsg := fmt.Sprintf("Invalid task ID passed: %v", tID)
		a.APIError(w, http.StatusBadRequest, errMsg)
		return
	}

	uuid, err := uuid2.Parse(tID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse the uuid: %v", err)
		a.APIError(w, http.StatusBadRequest, errMsg)
		return
	}

	t, ok := a.Worker.Db[uuid]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		errMsg := fmt.Sprintf("Task not found: %v", err)
		a.APIError(w, http.StatusBadRequest, errMsg)
		return
	}

	copiedTask := *t
	copiedTask.State = task.Completed
	a.Worker.AddTask(copiedTask)

	log.Printf("Added task %v\n to be stopped.", copiedTask)
	w.WriteHeader(http.StatusNoContent)
}

//This is my own implementation of 'StartTaskHandler' based on my 'now' knowledge of Go
// but comparing my implementation to the writer's own, His is better and this is what my
// co-tutor said.
// Check it out: https://chatgpt.com/share/66f2352c-9300-8005-8f33-ee758c0eac97
//func (a *API) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
//	payload, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	defer r.Body.Close()
//
//	err = json.Unmarshal(payload, &taskEvent)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	a.Worker.AddTask(taskEvent.Task)
//	w.WriteHeader(http.StatusCreated)
//	json.NewEncoder(w).Encode(taskEvent.Task)
//}
