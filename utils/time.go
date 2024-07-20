package utils

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"

	"loan-payment/constants"
)

func GetNextTenureSchedule(startTime time.Time, timeUnit constants.TenureUnit) *time.Time {
	var (
		expr     *cronexpr.Expression
		nextTime time.Time
	)

	switch timeUnit {
	case constants.TenureUnit_Day:
		expr = cronexpr.MustParse("59 23 * * *")
	case constants.TenureUnit_Week:
		expr = cronexpr.MustParse(fmt.Sprintf("59 23 * * %d", startTime.Weekday()))
	case constants.TenureUnit_Month:
		expr = cronexpr.MustParse(fmt.Sprintf("59 23 %d * *", startTime.Day()))
	case constants.TenureUnit_Year:
		expr = cronexpr.MustParse(fmt.Sprintf("59 23 %d %d *", startTime.Day(), startTime.Month()))
	}

	if expr != nil {
		nextTime = expr.Next(startTime)
		return &nextTime
	}
	return nil
}
