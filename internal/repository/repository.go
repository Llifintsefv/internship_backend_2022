package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"internship_backend_2022/internal/models"
	"math/big"
	"time"

	_ "github.com/lib/pq"
)

var (
	ErrNoRows = sql.ErrNoRows
)

type Repository interface {
	GetUserBalance(ctx context.Context,userID int) (*big.Float,error)
	GetUserReservedFunds(ctx context.Context,userId int) (*big.Float,error)
	CreateUser(ctx context.Context,userID int) (error)
	CreateTransaction(ctx context.Context,userId int, serviceId int,orderId int,amount *big.Float,txType models.TransactionType,descriptions string ) (int,error)
	UpdateUserBalance(ctx context.Context,userID int,amount *big.Float) (*big.Float,error)
	ReserveFunds(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) (int,error)
	DeleteReservation(ctx context.Context,ReservedID int) (error)
	DeleteReservationByServiceAndOrder(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) error
	GetReserveFundsByServiceAndOrder(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) (bool,error)
	AddRevenueRecord(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) error
	Transfer(ctx context.Context,fromUserId int,toUserId int,amount *big.Float) error
	GetMonthlyReportData(ctx context.Context, year, month int) ([]models.MonthlyReportData, error)
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


func (r *repository)GetUserReservedFunds(ctx context.Context,userId int) (*big.Float,error) {
    var totalReservedStr string
    err := r.db.QueryRowContext(ctx, `
        SELECT COALESCE(SUM(amount), '0')
        FROM reserved_funds
        WHERE user_id = $1`,
        userId,
    ).Scan(&totalReservedStr)
    if err != nil {
        return nil, fmt.Errorf("failed to get reserved balance: %w", err)
    }

    totalReserved, ok := new(big.Float).SetString(totalReservedStr)
    if !ok {
        return nil, fmt.Errorf("failed to parse total reserved balance: %s", totalReservedStr)
    }

    return totalReserved, nil
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
	stmt,err := r.db.Prepare(`INSERT INTO transactions (user_id,service_id,order_id,amount,type,description)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`)
	if err != nil {
		return 0,fmt.Errorf("failed to create transaction: %w", err)
	}
	defer stmt.Close()
	_ = stmt.QueryRowContext(ctx,userId,serviceId,orderId,amount.Text('f', 2) ,txType,descriptions).Scan(&transactionsID)
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

func (r *repository)ReserveFunds(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) (int,error) {
	var ReservedID int
	stmt,err := r.db.Prepare(`INSERT INTO reserved_funds (user_id,service_id,order_id,amount)
	VALUES ($1,$2,$3,$4)
	RETURNING id`)
	if err != nil {
		return 0,fmt.Errorf("failed to create transaction: %w", err)
	}
	defer stmt.Close()
	_ = stmt.QueryRowContext(ctx,userId,serviceId,orderId,amount.Text('f', 2) ).Scan(&ReservedID)
	return ReservedID, nil
}

func (r *repository)DeleteReservation(ctx context.Context,ReservedID int) (error) {
	stmt,err := r.db.Prepare("DELETE FROM reserved_funds WHERE id = $1")
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	_,err = stmt.ExecContext(ctx,ReservedID)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	return nil
	}


func (r *repository)GetReserveFundsByServiceAndOrder(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) (bool,error) {
	var Exist bool
	stmt,err := r.db.Prepare("SELECT EXISTS(SELECT 1 FROM reserved_funds WHERE user_id = $1 AND service_id = $2 AND order_id = $3 AND amount = $4)")
	if err != nil {
		return false,fmt.Errorf("failed to delete reservation: %w", err)
	}
	_ = stmt.QueryRowContext(ctx,userId,serviceId,orderId,amount.Text('f', 2) ).Scan(&Exist)
	if err != nil {
		return false,fmt.Errorf("failed to delete reservation: %w", err)
	}
	return Exist, nil
}

func (r *repository)DeleteReservationByServiceAndOrder(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) error {
	stmt,err := r.db.Prepare("DELETE FROM reserved_funds WHERE user_id = $1 AND service_id = $2 AND order_id = $3 AND amount = $4")
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	_,err = stmt.ExecContext(ctx,userId,serviceId,orderId,amount.Text('f', 2))
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	return nil
}


func (r *repository)AddRevenueRecord(ctx context.Context,userId int,serviceId int,orderId int,amount *big.Float) error {
	stmt,err := r.db.Prepare("INSERT INTO revenue_report (user_id,service_id,order_id,revenue) VALUES ($1,$2,$3,$4)")
	if err != nil {
		return fmt.Errorf("failed to add revenue record: %w", err)
	}
	_,err = stmt.ExecContext(ctx,userId,serviceId,orderId,amount.Text('f', 2))
	if err != nil {
		return fmt.Errorf("failed to add revenue record: %w", err)
	}
	return nil

}

func (r *repository)Transfer (ctx context.Context,fromUserId int,toUserId int,amount *big.Float) error {
	tx,err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	_,err = tx.ExecContext(ctx,"UPDATE users SET balance = balance + $1 WHERE id = $2",amount.Text('f', 2),toUserId)
	if err != nil {
		return fmt.Errorf("failed to update toUserId balance: %w", err)
	}
	_,err = tx.ExecContext(ctx,"UPDATE users SET balance = balance - $1 WHERE id = $2",amount.Text('f', 2),fromUserId)
	if err != nil {
		return fmt.Errorf("failed to update fromUserId balance: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w",err)
	}
	return nil
}

func (r *repository) GetMonthlyReportData(ctx context.Context, year, month int) ([]models.MonthlyReportData, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0) 

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT service_id, SUM(revenue) AS total_revenue
		FROM revenue_report
		WHERE created_at >= $1 AND created_at < $2
		GROUP BY service_id
	`)
	if err != nil {
		return []models.MonthlyReportData{}, fmt.Errorf("database prepare error: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, startOfMonth, endOfMonth)
	if err != nil {
		return []models.MonthlyReportData{}, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var reportData []models.MonthlyReportData
	for rows.Next() {
		var data models.MonthlyReportData
		err := rows.Scan(&data.ServiceName, &data.TotalRevenue)
		if err != nil {
			return []models.MonthlyReportData{}, fmt.Errorf("database scan error: %w", err)
		}
		reportData = append(reportData, data)
	}

	if err := rows.Err(); err != nil {
		return []models.MonthlyReportData{}, fmt.Errorf("database rows error: %w", err)
	}

	return reportData, nil
}