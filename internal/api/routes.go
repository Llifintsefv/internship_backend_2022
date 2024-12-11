package api

import "github.com/gorilla/mux"

func SetupRouter(handler *handler) *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/deposit", handler.Deposit).Methods("POST")
	router.HandleFunc("/balance/{user_id}", handler.GetUserBalance).Methods("GET")
	router.HandleFunc("/reserve",handler.Reserve).Methods("POST")
	router.HandleFunc("/confirm",handler.Confirm).Methods("POST")
	router.HandleFunc("/transfer",handler.Transfer).Methods("POST")
	

	return router
}