package api

import "github.com/gorilla/mux"

func SetupRouter(handler *handler) *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/deposit", handler.Deposit).Methods("POST")

	return router
}