package constants

import "github.com/shopspring/decimal"

type TenureUnit int8

const (
	TenureUnit_Day TenureUnit = iota + 1
	TenureUnit_Week
	TenureUnit_Month
	TenureUnit_Year
)

func (c TenureUnit) IsValid() bool {
	for i := TenureUnit_Day; i <= TenureUnit_Year; i++ {
		if i == c {
			return true
		}
	}
	return false
}

type LoanStatus int8

const (
	LoanStatus_InRepayment LoanStatus = iota + 1
	LoanStatus_Defaulted
	LoanStatus_Completed
)

type PaymentStatus int8

const (
	PaymentStatus_Pending PaymentStatus = iota + 1
	PaymentStatus_Completed
)

var (
	Percent = decimal.NewFromInt(100)
)
