package models

import "math/big"

type DepositRequest struct {
    UserID int      `json:"user_id"`
    Amount *big.Float `json:"amount"`
}

type DepositResponse struct {
    Status        string     `json:"status"`
    Message       string     `json:"message"`
    Balance       *big.Float `json:"balance"`
    TransactionID int        `json:"transaction_id"`
}

type ReserveRequest struct {
    UserID    int      `json:"user_id"`
    ServiceID int      `json:"service_id"`
    OrderID   int      `json:"order_id"`
    Amount    *big.Float `json:"amount"`
}
