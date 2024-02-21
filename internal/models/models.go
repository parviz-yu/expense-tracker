package models

import (
	"time"

	"github.com/parviz-yu/expense-tracker/pkg/types"
)

type ExpenseReq struct {
	UserID      string           `json:"user_id"`
	Category    string           `json:"category"`
	Amount      types.Money      `json:"amount"`
	Description string           `json:"description,omitempty"`
	Date        types.CustomTime `json:"date"`
}

type ExpenseInner struct {
	UserID      string
	CategoryID  int
	Amount      int
	Description string
	Date        time.Time
}
