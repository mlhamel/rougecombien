package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mlhamel/rougecombien/pkg/config"
)

// Controller is handling everything needed for handling http requests
type Controller struct {
	cfg    *config.Config
	router *mux.Router
}

// NewController is initiating a default controller
func NewController(cfg *config.Config) *Controller {
	router := mux.NewRouter().StrictSlash(true)

	instance := Controller{cfg, router}

	router.HandleFunc("/ping", instance.ping).Methods(http.MethodGet)

	return &instance
}

// Run is responsible for running the web server
func (controller *Controller) Run(ctx context.Context) error {
	return http.ListenAndServe(":8000", controller.router)
}

func (controller *Controller) ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong")
}
