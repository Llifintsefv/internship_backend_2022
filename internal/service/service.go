package service

import (
	"context"
	"errors"
	"fmt"
	"internship_backend_2022/internal/models"
	"internship_backend_2022/internal/repository"
	"math/big"
)

type Service interface {
	Deposit(ctx context.Context,request models.DepositRequest) (models.DepositResponse, error)
	GetUserBalance(ctx context.Context,userID int) (*big.Float,error)
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
	fmt.Println(depositResponse)
	return depositResponse, nil
}

func (s *service) GetUserBalance(ctx context.Context, userID int) (*big.Float, error) {
	return s.repository.GetUserBalance(ctx, userID)
}