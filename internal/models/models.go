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


type ReserveResponse struct {
    Status        string     `json:"status"`
    Message       string     `json:"message"`
    Balance       *big.Float `json:"balance"`
    Reserved      *big.Float `json:"reserved"`
    TransactionID int        `json:"transaction_id"`
}

type ConfirmRequest struct {
    UserID    int      `json:"user_id"`
    ServiceID int      `json:"service_id"`
    OrderID   int      `json:"order_id"`
    Amount    *big.Float `json:"amount"`
}


type ConfirmResponse struct {
    Status        string     `json:"status"`
    Message       string     `json:"message"`
    TransactionID int        `json:"transaction_id"`
}

type BalanceResponse struct {
    Balance *big.Float `json:"balance"`
    Reserved *big.Float `json:"reserved"`
}


type TransactionType string

const (
    Deposit           TransactionType = "deposit"
    Withdrawal        TransactionType = "withdrawal"
    Reserve           TransactionType = "reserve"
    Confirm            TransactionType = "confirm"
    TransferFrom       TransactionType = "transfer_from"
    TransferTo         TransactionType = "transfer_to"
)