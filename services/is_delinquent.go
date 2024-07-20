package services

import (
	"context"

	"loan-payment/clients"
	"loan-payment/dtos"
)

func IsDelinquent(ctx context.Context, param dtos.IsDelinquentParam) (bool, error) {
	overdueBillings, err := clients.DBGetOverdueBillings(ctx, param.LoanID)
	if err != nil {
		return false, err
	}

	if len(overdueBillings) > 2 {
		return true, nil
	}
	return false, nil
}
