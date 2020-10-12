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
	hostname := fmt.Sprintf(":%d", controller.cfg.HTTPPort())
	return http.ListenAndServe(hostname, controller.router)
}

func (controller *Controller) ping(w http.ResponseWriter, req *http.Request) {
	controller.cfg.Logger().Info().Str("uri", req.RequestURI).Str("remote", req.RemoteAddr).Msg("Request received")
	fmt.Fprintf(w, "pong")
}
