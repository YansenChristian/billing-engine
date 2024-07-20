package dtos

import (
	"loan-payment/constants"

	"github.com/shopspring/decimal"
)

type UserModel struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	CreatedAt uint64 `db:"created_at"`
	UpdatedAt uint64 `db:"updated_at"`
	DeletedAt uint64 `db:"deleted_at"`
}

func (m *UserModel) GetAll() []interface{} {
	return []interface{}{
		&m.ID,
		&m.Name,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	}
}

func (m *UserModel) GetTableName() string {
	return "users_tab"
}

type LoanRequestModel struct {
	ID                  int64                `db:"id"`
	UserID              int64                `db:"user_id"`
	LoanAmount          decimal.Decimal      `db:"loan_amount"`
	PrincipalPaidAmount decimal.Decimal      `db:"principal_paid_amount"`
	InterestPaidAmount  decimal.Decimal      `db:"interest_paid_amount"`
	DisbursementTime    int64                `db:"disbursement_time"`
	TenureValue         int                  `db:"tenure_value"`
	TenureUnit          constants.TenureUnit `db:"tenure_unit"`
	Status              constants.LoanStatus `db:"status"`
	AnnualInterestRate  decimal.Decimal      `db:"annual_interest_rate"`
	CreatedAt           int64                `db:"created_at"`
	UpdatedAt           int64                `db:"updated_at"`
	DeletedAt           int64                `db:"deleted_at"`
}

func (m *LoanRequestModel) GetAll() []interface{} {
	return []interface{}{
		&m.ID,
		&m.UserID,
		&m.LoanAmount,
		&m.PrincipalPaidAmount,
		&m.InterestPaidAmount,
		&m.DisbursementTime,
		&m.TenureValue,
		&m.TenureUnit,
		&m.Status,
		&m.AnnualInterestRate,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	}
}

func (m *LoanRequestModel) GetTableName() string {
	return "loan_requests_tab"
}

type BillingModel struct {
	ID                 int64                   `db:"id"`
	BillingID          string                  `db:"billing_id"`
	LoanID             int64                   `db:"loan_id"`
	PaymentID          int64                   `db:"payment_id"`
	RecurringIndex     int                     `db:"recurring_index"`
	PrincipalAmount    decimal.Decimal         `db:"principal_amount"`
	InterestAmount     decimal.Decimal         `db:"interest_amount"`
	TotalAmount        decimal.Decimal         `db:"total_amount"`
	DueTime            int64                   `db:"due_time"`
	PaymentCompletedAt int64                   `db:"payment_completed_at"`
	Status             constants.PaymentStatus `db:"status"`
	CreatedAt          int64                   `db:"created_at"`
	UpdatedAt          int64                   `db:"updated_at"`
	DeletedAt          int64                   `db:"deleted_at"`
}

func (m *BillingModel) GetAll() []interface{} {
	return []interface{}{
		&m.ID,
		&m.BillingID,
		&m.LoanID,
		&m.PaymentID,
		&m.RecurringIndex,
		&m.PrincipalAmount,
		&m.InterestAmount,
		&m.TotalAmount,
		&m.DueTime,
		&m.PaymentCompletedAt,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	}
}

func (m *BillingModel) GetTableName() string {
	return "billings_tab"
}

type PaymentModel struct {
	ID        int64           `db:"id"`
	UserID    int64           `db:"user_id"`
	Amount    decimal.Decimal `db:"amount"`
	CreatedAt int64           `db:"created_at"`
	UpdatedAt int64           `db:"updated_at"`
	DeletedAt int64           `db:"deleted_at"`
}

func (m *PaymentModel) GetAll() []interface{} {
	return []interface{}{
		&m.ID,
		&m.UserID,
		&m.Amount,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	}
}

func (m *PaymentModel) GetTableName() string {
	return "payments_tab"
}

type LoanRequestHistory struct {
	ID                  int64                `db:"id"`
	LoanID              int64                `db:"loan_id"`
	PrincipalPaidAmount decimal.Decimal      `db:"principal_paid_amount"`
	InterestPaidAmount  decimal.Decimal      `db:"interest_paid_amount"`
	Status              constants.LoanStatus `db:"status"`
	CreatedAt           int64                `db:"created_at"`
}

func (m *LoanRequestHistory) GetAll() []interface{} {
	return []interface{}{
		&m.ID,
		&m.LoanID,
		&m.PrincipalPaidAmount,
		&m.InterestPaidAmount,
		&m.Status,
		&m.CreatedAt,
	}
}

func (m *LoanRequestHistory) GetTableName() string {
	return "loan_request_histories_tab"
}

type BillingHistoryModel struct {
	ID                 int64                   `db:"id"`
	BillingID          string                  `db:"billing_id"`
	PaymentCompletedAt int64                   `db:"payment_completed_at"`
	Status             constants.PaymentStatus `db:"status"`
	CreatedAt          int64                   `db:"created_at"`
}

func (m *BillingHistoryModel) GetAll() []interface{} {
	return []interface{}{
		&m.ID,
		&m.BillingID,
		&m.PaymentCompletedAt,
		&m.Status,
		&m.CreatedAt,
	}
}

func (m *BillingHistoryModel) GetTableName() string {
	return "billing_histories_tab"
}
