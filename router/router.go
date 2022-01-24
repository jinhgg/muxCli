package router

import (
	"github.com/gorilla/mux"
	"muxCli/api"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", api.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", api.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", api.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", api.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users/{id}", api.UpdateUserHandler).Methods("PUT")

	return r
}
