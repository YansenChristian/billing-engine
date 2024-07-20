package services

import (
	"context"

	"loan-payment/clients"
	"loan-payment/constants"
	"loan-payment/dtos"

	"github.com/shopspring/decimal"
)

func GetOutstanding(ctx context.Context, param dtos.GetOutstandingParam) (decimal.Decimal, error) {
	if _, err := clients.DBGetUserByID(ctx, param.UserID); err != nil {
		return decimal.Zero, err
	}

	loanRequestModel, err := clients.DBGetLoanRequestByID(ctx, param.LoanID)
	if err != nil {
		return decimal.Zero, err
	}

	if loanRequestModel.Status == constants.LoanStatus_Completed {
		return decimal.Zero, nil
	}

	outstandingPrincipal := loanRequestModel.LoanAmount.Sub(loanRequestModel.PrincipalPaidAmount)
	outstandingInterest := outstandingPrincipal.Mul(loanRequestModel.AnnualInterestRate).Div(constants.Percent)
	return outstandingPrincipal.Add(outstandingInterest), nil
}
