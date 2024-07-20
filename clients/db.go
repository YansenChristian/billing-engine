package clients

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"loan-payment/configs"
	"loan-payment/constants"
	"loan-payment/dtos"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var dbClient *sqlx.DB

func getDatabase() *sqlx.DB {
	if dbClient != nil {
		return dbClient
	}

	conf := configs.Get()

	if conf.DBMaster == nil {
		log.Fatalf("failed to get DB config")
	}

	db, err := sqlx.Open("mysql", conf.DBMaster.ConnectionString)
	if err != nil {
		log.Fatalf("failed to open DB master connection. %+v", err)
	}
	db.SetMaxIdleConns(conf.DBMaster.MaxIdle)
	db.SetMaxOpenConns(conf.DBMaster.MaxOpen)
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping DB master. %+v", err)
	}

	dbClient = db
	return dbClient
}

func DBBeginTransaction(ctx context.Context) (*sqlx.Tx, error) {
	var err error
	if tx, err := getDatabase().BeginTxx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	}); err == nil {
		return tx, nil
	}
	return nil, errors.Wrap(err, "failed to start tx")
}

func DBRollbackTransaction(tx *sqlx.Tx) error {
	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, "failed to rollback tx")
	}
	return nil
}

func DBCommitTransaction(tx *sqlx.Tx) error {
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}
	return nil
}

func DBGetUserByID(ctx context.Context, userID int64) (*dtos.UserModel, error) {
	var (
		userModel dtos.UserModel
		err       error

		args = []interface{}{
			userID,
		}
		query = `
			SELECT 
				id, name, 
				created_at, updated_at, deleted_at
			FROM users_tab
			WHERE 
			    id = ? 
				AND deleted_at = 0
			LIMIT 1`
	)

	if err = getDatabase().QueryRowContext(ctx, query, args...).Scan(userModel.GetAll()...); err == sql.ErrNoRows {
		return nil, constants.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}
	return &userModel, nil
}

func DBGetLoanRequestByID(ctx context.Context, loanID int64) (*dtos.LoanRequestModel, error) {
	var (
		loanRequestModel dtos.LoanRequestModel
		err              error

		args = []interface{}{
			loanID,
		}
		query = `
			SELECT 
				id, user_id,
				loan_amount, principal_paid_amount, interest_paid_amount,
			    disbursement_time, tenure_value, tenure_unit, 
			    status, annual_interest_rate,
				created_at, updated_at, deleted_at
			FROM loan_requests_tab
			WHERE 
			    id = ? 
			  	AND deleted_at = 0
			LIMIT 1`
	)

	if err = getDatabase().QueryRowContext(ctx, query, args...).Scan(loanRequestModel.GetAll()...); err == sql.ErrNoRows {
		return nil, constants.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}
	return &loanRequestModel, nil
}

func DBGetLoanRequestByIDAndUserIDForUpdate(ctx context.Context, tx *sqlx.Tx, loanID, userID int64) (*dtos.LoanRequestModel, error) {
	var (
		loanRequestModel dtos.LoanRequestModel
		err              error

		args = []interface{}{
			loanID,
			userID,
		}
		query = `
			SELECT 
				id, user_id,
				loan_amount, principal_paid_amount, interest_paid_amount,
				disbursement_time, tenure_value, tenure_unit, 
			    status, annual_interest_rate,
				created_at, updated_at, deleted_at
			FROM loan_requests_tab
			WHERE 
			    id = ? 
			  	AND user_id = ?
			  	AND deleted_at = 0
			LIMIT 1
			FOR UPDATE`
	)

	if err = tx.QueryRowContext(ctx, query, args...).Scan(loanRequestModel.GetAll()...); err == sql.ErrNoRows {
		return nil, constants.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}
	return &loanRequestModel, nil
}

func DBInsertLoanRequest(ctx context.Context, tx *sqlx.Tx, model *dtos.LoanRequestModel) (int64, error) {
	var (
		loanID int64
		res    sql.Result
		err    error

		now   = time.Now().UnixMilli()
		query = `INSERT INTO 
			loan_requests_tab 
			(user_id, 
			 loan_amount, principal_paid_amount,interest_paid_amount, 
			 disbursement_time, tenure_value, tenure_unit, 
			 status, annual_interest_rate,
			 created_at, updated_at, deleted_at) VALUES 
			(?,
			 ?, ?, ?,
			 ?, ?, ?,
			 ?, ?,
			 ?, ?, ?)`
	)

	if tx != nil {
		res, err = tx.ExecContext(ctx, query,
			model.UserID,
			model.LoanAmount, model.PrincipalPaidAmount, model.InterestPaidAmount,
			model.DisbursementTime, model.TenureValue, model.TenureUnit,
			model.Status, model.AnnualInterestRate,
			now, now, 0)
	} else {
		res, err = getDatabase().ExecContext(ctx, query,
			model.UserID,
			model.LoanAmount, model.PrincipalPaidAmount, model.InterestPaidAmount,
			model.DisbursementTime, model.TenureValue, model.TenureUnit,
			model.Status, model.AnnualInterestRate,
			now, now, 0)
	}

	if err != nil {
		return 0, err
	}

	loanID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return loanID, nil
}

func DBInsertPayment(ctx context.Context, tx *sqlx.Tx, model *dtos.PaymentModel) (int64, error) {
	var (
		paymentID int64
		res       sql.Result
		err       error

		now   = time.Now().UnixMilli()
		query = `INSERT INTO 
			payments_tab 
			(user_id, amount, 
			 created_at, updated_at, deleted_at) VALUES 
			(?, ?,
			 ?, ?, ?)`
	)

	if tx != nil {
		res, err = tx.ExecContext(ctx, query,
			model.UserID, model.Amount,
			now, now, 0)
	} else {
		res, err = getDatabase().ExecContext(ctx, query,
			model.UserID, model.Amount,
			now, now, 0)
	}

	if err != nil {
		return 0, err
	}

	paymentID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return paymentID, nil
}

func DBBatchInsertBillings(ctx context.Context, tx *sqlx.Tx, models []dtos.BillingModel) error {
	var (
		err error

		placeholders = make([]string, 0, len(models))
		args         = make([]interface{}, 0)
	)

	queryTemplate := `INSERT INTO billings_tab 
		(billing_id, loan_id, payment_id, recurring_index,
		principal_amount, interest_amount, total_amount,
		due_time, payment_completed_at, status,
		created_at, updated_at, deleted_at) VALUES %s`
	insertPlaceholder := `(
		?, ?, ?, ?,
		?, ?, ?,
		?, ?, ?,
	    ?, ?, ?)`

	for _, model := range models {
		placeholders = append(placeholders, insertPlaceholder)
		args = append(args,
			model.BillingID, model.LoanID, model.PaymentID, model.RecurringIndex,
			model.PrincipalAmount, model.InterestAmount, model.TotalAmount,
			model.DueTime, model.PaymentCompletedAt, model.Status,
			model.CreatedAt, model.UpdatedAt, model.DeletedAt,
		)
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ",")), args...)
	} else {
		_, err = getDatabase().ExecContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ",")), args)
	}
	return err
}

func DBBatchInsertLoanRequestHistories(ctx context.Context, tx *sqlx.Tx, models []dtos.LoanRequestHistory) error {
	var (
		err error

		placeholders = make([]string, 0, len(models))
		args         = make([]interface{}, 0)
	)

	queryTemplate := `INSERT INTO loan_request_histories_tab 
		(loan_id, principal_paid_amount, 
		interest_paid_amount, status, created_at) VALUES %s`
	insertPlaceholder := `(
		?, ?,
		?, ?, ?)`

	for _, model := range models {
		placeholders = append(placeholders, insertPlaceholder)
		args = append(args,
			model.LoanID,
			model.PrincipalPaidAmount,
			model.InterestPaidAmount,
			model.Status,
			model.CreatedAt,
		)
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ",")), args...)
	} else {
		_, err = getDatabase().ExecContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ",")), args)
	}
	return err
}

func DBBatchInsertBillingHistories(ctx context.Context, tx *sqlx.Tx, models []dtos.BillingHistoryModel) error {
	var (
		err error

		placeholders = make([]string, 0, len(models))
		args         = make([]interface{}, 0)
	)

	queryTemplate := `INSERT INTO billing_histories_tab 
		(billing_id, payment_completed_at, status, created_at) VALUES %s`
	insertPlaceholder := `(?, ?, ?, ?)`

	for _, model := range models {
		placeholders = append(placeholders, insertPlaceholder)
		args = append(args,
			model.BillingID,
			model.PaymentCompletedAt,
			model.Status,
			model.CreatedAt,
		)
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ",")), args...)
	} else {
		_, err = getDatabase().ExecContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ",")), args)
	}
	return err
}

func DBGetOverdueBillings(ctx context.Context, loanID int64) ([]dtos.BillingModel, error) {
	var (
		billingModels []dtos.BillingModel
		err           error

		args = []interface{}{
			loanID,
			constants.PaymentStatus_Pending,
			time.Now().UnixMilli(),
		}
		query = `
			SELECT 
				id, billing_id,
			    loan_id, payment_id, recurring_index,
				principal_amount, interest_amount, total_amount,
				due_time, payment_completed_at, status,
				created_at, updated_at, deleted_at
			FROM billings_tab
			WHERE 
			  	loan_id = ?
			  	AND status = ? 
			  	AND due_time < ?`
	)

	if err = getDatabase().SelectContext(ctx, &billingModels, query, args...); err != nil {
		return nil, err
	}
	return billingModels, nil
}

func DBGetPendingBillingsWithDueTimeByLoanIdForUpdate(ctx context.Context, tx *sqlx.Tx, loanID, dueTime int64) ([]dtos.BillingModel, error) {
	var (
		models []dtos.BillingModel

		args = []interface{}{
			loanID,
			constants.PaymentStatus_Pending,
			dueTime,
		}
		query = `
			SELECT 
				id, billing_id,
				loan_id, payment_id, recurring_index,
				principal_amount, interest_amount, total_amount,
				due_time, payment_completed_at, status,
				created_at, updated_at, deleted_at
			FROM billings_tab
			WHERE 
			    loan_id = ?
			  	AND status = ? 
			  	AND due_time <= ?
			  	AND deleted_at = 0
			FOR UPDATE`
	)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var model dtos.BillingModel
		if err = rows.Scan(model.GetAll()...); err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return models, nil
}

func DBBulkPayBillingsByIDs(ctx context.Context, tx *sqlx.Tx, ids []int64, paymentID int64) error {
	var (
		err error

		now = time.Now().UnixMilli()
	)

	query := `UPDATE billings_tab 
		SET payment_id = :payment_id,
			payment_completed_at = :now,
			status = :status,
		    updated_at = :now
		WHERE id IN(:ids)`
	arg := map[string]interface{}{
		"payment_id": paymentID,
		"now":        now,
		"status":     constants.PaymentStatus_Completed,
		"ids":        ids,
	}

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = tx.Rebind(query)
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = getDatabase().ExecContext(ctx, query, args)
	}
	return err
}

func DBUpdateLoanRequestPaymentByID(ctx context.Context, tx *sqlx.Tx, loanID int64, principalPaid, interestPaid decimal.Decimal, status constants.LoanStatus) error {
	var err error

	query := `UPDATE loan_requests_tab 
		SET principal_paid_amount = principal_paid_amount + ?,
			interest_paid_amount = interest_paid_amount + ?,
			status = ?,
		    updated_at = ?
		WHERE id = ?`
	args := []interface{}{
		principalPaid,
		interestPaid,
		status,
		time.Now().UnixMilli(),
		loanID,
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = getDatabase().ExecContext(ctx, query, args)
	}
	return err
}
