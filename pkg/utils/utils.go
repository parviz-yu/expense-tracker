package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/parviz-yu/expense-tracker/pkg/errs"
	"github.com/parviz-yu/expense-tracker/pkg/types"
)

func ParseString2Time(s string) (time.Time, error) {
	const fn = "pkg.utils.ParseString2Time"

	var date time.Time
	var err error
	if s != "" {
		date, err = time.Parse(types.Layout, s)
		if err != nil {
			return time.Time{}, fmt.Errorf("%s: %w", fn, err)
		}
	}

	return date, nil
}

func ParseString2Float64(s string) (float64, error) {
	const fn = "pkg.utils.ParseString2Float64"

	var val float64
	var err error
	if s != "" {
		val, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", fn, err)
		}
	}

	return val, nil
}

func VerifyTimes(s1 string, s2 string) (time.Time, time.Time, error) {
	const fn = "pkg.utils.VerifyTimes"

	startTime, err := ParseString2Time(s1)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("%s: %w", fn, err)
	}

	endTime, err := ParseString2Time(s2)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("%s: %w", fn, err)
	}

	if !startTime.IsZero() && !endTime.IsZero() && !startTime.Before(endTime) {
		return time.Time{}, time.Time{}, fmt.Errorf("%s: %w", fn, errs.ErrInvalidDateRange)
	}

	return startTime, endTime, nil
}

func VerifyMinMax(min string, max string) (float64, float64, error) {
	const fn = "pkg.utils.VerifyTimes"

	minVal, err := ParseString2Float64(min)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", fn, err)
	}

	maxVal, err := ParseString2Float64(max)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", fn, err)
	}

	if minVal == 0 && maxVal == 0 {
		return minVal, maxVal, nil
	}

	if minVal < 0 || maxVal < 0 || (maxVal < minVal && maxVal != 0) {
		return 0, 0, fmt.Errorf("%s: %w", fn, errs.ErrInvalidMinMaxAmounts)
	}

	return minVal, maxVal, nil
}
