package models

import (
	"math/big"
	"time"
)

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

type TransferRequest struct {
    FromUserID int      `json:"from_user_id"`
    ToUserID   int      `json:"to_user_id"`
    Amount     *big.Float `json:"amount"`
}

type TransferResponse struct {
    Status        string     `json:"status"`
    Message       string     `json:"message"`
    TransactionID int        `json:"transaction_id"`
    UserToBalance *big.Float `json:"user_to_balance"`
    UserFromBalance *big.Float `json:"user_from_balance"`
}

type MonthlyReportRequest struct {
    Month int `json:"month"`
    Year int `json:"year"`
}

type MonthlyReportData struct {
	ServiceName   string
	TotalRevenue float64
}

type MonthlyReportResponse struct {
    FilePath string
}


type Transaction struct {
    ID          int             `json:"id"`
    UserID      int             `json:"user_id"`
    ServiceID   int             `json:"service_id,omitempty"` 
    OrderID     int             `json:"order_id,omitempty"`   
    Amount      *big.Float      `json:"amount"`
    Type        TransactionType `json:"type"`
    Description string          `json:"description"`
    CreatedAt   time.Time       `json:"created_at"`
}

type TransactionsResponse struct {
    Transactions []Transaction `json:"transactions"`
    Total int `json:"total"`
    Page int `json:"page"`
    Limit int `json:"limit"`
}

type TransactionRequest struct{
    UserId int `json:"user_id"`
    Page int `json:"page"`
    Limit int `json:"limit"`
    SortBy string `json:"sort_by"`
    SortOrder string `json:"sort_order"`
}




type TransactionType string

const (
    Deposit           TransactionType = "deposit"
    Withdrawal        TransactionType = "withdrawal"
    Reserve           TransactionType = "reserve"
    Confirm            TransactionType = "confirm"
    Transfer             TransactionType = "transfer"
    
)
