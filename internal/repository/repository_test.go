package repository

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_DeleteReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	type args struct {
		ctx        context.Context
		ReservedID int
	}
	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr bool
	}{
		{
			name: "Successful deletion",
			mock: func() {
				mock.ExpectPrepare("DELETE FROM reserved_funds WHERE id = \\$1").
					ExpectExec().
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx:        context.Background(),
				ReservedID: 1,
			},
			wantErr: false,
		},
		{
			name: "Prepare statement error",
			mock: func() {
				mock.ExpectPrepare("DELETE FROM reserved_funds WHERE id = \\$1").
					WillReturnError(fmt.Errorf("prepare error"))
			},
			args: args{
				ctx:        context.Background(),
				ReservedID: 1,
			},
			wantErr: true,
		},
		{
			name: "Exec error",
			mock: func() {
				mock.ExpectPrepare("DELETE FROM reserved_funds WHERE id = \\$1").
					ExpectExec().
					WithArgs(1).
					WillReturnError(fmt.Errorf("exec error"))
			},
			args: args{
				ctx:        context.Background(),
				ReservedID: 1,
			},
			wantErr: true,
		},
		{
			name: "No rows affected",
			mock: func() {
				mock.ExpectPrepare("DELETE FROM reserved_funds WHERE id = \\$1").
					ExpectExec().
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			args: args{
				ctx:        context.Background(),
				ReservedID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.DeleteReservation(tt.args.ctx, tt.args.ReservedID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.DeleteReservation() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepositoryGetUserBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewRepository(db)

	tests := []struct {
		name    string
		mock    func()
		args    int
		wantErr bool
	}{
		{
			name: "Successfully get user balance",
			mock: func() {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT balance FROM users WHERE id = $1")).
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow("100.50"))
			},
			args:    1,
			wantErr: false,
		},
		{
			name: "Prepare statement error",
			mock: func() {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT balance FROM users WHERE id = $1")).
					ExpectQuery().
					WithArgs(1).
					WillReturnError(fmt.Errorf("prepare error"))
			},
			args:    1,
			wantErr: true,
		},
		{
			name: "No rows",
			mock: func() {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT balance FROM users WHERE id = $1")).
					ExpectQuery().
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			args:    1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := repo.GetUserBalance(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
