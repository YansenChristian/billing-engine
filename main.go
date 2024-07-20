package main

import (
	"context"

	"loan-payment/configs"
	"loan-payment/dtos"
	"loan-payment/services"

	"github.com/sirupsen/logrus"
)

func main() {
	configs.Init("webservice")

	var (
		ctx    = context.Background()
		userID = int64(1)
	)

	loanID, err := services.CreateLoanRequest(ctx, dtos.CreateLoanRequestParam{
		UserID:             userID,
		LoanAmount:         "10000000",
		TenureValue:        50,
		TenureUnit:         2,
		AnnualInterestRate: "10.75",
	})
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("loan_id: %d", loanID)

	outstandingAmount, err := services.GetOutstanding(ctx, dtos.GetOutstandingParam{
		UserID: userID,
		LoanID: loanID,
	})
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("outstanding amount: %s", outstandingAmount.String())

	isDelinquent, err := services.IsDelinquent(ctx, dtos.IsDelinquentParam{LoanID: loanID})
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("is delinquent: %t", isDelinquent)

	if err = services.MakePayment(ctx, dtos.MakePaymentParam{
		UserID: userID,
		LoanID: loanID,
		Amount: "221500",
	}); err != nil {
		logrus.Error(err)
	}
}
