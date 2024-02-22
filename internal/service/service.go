package service

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/expense-tracker/internal/models"
	"github.com/parviz-yu/expense-tracker/internal/storage"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
	"github.com/parviz-yu/expense-tracker/pkg/types"
)

type ServiceI interface {
	AddExpense(ctx context.Context, expenseReq *models.ExpenseReq) error
	GetCategoriesExpenses(ctx context.Context, filters *models.Filters) ([]*models.CategoryExpensesResp, error)
}

type service struct {
	log  logger.LoggerI
	strg storage.StorageI
}

func NewService(log logger.LoggerI, strg storage.StorageI) ServiceI {
	return &service{
		log:  log,
		strg: strg,
	}
}

func (s *service) AddExpense(ctx context.Context, expenseReq *models.ExpenseReq) error {
	const fn = "service.AddExpense"

	s.log = logger.With(
		s.log,
		logger.String("request_id", middleware.GetReqID(ctx)),
		logger.Any("expenseReq", expenseReq),
	)

	s.log.Info("getting category_id...")

	categoryID, err := s.strg.Category().GetCategoryID(ctx, expenseReq.Category)
	if err != nil {
		s.log.Error("failed to get category_id", logger.Error(err))

		return fmt.Errorf("%s: %w", fn, err)
	}

	expenseInner := &models.ExpenseInner{
		UserID:      expenseReq.UserID,
		CategoryID:  categoryID,
		Amount:      expenseReq.Amount.ToSmallUnit(),
		Description: expenseReq.Description,
		Date:        expenseReq.Date.Time,
	}

	s.log.Info("adding expense...")

	if err := s.strg.Expense().AddExpense(ctx, expenseInner); err != nil {
		s.log.Error("failed to add expense", logger.Error(err))

		return fmt.Errorf("%s: %w", fn, err)
	}

	s.log.Info("expense successfully added")

	return nil
}

func (s *service) GetCategoriesExpenses(ctx context.Context, filters *models.Filters) ([]*models.CategoryExpensesResp, error) {
	const fn = "service.AddExpense"

	s.log = logger.With(
		s.log,
		logger.String("request_id", middleware.GetReqID(ctx)),
		logger.Any("filters", *filters),
	)

	s.log.Info("checking category for existance...")

	if filters.Category != "" {

		_, err := s.strg.Category().GetCategoryID(ctx, filters.Category)
		if err != nil {
			s.log.Error("failed to get category_id", logger.Error(err))

			return nil, fmt.Errorf("%s: %w", fn, err)
		}
	}

	s.log.Info("getting all users stats...")

	usersStats, err := s.strg.Expense().GetAllUsersStats(ctx, filters)
	if err != nil {
		s.log.Error("failed to get users stats", logger.Error(err))

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	categoryExpensesResp := make([]*models.CategoryExpensesResp, 0)
	if len(usersStats) == 0 {
		s.log.Info("empty result with these filters")

		return categoryExpensesResp, nil
	}

	prevCategory := ""
	var result *models.CategoryExpensesResp
	for _, item := range usersStats {

		if item.Category != prevCategory {
			result = &models.CategoryExpensesResp{}
			result.Category = item.Category
			categoryExpensesResp = append(categoryExpensesResp, result)
		}

		if item.Count == 0 {
			result.UsersExpenses = make([]models.UserExpense, 0)
			continue
		}

		currUserExpense := models.UserExpense{
			UserID: item.UserID,
			Total:  types.Money(item.Sum) / 100,
			Count:  item.Count,
		}

		result.UsersExpenses = append(result.UsersExpenses, currUserExpense)
		result.TotalExpenses += types.Money(item.Sum) / 100
		result.ExpensesCount += item.Count

		prevCategory = item.Category
	}

	return categoryExpensesResp, nil
}
