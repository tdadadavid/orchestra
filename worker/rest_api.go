package worker

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

// References
// https://kubernetes.io/docs/concepts/overview/kubernetes-api/
// https://kubernetes.io/docs/reference/using-api/api-concepts/

type API struct {
	Router  *chi.Mux
	Port    int
	Address string
	Worker  *Worker
}

type ErrorResponse struct {
	Message        string `json:"message"`
	HttpStatusCode int    `json:"http_status_code"`
}

func (a *API) APIError(w http.ResponseWriter, code int, errMsg string) {
	log.Println(errMsg)
	w.WriteHeader(code)
	e := ErrorResponse{
		HttpStatusCode: code,
		Message:        errMsg,
	}
	json.NewEncoder(w).Encode(e)
	return
}

func (a *API) initRouter() {
	a.Router = chi.NewRouter()

	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
}

func (a *API) Start() {
	a.initRouter()
	addr := fmt.Sprintf("%s:%d", a.Address, a.Port)
	fmt.Printf("Server running on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("Error starting http server: %s", err)
		return
	}
}
