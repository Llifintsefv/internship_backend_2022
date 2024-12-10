package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"internship_backend_2022/internal/models"
	"math/big"

	_ "github.com/lib/pq"
)

var (
	ErrNoRows = sql.ErrNoRows
)

type Repository interface {
	GetUserBalance(ctx context.Context,userID int) (*big.Float,error)
	CreateUser(ctx context.Context,userID int) (error)
	CreateTransaction(ctx context.Context,userId int, serviceId int,orderId int,amount *big.Float,txType models.TransactionType,descriptions string ) (int,error)
	UpdateUserBalance(ctx context.Context,userID int,amount *big.Float) (*big.Float,error)
	
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func InitDB(connStr string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}


	return db, nil
}

func (r *repository)GetUserBalance(ctx context.Context,userID int) (*big.Float,error) {
	stmt,err := r.db.Prepare("SELECT balance FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}
	defer stmt.Close()

	var balanceStr string
	err = stmt.QueryRowContext(ctx,userID).Scan(&balanceStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	balance, ok := new(big.Float).SetString(balanceStr)
	if !ok {
		return nil, fmt.Errorf("failed to convert balance to big.Float: %w", err)
	}

	return balance, nil
}

func (r *repository)CreateUser(ctx context.Context,userID int) (error) {
	stmt,err := r.db.Prepare("INSERT INTO users (id,balance) VALUES ($1,0.00)")
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	_,err = stmt.ExecContext(ctx,userID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}


func (r *repository)CreateTransaction(ctx context.Context,userId int, serviceId int,orderId int,amount *big.Float,txType models.TransactionType,descriptions string ) (int,error) {
	var transactionsID int
	amountStr := amount.Text('f', 2) 
	stmt,err := r.db.Prepare(`INSERT INTO transactions (user_id,service_id,order_id,amount,type,description)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`)
	if err != nil {
		return 0,fmt.Errorf("failed to create transaction: %w", err)
	}
	defer stmt.Close()
	_ = stmt.QueryRowContext(ctx,userId,serviceId,orderId,amountStr,txType,descriptions).Scan(&transactionsID)
	return transactionsID, nil

}

func (r *repository)UpdateUserBalance(ctx context.Context,userID int,amount *big.Float) (*big.Float,error) {
	tx,err := r.db.BeginTx(ctx,nil)
	if err != nil {
		return nil,fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentBalanceStr string
	err = tx.QueryRowContext(ctx,"SELECT balance FROM users WHERE id = $1 FOR UPDATE",userID).Scan(&currentBalanceStr)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows) {
			return nil,ErrNoRows
		}
		return nil,fmt.Errorf("failed to get user balance: %w", err)
	} 
		currentBalance,ok := new(big.Float).SetString(currentBalanceStr)
		if !ok {
			return nil,fmt.Errorf("failed to convert balance to big.Float: %w", err)
		}
		newBalance := new(big.Float).Add(currentBalance,amount)

		_,err = tx.ExecContext(ctx,"UPDATE users SET balance = $1 WHERE id = $2",newBalance.String(),userID)
		if err != nil {
			return nil,fmt.Errorf("failed to update user balance: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return nil,fmt.Errorf("failed to commit transaction: %w",err)
		}

		return newBalance,nil
	}