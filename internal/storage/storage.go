package storage

import (
	"context"

	"github.com/parviz-yu/expense-tracker/internal/models"
)

type StorageI interface {
	CloseDB()
	Category() CategoryRepoI
	Expense() ExpenseRepoI
}

type CategoryRepoI interface {
	GetCategoryID(ctx context.Context, name string) (int, error)
}

type ExpenseRepoI interface {
	AddExpense(ctx context.Context, expense *models.ExpenseInner) error
}
