package service

import (
	"context"
	"errors"
	"fmt"
	"internship_backend_2022/internal/models"
	"internship_backend_2022/internal/repository"
)

type Service interface {
	Deposit(ctx context.Context,request models.DepositRequest) (models.DepositResponse, error)
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Deposit(ctx context.Context,DepositRequest models.DepositRequest) (models.DepositResponse,error) {

	if DepositRequest.Amount.Cmp(big.NewFloat(0)) <= 0 {
		return models.DepositResponse{},errors.New("amount must be greater than 0")
	}

	_,err := s.repository.GetUserBalance(ctx,DepositRequest.UserID)
	if err != nil {
		if errors.Is(err,repository.ErrNoRows) {
			err = s.repository.CreateUser(ctx,DepositRequest.UserID)
			if err != nil {
				return models.DepositResponse{},fmt.Errorf("failed to create user: %w", err)
			}
	} else {
		return models.DepositResponse{},fmt.Errorf("failed to get user balance: %w", err)
	}

	transactionId,err := s.repository.CreateTransaction(tx,DepositRequest.UserID,DepositRequest.ServiceID,DepositRequest.OrderID,DepositRequest.Amount,DepositRequest.Type,DepositRequest.Description)
	if err != nil {
		return models.DepositResponse{},fmt.Errorf("failed to create transaction: %w", err)
	}

	newBalance,err := s.repository.UpdateUserBalance(ctx,DepositRequest.UserID,DepositRequest.Amount)

	
	
}