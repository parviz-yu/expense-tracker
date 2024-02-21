package types

import (
	"encoding/json"
	"strconv"
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
