package services

import (
	"context"
	"fmt"
	"time"

	"loan-payment/clients"
	"loan-payment/constants"
	"loan-payment/dtos"
	"loan-payment/utils"

	"github.com/shopspring/decimal"
)

func MakePayment(ctx context.Context, param dtos.MakePaymentParam) error {
	if err := validatePayment(param); err != nil {
		return err
	}

	txn, err := clients.DBBeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer clients.DBRollbackTransaction(txn)

	// prevent update racing with pessimistic lock
	loanRequestModel, err := clients.DBGetLoanRequestByIDAndUserIDForUpdate(ctx, txn, param.LoanID, param.UserID)
	if err != nil {
		return err
	}

	nearestBillingSchedule, err := getNextNearestBillingSchedule(time.UnixMilli(loanRequestModel.DisbursementTime), loanRequestModel.TenureUnit)
	if err != nil {
		return err
	}

	// prevent update racing with pessimistic lock
	pendingBillings, err := clients.DBGetPendingBillingsWithDueTimeByLoanIdForUpdate(ctx, txn, param.LoanID, nearestBillingSchedule.UnixMilli())
	if err != nil {
		return err
	}

	var (
		now = time.Now().UnixMilli()

		includeLastBilling bool
		billingIDs         []int64
		outstandingAmount  decimal.Decimal
		principalAmount    decimal.Decimal
		interestAmount     decimal.Decimal
		billingHistories   []dtos.BillingHistoryModel
	)
	for _, billing := range pendingBillings {
		outstandingAmount = outstandingAmount.Add(billing.TotalAmount)
		principalAmount = principalAmount.Add(billing.PrincipalAmount)
		interestAmount = interestAmount.Add(billing.InterestAmount)

		billingIDs = append(billingIDs, billing.ID)
		billingHistories = append(billingHistories, dtos.BillingHistoryModel{
			BillingID:          billing.BillingID,
			PaymentCompletedAt: now,
			Status:             constants.PaymentStatus_Completed,
			CreatedAt:          now,
		})

		if billing.RecurringIndex == loanRequestModel.TenureValue {
			includeLastBilling = true
		}
	}

	paymentAmount, _ := decimal.NewFromString(param.Amount)
	if !paymentAmount.Equal(outstandingAmount) {
		return fmt.Errorf("%w. incorrect payment amount", constants.ErrInvalidValue)
	}

	paymentID, err := clients.DBInsertPayment(ctx, txn, &dtos.PaymentModel{
		UserID: param.UserID,
		Amount: paymentAmount,
	})
	if err != nil {
		return err
	}

	if err = clients.DBBulkPayBillingsByIDs(ctx, txn, billingIDs, paymentID); err != nil {
		return err
	}

	loanRequestStatus := loanRequestModel.Status
	if includeLastBilling {
		loanRequestStatus = constants.LoanStatus_Completed
	}
	if err = clients.DBUpdateLoanRequestPaymentByID(ctx, txn, loanRequestModel.ID, principalAmount, interestAmount, loanRequestStatus); err != nil {
		return err
	}

	if err = clients.DBBatchInsertLoanRequestHistories(ctx, txn, []dtos.LoanRequestHistory{
		{
			LoanID:              loanRequestModel.ID,
			PrincipalPaidAmount: principalAmount,
			InterestPaidAmount:  interestAmount,
			Status:              loanRequestStatus,
			CreatedAt:           now,
		},
	}); err != nil {
		return err
	}

	if err = clients.DBBatchInsertBillingHistories(ctx, txn, billingHistories); err != nil {
		return err
	}
	return clients.DBCommitTransaction(txn)
}

func validatePayment(param dtos.MakePaymentParam) error {
	paymentAmount, err := decimal.NewFromString(param.Amount)
	if err != nil {
		return fmt.Errorf("%w. unable to parse payment's amount", constants.ErrInvalidValue)
	} else if paymentAmount.LessThan(decimal.Zero) {
		return fmt.Errorf("%w. payment's amount should be greater than 0", constants.ErrInvalidValue)
	}
	return nil
}

func getNextNearestBillingSchedule(disbursementTime time.Time, tenureUnit constants.TenureUnit) (*time.Time, error) {
	var (
		now                 = time.Now()
		nextNearestSchedule = utils.GetNextTenureSchedule(
			time.Date(
				disbursementTime.Year(),
				disbursementTime.Month(),
				disbursementTime.Day(),
				23,
				59,
				59,
				0,
				disbursementTime.Location(),
			),
			tenureUnit,
		)
	)

	if nextNearestSchedule == nil {
		return nil, fmt.Errorf("unable to get next repayment schedule")
	}

	for nextNearestSchedule.Before(now) {
		nextNearestSchedule = utils.GetNextTenureSchedule(*nextNearestSchedule, tenureUnit)
	}
	return nextNearestSchedule, nil
}
