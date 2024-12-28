package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"internship_backend_2022/internal/models"
	"internship_backend_2022/internal/repository"
	"math/big"
	"os"
	"strconv"
)

type Service interface {
	Deposit(ctx context.Context,request models.DepositRequest) (models.DepositResponse, error)
	GetUserBalance(ctx context.Context,userID int) (models.BalanceResponse,error)
	Reserve(ctx context.Context, request models.ReserveRequest) (models.ReserveResponse, error)
	Confirm(ctx context.Context,request models.ConfirmRequest) (models.ConfirmResponse, error)
	Transfer(ctx context.Context,request models.TransferRequest) (models.TransferResponse,error)
	MonthlyReport(ctx context.Context,MonthlyReportRequest models.MonthlyReportRequest) (models.MonthlyReportResponse, error)
	Transactions(ctx context.Context, request models.TransactionRequest) (models.TransactionsResponse,error)
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Deposit(ctx context.Context, depositRequest models.DepositRequest) (models.DepositResponse, error) {
	fmt.Println(depositRequest)
	if depositRequest.Amount == nil {
        return models.DepositResponse{}, errors.New("amount is required")
    }

    if depositRequest.Amount.Cmp(big.NewFloat(0)) <= 0 {
        return models.DepositResponse{}, errors.New("amount must be greater than 0")
    }
	fmt.Println(depositRequest)
	_, err := s.repository.GetUserBalance(ctx, depositRequest.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			err = s.repository.CreateUser(ctx, depositRequest.UserID)
			if err != nil {
				return models.DepositResponse{}, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return models.DepositResponse{}, fmt.Errorf("failed to get user balance: %w", err)
		}
	}

	transactionId, err := s.repository.CreateTransaction(
		ctx,
		depositRequest.UserID,
		0,
		0,
		depositRequest.Amount,
		models.Deposit,
		"deposit",
	)
	if err != nil {
		return models.DepositResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	newBalance, err := s.repository.UpdateUserBalance(ctx, depositRequest.UserID, depositRequest.Amount)
	if err != nil {
		return models.DepositResponse{}, fmt.Errorf("failed to update user balance: %w", err)
	}

	depositResponse := models.DepositResponse{
		Status:        "success",
		Message:       "funds deposited successfully",
		Balance:       newBalance,
		TransactionID: transactionId,
	}
	return depositResponse, nil
}

func (s *service) GetUserBalance(ctx context.Context, userID int) (models.BalanceResponse, error) {

	balance,err := s.repository.GetUserBalance(ctx, userID)
	if err != nil {
		return models.BalanceResponse{}, fmt.Errorf("failed to get user balance: %w", err)
	}

	reserved,err := s.repository.GetUserReservedFunds(ctx, userID)
	if err != nil {
		return models.BalanceResponse{}, fmt.Errorf("failed to get user reserved funds: %w", err)
	}

	return models.BalanceResponse{Balance: balance, Reserved: reserved},nil 
}

func (s *service) Reserve (ctx context.Context, reserveRequest models.ReserveRequest) (models.ReserveResponse, error) {
	if reserveRequest.Amount == nil {
        return models.ReserveResponse{}, errors.New("amount is required")
    }

    if reserveRequest.Amount.Cmp(big.NewFloat(0)) <= 0 {
        return models.ReserveResponse{}, errors.New("amount must be greater than 0")
    }

	userBalance,err := s.repository.GetUserBalance(ctx,reserveRequest.UserID) 
	if err != nil {
		return models.ReserveResponse{}, fmt.Errorf("failed to get user balance: %w", err)
	}

	if userBalance.Cmp(reserveRequest.Amount) < 0 {
        return models.ReserveResponse{}, errors.New("insufficient funds")
    }

	reservedId,err := s.repository.ReserveFunds(ctx,reserveRequest.UserID,reserveRequest.ServiceID,reserveRequest.OrderID,reserveRequest.Amount)
	if err != nil {
		return models.ReserveResponse{}, fmt.Errorf("failed to reserve funds: %w", err)
	}

	NewUserBalance,err := s.repository.UpdateUserBalance(ctx,reserveRequest.UserID,new(big.Float).Neg(reserveRequest.Amount))
	if err != nil {
		if err := s.repository.DeleteReservation(ctx, reservedId); err != nil {
            return models.ReserveResponse{}, fmt.Errorf("failed to rollback reservation after balance update error: %w", err)
        }
        return models.ReserveResponse{}, fmt.Errorf("failed to update user balance: %w", err)
	}

	transactionId, err := s.repository.CreateTransaction(
		ctx,
		reserveRequest.UserID,
		reserveRequest.ServiceID,
		reserveRequest.OrderID,
		new(big.Float).Neg(reserveRequest.Amount),
		models.Reserve,
		"reserve",
	)
	
	if err != nil {
		return models.ReserveResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	fmt.Println(transactionId,NewUserBalance,reservedId)

	ReserveResponse := models.ReserveResponse{
		Status:        "success",
		Message:       "funds reserved successfully",
		Balance:       NewUserBalance,
		Reserved:      reserveRequest.Amount,
		TransactionID: transactionId,
	}
	return ReserveResponse, nil
}

func (s *service) Confirm (ctx context.Context, ConfirmRequest models.ConfirmRequest) (models.ConfirmResponse,error) {
	if ConfirmRequest.Amount == nil {
        return models.ConfirmResponse{}, errors.New("amount is required")
    }

    if ConfirmRequest.Amount.Cmp(big.NewFloat(0)) <= 0 {
        return models.ConfirmResponse{}, errors.New("amount must be greater than 0")
    }

	ReserveExist,err := s.repository.GetReserveFundsByServiceAndOrder(ctx,ConfirmRequest.UserID,ConfirmRequest.ServiceID,ConfirmRequest.OrderID,ConfirmRequest.Amount)
	if err != nil {
		return models.ConfirmResponse{}, fmt.Errorf("failed to check reservation existence: %w", err)
	}
	if !ReserveExist {
		return models.ConfirmResponse{}, errors.New("no corresponding reservation found")
	}

	err = s.repository.DeleteReservationByServiceAndOrder(ctx,ConfirmRequest.UserID,ConfirmRequest.ServiceID,ConfirmRequest.OrderID,ConfirmRequest.Amount)
	if err != nil {
		return models.ConfirmResponse{}, fmt.Errorf("failed to delete reservation: %w", err)
	}

	
	transactionId, err := s.repository.CreateTransaction(
		ctx,
		ConfirmRequest.UserID,
		ConfirmRequest.ServiceID,
		ConfirmRequest.OrderID,
		new(big.Float).Neg(ConfirmRequest.Amount),
		models.Confirm,
		"confirm",
	)
	if err != nil {
		return models.ConfirmResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	err = s.repository.AddRevenueRecord(ctx,ConfirmRequest.UserID,ConfirmRequest.ServiceID,ConfirmRequest.OrderID,ConfirmRequest.Amount)
	if err != nil {
		return models.ConfirmResponse{}, fmt.Errorf("failed to add revenue record: %w", err)
	}



	ConfirmResponse := models.ConfirmResponse{
		Status:        "success",
		Message:       "funds confirmed successfully",
		TransactionID: transactionId,
	}

	
	return ConfirmResponse, nil
}

func (s *service)Transfer(ctx context.Context,transferRequest models.TransferRequest) (models.TransferResponse,error) {
	if transferRequest.Amount == nil {
        return models.TransferResponse{}, errors.New("amount is required")
    }

    if transferRequest.Amount.Cmp(big.NewFloat(0)) <= 0 {
        return models.TransferResponse{}, errors.New("amount must be greater than 0")
    }

	if transferRequest.FromUserID == transferRequest.ToUserID {
        return models.TransferResponse{}, errors.New("cannot transfer to self")
    }

	FromUserBalance,err := s.repository.GetUserBalance(ctx,transferRequest.FromUserID)
	if err != nil {
		return models.TransferResponse{}, fmt.Errorf("failed to get user balance: %w", err)
	}

	if FromUserBalance.Cmp(transferRequest.Amount) < 0 {
		return models.TransferResponse{}, errors.New("insufficient funds")
	}

	err = s.repository.Transfer(ctx,transferRequest.FromUserID,transferRequest.ToUserID,transferRequest.Amount)
	if err != nil {
		return models.TransferResponse{}, fmt.Errorf("failed to transfer funds: %w", err)
	}

	transactionId, err := s.repository.CreateTransaction(
		ctx,
		transferRequest.FromUserID,
		transferRequest.ToUserID,
		0,
		transferRequest.Amount,
		models.Transfer,
		"transfer",
	)
	if err != nil {
		return models.TransferResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	NewUserToBalance ,err := s.repository.GetUserBalance(ctx,transferRequest.ToUserID)
	if err != nil {
		return models.TransferResponse{}, fmt.Errorf("failed to get user balance: %w", err)
	}

	NewUserFromBalance ,err := s.repository.GetUserBalance(ctx,transferRequest.FromUserID)
	if err != nil {
		return models.TransferResponse{}, fmt.Errorf("failed to get user balance: %w", err)
	}

	TransferResponse := models.TransferResponse{
		Status:        "success",
		Message:       "funds transferred successfully",
		TransactionID: transactionId,
		UserToBalance: NewUserToBalance,
		UserFromBalance: NewUserFromBalance,
	}

	return TransferResponse, nil

}

func (s *service) MonthlyReport(ctx context.Context,MonthlyReportRequest models.MonthlyReportRequest) (models.MonthlyReportResponse, error) {
	reportData,err := s.repository.GetMonthlyReportData(ctx,MonthlyReportRequest.Year,MonthlyReportRequest.Month)
	if err != nil {
		return models.MonthlyReportResponse{}, fmt.Errorf("failed to get monthly report data: %w", err)
	}

	tempFile, err := os.CreateTemp("", "monthly_report_*.csv")
	if err != nil {
		return models.MonthlyReportResponse{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	writer := csv.NewWriter(tempFile)
	defer writer.Flush()

	err = writer.Write([]string{"Service Name", "Total Revenue"})
	if err != nil {
		return models.MonthlyReportResponse{}, fmt.Errorf("failed to write csv header: %w", err)
	}

	for _, data := range reportData {
		err = writer.Write([]string{data.ServiceId, strconv.FormatFloat(data.TotalRevenue, 'f', 2, 64)})
		if err != nil {
			return models.MonthlyReportResponse{}, fmt.Errorf("failed to write csv row: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return models.MonthlyReportResponse{}, fmt.Errorf("csv writer error: %w", err)
	}

	return models.MonthlyReportResponse{FilePath: tempFile.Name()}, nil
}

func (s *service) Transactions(ctx context.Context,TransactionsRequest models.TransactionRequest) (models.TransactionsResponse, error) {
	Transactions,total,err := s.repository.GetTransactions(ctx,TransactionsRequest.UserId,TransactionsRequest.Page,TransactionsRequest.Limit,TransactionsRequest.SortBy,TransactionsRequest.SortOrder)
	if err != nil {
		return models.TransactionsResponse{},fmt.Errorf("failed to get transactions: %w", err)
	}

	TransactionsResponse := models.TransactionsResponse{
		Transactions: Transactions,
		Total:        total,
		Page:         TransactionsRequest.Page,
		Limit:        TransactionsRequest.Limit,
	}
	return TransactionsResponse, nil
}