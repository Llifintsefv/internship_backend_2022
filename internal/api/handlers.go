package api

import (
	"encoding/json"
	"fmt"
	"internship_backend_2022/internal/models"
	"internship_backend_2022/internal/service"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		return 
	}
	
type DepositRequestDTO struct {
    UserID int           `json:"user_id"`
    Amount json.Number `json:"amount"`
}
	var dto DepositRequestDTO
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w,"bad request",http.StatusBadRequest)
		return
	}

	amount, err := dto.Amount.Float64()
    if err != nil {
        http.Error(w, "invalid amount format", http.StatusBadRequest)
        return
    }

	DepositRequest := models.DepositRequest{
        UserID: dto.UserID,
        Amount: big.NewFloat(amount),
    }

	
	fmt.Println(DepositRequest)
	DepositResponse,err := h.service.Deposit(ctx,DepositRequest)
	if err != nil { 
		log.Print(err)
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
	

	if err := json.NewEncoder(w).Encode(DepositResponse); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return 
	}
	ctx := r.Context()
	params := mux.Vars(r)
	fmt.Println(params["user_id"])
	userID,err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w,"bad request",http.StatusBadRequest)
		return
	}

	balance, err := h.service.GetUserBalance(ctx,userID)
	if err != nil { 
		log.Print(err)
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}
	