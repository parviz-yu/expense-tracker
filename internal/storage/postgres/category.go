package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/parviz-yu/expense-tracker/pkg/errs"
)

type categoriesRepo struct {
	db *sql.DB
}

func newCategoriesRepo(db *sql.DB) *categoriesRepo {
	return &categoriesRepo{
		db: db,
	}
}

func (r *categoriesRepo) GetCategoryID(ctx context.Context, name string) (int, error) {
	const fn = "storage.postgres.NewStorage"

	query := `SELECT id FROM categories WHERE name = $1`

	var id int
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	err = stmt.QueryRowContext(ctx, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", fn, errs.ErrCategoryNotExists)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (r *categoriesRepo) AddNewCategory(ctx context.Context, name string) error {
	const fn = "storage.postgres.AddNewCategory"

	query := "INSERT INTO categories (name) VALUES ($1)"

	_, err := r.db.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
