package api

import (
	"encoding/json"
	"internship_backend_2022/internal/models"
	"internship_backend_2022/internal/service"
	"net/http"
)

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) *handler {
	return &handler{
		service: service,
	}
}


func (h *handler) Deposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		
	}

	ctx := r.Context()
	var DepositResponse models.DepositResponse
	err := json.NewDecoder(r.Body).Decode(&DepositResponse)
	if err != nil {

	}

	DepositRequest,err := h.service.Deposit(ctx,DepositResponse)

	if err := json.NewEncoder(w).Encode(DepositRequest); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}