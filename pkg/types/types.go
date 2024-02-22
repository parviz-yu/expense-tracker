package types

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Money float64

func (m Money) MarshalJSON() ([]byte, error) {
	toMarshal := strconv.FormatFloat(float64(m), 'f', 2, 64)

	return []byte(toMarshal), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var res float64

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	*m = Money(res)

	return nil
}

func (m Money) ToSmallUnit() int {
	return int(m * 100)
}

func (m Money) IsPositive() bool {
	return m >= 0
}

type CustomTime struct {
	time.Time
}

const Layout = "02-01-2006"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}

	ct.Time, err = time.Parse(Layout, s)
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Format(Layout)), nil
}
