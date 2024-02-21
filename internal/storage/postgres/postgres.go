package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/parviz-yu/expense-tracker/internal/config"
	"github.com/parviz-yu/expense-tracker/internal/storage"
)

type Storage struct {
	db             *sql.DB
	categoriesRepo *categoriesRepo
	expensesRepo   *expensesRepo
}

func NewStorage(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	const fn = "storage.postgres.NewStorage"

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.PostgresUser,
		cfg.Database.PostgresPassword,
		cfg.Database.PostgresHost,
		cfg.Database.PostgresPort,
		cfg.Database.PostgresDatabase,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return newStorage(db), nil
}

func newStorage(db *sql.DB) *Storage {
	return &Storage{
		db:             db,
		categoriesRepo: newCategoriesRepo(db),
		expensesRepo:   newExpensesRepo(db),
	}
}

func (s *Storage) CloseDB() {
	s.db.Close()
}

func (s *Storage) Category() storage.CategoryRepoI {
	return s.categoriesRepo
}

func (s *Storage) Expense() storage.ExpenseRepoI {
	return s.expensesRepo
}
