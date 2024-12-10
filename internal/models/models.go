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

type TransactionType string

const (
    Deposit           TransactionType = "deposit"
    Withdrawal        TransactionType = "withdrawal"
    Reserve           TransactionType = "reserve"
    RevenueRecognition TransactionType = "revenue_recognition"
    TransferFrom       TransactionType = "transfer_from"
    TransferTo         TransactionType = "transfer_to"
)