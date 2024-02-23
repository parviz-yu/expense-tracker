package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (r *expensesRepo) GetAllUsersStats(ctx context.Context, filters *models.Filters) ([]*models.UsersStats, error) {
	const fn = "storage.postgres.GetAllUsersStats"

	queryStart := `SELECT c.name, COALESCE(e.user_id, ''), COALESCE(SUM(e.amount), 0), COUNT(e.amount)
	FROM categories c LEFT JOIN expenses e ON c.id=e.category_id`

	queryEnd := `GROUP BY c.name, e.user_id ORDER BY c.name`

	filtersSlice, params := checkFilters(filters)
	whereClause := ""
	if len(filtersSlice) != 0 {
		whereClause = fmt.Sprintf("WHERE %s", strings.Join(filtersSlice, " AND "))
		whereClause = putPlaceholders(whereClause, len(filtersSlice))
	}
	query := fmt.Sprintf("%s %s %s", queryStart, whereClause, queryEnd)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	stats := make([]*models.UsersStats, 0)
	for rows.Next() {
		stat := &models.UsersStats{}
		err := rows.Scan(&stat.Category, &stat.UserID, &stat.Sum, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return stats, nil
}

func (r *expensesRepo) GetUserStats(ctx context.Context, userID string, filters *models.Filters) ([]*models.UserStats, error) {
	const fn = "storage.postgres.GetAllUsersStats"

	queryStart := `
	SELECT c.name, e.date, e.description, e.amount,
    	COALESCE(SUM(e.amount) OVER w, 0) AS total_sum,
    	COALESCE(COUNT(e.amount) OVER w, 0) AS total_count
	FROM categories c LEFT JOIN expenses e ON c.id=e.category_id`

	queryEnd := "WINDOW w AS (PARTITION BY c.name)"

	filtersSlice, params := checkFilters(filters)
	filtersSlice = append(filtersSlice, "e.user_id=?")
	params = append(params, userID)

	whereClause := fmt.Sprintf("WHERE %s", strings.Join(filtersSlice, " AND "))
	whereClause = putPlaceholders(whereClause, len(filtersSlice))

	query := fmt.Sprintf("%s %s %s", queryStart, whereClause, queryEnd)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	stats := make([]*models.UserStats, 0)
	for rows.Next() {
		stat := &models.UserStats{}
		err := rows.Scan(&stat.Category, &stat.Date, &stat.Description, &stat.Amount, &stat.Total, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return stats, nil
}

func checkFilters(filters *models.Filters) ([]string, []interface{}) {
	filtersSlice := make([]string, 0)
	params := make([]interface{}, 0)
	if filters.Category != "" {
		filtersSlice = append(filtersSlice, "c.name = ?")
		params = append(params, filters.Category)
	}

	if !filters.StartDate.IsZero() {
		filtersSlice = append(filtersSlice, "e.date >= ?")
		params = append(params, filters.StartDate)
	}

	if !filters.EndDate.IsZero() {
		filtersSlice = append(filtersSlice, "e.date <= ?")
		params = append(params, filters.EndDate)
	}

	if filters.MinAmount > 0 {
		filtersSlice = append(filtersSlice, "e.amount >= ?")
		params = append(params, filters.MinAmount)
	}

	if filters.MaxAmount != 0 {
		filtersSlice = append(filtersSlice, "e.amount <= ?")
		params = append(params, filters.MaxAmount)
	}

	return filtersSlice, params
}

func putPlaceholders(line string, n int) string {
	for i := 1; i <= n; i++ {
		s := fmt.Sprintf("$%d", i)
		line = strings.Replace(line, "?", s, 1)
	}

	return line
}
