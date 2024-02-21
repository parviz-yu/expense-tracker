package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/parviz-yu/expense-tracker/internal/models"
)

type expensesRepo struct {
	db *sql.DB
}

func newExpensesRepo(db *sql.DB) *expensesRepo {
	return &expensesRepo{
		db: db,
	}
}

func (r *expensesRepo) AddExpense(ctx context.Context, expense *models.ExpenseInner) error {
	const fn = "storage.postgres.AddExpense"

	query := `INSERT INTO expenses (user_id, category_id, amount, description, date) VALUES ($1, $2, $3, $4, $5)`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.ExecContext(
		ctx,
		expense.UserID,
		expense.CategoryID,
		expense.Amount,
		expense.Description,
		expense.Date,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
