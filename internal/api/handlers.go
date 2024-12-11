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
		http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
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
		http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
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

	BalanceResponse, err := h.service.GetUserBalance(ctx,userID)
	if err != nil { 
		log.Print(err)
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(BalanceResponse); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}


func (h *handler)Reserve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		return 
	}

	type ReserveRequestDTO struct {
    UserID    int      `json:"user_id"`
    ServiceID int      `json:"service_id"`
    OrderID   int      `json:"order_id"`
    Amount    json.Number  `json:"amount"`
	}
	var dto ReserveRequestDTO
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

	ReserveRequest := models.ReserveRequest{
		UserID:    dto.UserID,
		ServiceID: dto.ServiceID,
		OrderID:   dto.OrderID,
		Amount:    big.NewFloat(amount),
	}
	ReserveResponse,err := h.service.Reserve(ctx,ReserveRequest)
	if err != nil { 
		log.Print(err)
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(ReserveResponse); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}
	
func (h *handler)Confirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		return 
	}

	type ConfirmRequestDTO struct {
	UserID    int      `json:"user_id"`
	ServiceID int      `json:"service_id"`
	OrderID   int      `json:"order_id"`
	Amount    json.Number  `json:"amount"`
	}
	var dto ConfirmRequestDTO
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

	ConfirmRequest := models.ConfirmRequest{
		UserID:    dto.UserID,
		ServiceID: dto.ServiceID,
		OrderID:   dto.OrderID,
		Amount:    big.NewFloat(amount),
	}
	ConfirmResponse,err := h.service.Confirm(ctx,ConfirmRequest)
	if err != nil { 
		log.Print(err)	
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(ConfirmResponse); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}


func (h *handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		return 
	}

	type TransferRequestDTO struct {
	FromUserID    int      `json:"from_user_id"`
	ToUserID   int      `json:"to_user_id"`
	Amount    json.Number  `json:"amount"`
	}
	var dto TransferRequestDTO
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

	TransferRequest := models.TransferRequest{
		FromUserID:    dto.FromUserID,
		ToUserID:   dto.ToUserID,
		Amount:    big.NewFloat(amount),
	}
	TransferResponse,err := h.service.Transfer(ctx,TransferRequest)
	if err != nil { 
		log.Print(err)	
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(TransferResponse); err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}