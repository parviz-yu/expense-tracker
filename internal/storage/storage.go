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
	GetAllUsersStats(ctx context.Context, filters *models.Filters) ([]*models.UsersStats, error)
	GetUserStats(ctx context.Context, userID string, filters *models.Filters) ([]*models.UserStats, error)
}
