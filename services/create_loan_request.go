package services

import (
	"context"
	"fmt"
	"time"

	"loan-payment/clients"
	"loan-payment/constants"
	"loan-payment/dtos"
	"loan-payment/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func CreateLoanRequest(ctx context.Context, param dtos.CreateLoanRequestParam) (int64, error) {
	if err := validateLoanRequest(ctx, param); err != nil {
		return 0, err
	}

	now := time.Now().UnixMilli()
	loanAmount, _ := decimal.NewFromString(param.LoanAmount)
	annualInterestRate, _ := decimal.NewFromString(param.AnnualInterestRate)

	txn, err := clients.DBBeginTransaction(ctx)
	if err != nil {
		return 0, err
	}
	defer clients.DBRollbackTransaction(txn)

	var loanModel = dtos.LoanRequestModel{
		UserID:              param.UserID,
		LoanAmount:          loanAmount,
		PrincipalPaidAmount: decimal.NewFromUint64(0),
		InterestPaidAmount:  decimal.NewFromUint64(0),
		DisbursementTime:    now, // assume that all loan request is disbursed that day
		TenureValue:         param.TenureValue,
		TenureUnit:          constants.TenureUnit(param.TenureUnit),
		Status:              constants.LoanStatus_InRepayment,
		AnnualInterestRate:  annualInterestRate,
	}
	loanID, err := clients.DBInsertLoanRequest(ctx, txn, &loanModel)
	if err != nil {
		return 0, err
	}
	loanModel.ID = loanID

	billingModels, billingHistories, err := createRepaymentSchedule(loanModel)
	if err != nil {
		return 0, err
	}
	if err = clients.DBBatchInsertBillings(ctx, txn, billingModels); err != nil {
		return 0, err
	}

	if err = clients.DBBatchInsertLoanRequestHistories(ctx, txn, []dtos.LoanRequestHistory{
		{
			LoanID:              loanID,
			PrincipalPaidAmount: decimal.Zero,
			InterestPaidAmount:  decimal.Zero,
			Status:              constants.LoanStatus_InRepayment,
			CreatedAt:           now,
		},
	}); err != nil {
		return 0, err
	}

	if err = clients.DBBatchInsertBillingHistories(ctx, txn, billingHistories); err != nil {
		return 0, err
	}

	if err = clients.DBCommitTransaction(txn); err != nil {
		return 0, err
	}
	return loanID, nil
}

func createRepaymentSchedule(loanModel dtos.LoanRequestModel) ([]dtos.BillingModel, []dtos.BillingHistoryModel, error) {
	var (
		billingModels    []dtos.BillingModel
		billingHistories []dtos.BillingHistoryModel
		nextDueTime      *time.Time

		now              = time.Now().UnixMilli()
		tenure           = decimal.NewFromInt(int64(loanModel.TenureValue))
		disbursementDate = time.UnixMilli(loanModel.DisbursementTime)

		principalAmountPerPayment = loanModel.LoanAmount.Div(tenure)
		interestAmountPerPayment  = loanModel.LoanAmount.Mul(loanModel.AnnualInterestRate).Div(constants.Percent).Div(tenure)
		totalAmountPerPayment     = principalAmountPerPayment.Add(interestAmountPerPayment)
	)

	nextDueTime = utils.GetNextTenureSchedule(disbursementDate, loanModel.TenureUnit)
	for recurringIdx := 1; recurringIdx <= loanModel.TenureValue; recurringIdx++ {
		if nextDueTime = utils.GetNextTenureSchedule(*nextDueTime, loanModel.TenureUnit); nextDueTime == nil {
			return nil, nil, fmt.Errorf("unable to get next tenure schedule from %v", *nextDueTime)
		}

		billingID := uuid.NewString()
		billingModels = append(billingModels, dtos.BillingModel{
			BillingID:          billingID,
			LoanID:             loanModel.ID,
			PaymentID:          0,
			RecurringIndex:     recurringIdx,
			PrincipalAmount:    principalAmountPerPayment,
			InterestAmount:     interestAmountPerPayment,
			TotalAmount:        totalAmountPerPayment,
			DueTime:            nextDueTime.UnixMilli(),
			PaymentCompletedAt: 0,
			Status:             constants.PaymentStatus_Pending,
			CreatedAt:          now,
			UpdatedAt:          now,
		})
		billingHistories = append(billingHistories, dtos.BillingHistoryModel{
			BillingID:          billingID,
			PaymentCompletedAt: 0,
			Status:             constants.PaymentStatus_Pending,
			CreatedAt:          now,
		})
	}
	return billingModels, billingHistories, nil
}

func validateLoanRequest(ctx context.Context, param dtos.CreateLoanRequestParam) error {
	if _, err := clients.DBGetUserByID(ctx, param.UserID); err != nil {
		return err
	}

	loanAmount, err := decimal.NewFromString(param.LoanAmount)
	if err != nil {
		return fmt.Errorf("%w. %s", constants.ErrInvalidValue, err.Error())
	}
	if loanAmount.LessThan(decimal.NewFromInt(1)) {
		return fmt.Errorf("%w. outstanding_amount should be greater than 0", constants.ErrInvalidValue)
	}

	if param.TenureValue < 1 {
		return fmt.Errorf("%w. tenure_value should be greater than 0", constants.ErrInvalidValue)
	}

	if !constants.TenureUnit(param.TenureUnit).IsValid() {
		return fmt.Errorf("%w. tenure_unit", constants.ErrInvalidValue)
	}

	annualInterestRate, err := decimal.NewFromString(param.AnnualInterestRate)
	if err != nil {
		return fmt.Errorf("%w. %s", constants.ErrInvalidValue, err.Error())
	}
	if annualInterestRate.LessThan(decimal.NewFromInt(0)) {
		return fmt.Errorf("%w. annual_interest_rate should be greater or equals to 0", constants.ErrInvalidValue)
	}

	return nil
}
