package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	RouteHealthV1            = "/v1/health"
	RouteAuthRegisterV1      = "/v1/auth/register"
	RouteAuthLoginV1         = "/v1/auth/login"
	RouteAuthTokenV1         = "/v1/auth/token"
	RouteTokensV1            = "/v1/tokens"
	RouteTokenDeleteV1       = "/v1/tokens/{tokenId}"
	RouteProjectsV1          = "/v1/projects"
	RouteProjectDeleteV1     = "/v1/projects/{projectId}"
	RouteProjectDeployV1     = "/v1/projects/{projectId}/deploy"
	RouteProjectDeploymentV1 = "/v1/projects/{projectId}/deployments/{deploymentId}"
	RouteProjectEnvV1        = "/v1/projects/{projectId}/env"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(RouteHealthV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	r.HandleFunc(RouteAuthRegisterV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteAuthRegisterV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteAuthLoginV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteAuthLoginV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteAuthTokenV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteAuthTokenV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteTokensV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET " + RouteTokensV1))
	}).Methods(http.MethodGet)

	r.HandleFunc(RouteTokensV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST " + RouteTokensV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteTokenDeleteV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)

	r.HandleFunc(RouteProjectsV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteProjectsV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteProjectDeleteV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)

	r.HandleFunc(RouteProjectDeployV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteProjectDeployV1))
	}).Methods(http.MethodPost)

	r.HandleFunc(RouteProjectDeploymentV1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(RouteProjectDeploymentV1))
	}).Methods(http.MethodGet)

	r.HandleFunc(RouteProjectEnvV1, func(w http.ResponseWriter, r *http.Request) {
		// Note: The API spec says 204 OK for this one, using StatusNoContent.
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodPost)

	return r
}
