package models

import (
	"time"

	"github.com/parviz-yu/expense-tracker/pkg/types"
)

type ExpenseReq struct {
	UserID      string           `json:"user_id"`
	Category    string           `json:"category"`
	Description string           `json:"description,omitempty"`
	Amount      types.Money      `json:"amount"`
	Date        types.CustomTime `json:"date"`
}

type ExpenseInner struct {
	Date        time.Time
	UserID      string
	Description string
	CategoryID  int
	Amount      int
}

type UsersStats struct {
	Category string
	UserID   string
	Sum      int
	Count    int
}

type UserExpense struct {
	UserID string      `json:"user_id"`
	Total  types.Money `json:"total"`
	Count  int         `json:"count"`
}

type CategoryExpensesResp struct {
	Category      string        `json:"category"`
	TotalExpenses types.Money   `json:"total_expenses"`
	ExpensesCount int           `json:"expenses_count"`
	UsersExpenses []UserExpense `json:"users_expenses"`
}

type Filters struct {
	Category  string
	MinAmount int
	MaxAmount int
	StartDate time.Time
	EndDate   time.Time
}

type UserStats struct {
	Category    string
	Description string
	Amount      int
	Total       int
	Count       int
	Date        time.Time
}

type UserExpenseExtended struct {
	Amount      types.Money      `json:"amount"`
	Description string           `json:"description"`
	Date        types.CustomTime `json:"date"`
}
type UserExpensesResp struct {
	Category      string                `json:"category"`
	TotalExpenses types.Money           `json:"total_expenses"`
	ExpensesCount int                   `json:"expenses_count"`
	UserExpenses  []UserExpenseExtended `json:"expenses"`
}
