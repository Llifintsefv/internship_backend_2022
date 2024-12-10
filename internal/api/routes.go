package api

import "github.com/gorilla/mux"

func SetupRouter(handler *handler) *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/deposit", handler.Deposit).Methods("POST")
	router.HandleFunc("/balance/{user_id}", handler.GetUserBalance).Methods("GET")

	return router
}