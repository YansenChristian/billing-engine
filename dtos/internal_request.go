package dtos

type CreateLoanRequestParam struct {
	UserID             int64  `json:"user_id"`
	LoanAmount         string `json:"loan_amount"`
	TenureValue        int    `json:"tenure_value"`
	TenureUnit         int8   `json:"tenure_unit"`
	AnnualInterestRate string `json:"annual_interest_rate"` // inflated by 10^2
}

type GetOutstandingParam struct {
	UserID int64 `json:"user_id"`
	LoanID int64 `json:"loan_id"`
}

type IsDelinquentParam struct {
	LoanID int64 `json:"loan_id"`
}

type MakePaymentParam struct {
	UserID int64  `json:"user_id"`
	LoanID int64  `json:"loan_id"`
	Amount string `json:"amount"`
}
