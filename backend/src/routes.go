package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	RouteRoot     = "/"
	RouteHealth   = "/health"
	RouteCreateV1 = "/v1/app/create"
	RouteDeployV1 = "/v1/app/deploy"
	RouteDeleteV1 = "/v1/app/delete"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(RouteRoot, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	}).Methods(http.MethodGet)

	r.HandleFunc(RouteHealth, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	r.HandleFunc(RouteCreateV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteCreateV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteDeployV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteDeployV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteDeleteV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteDeleteV1))
	}).Methods(http.MethodDelete)

	return r
}
