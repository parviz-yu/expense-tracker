package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/expense-tracker/internal/models"
	"github.com/parviz-yu/expense-tracker/internal/storage"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
	"github.com/parviz-yu/expense-tracker/pkg/types"
)

type ServiceI interface {
	AddExpense(ctx context.Context, expenseReq *models.ExpenseReq) error
	GetCategoriesExpenses(ctx context.Context, filters *models.Filters) ([]*models.CategoryExpensesResp, error)
	GetUserExpenses(ctx context.Context, userID string, filters *models.Filters, verbose bool) ([]*models.UserExpensesResp, error)
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

	log := logger.With(
		s.log,
		logger.String("request_id", middleware.GetReqID(ctx)),
		logger.Any("expenseReq", expenseReq),
	)

	log.Info("getting category_id...")

	categoryID, err := s.strg.Category().GetCategoryID(ctx, expenseReq.Category)
	if err != nil {
		log.Error("failed to get category_id", logger.Error(err))

		return fmt.Errorf("%s: %w", fn, err)
	}

	expenseInner := &models.ExpenseInner{
		UserID:      expenseReq.UserID,
		CategoryID:  categoryID,
		Amount:      expenseReq.Amount.ToSmallUnit(),
		Description: expenseReq.Description,
		Date:        time.Time(expenseReq.Date),
	}

	log.Info("adding expense...")

	if err := s.strg.Expense().AddExpense(ctx, expenseInner); err != nil {
		log.Error("failed to add expense", logger.Error(err))

		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("expense successfully added")

	return nil
}

func (s *service) GetCategoriesExpenses(ctx context.Context, filters *models.Filters) ([]*models.CategoryExpensesResp, error) {
	const fn = "service.AddExpense"

	log := logger.With(
		s.log,
		logger.String("request_id", middleware.GetReqID(ctx)),
		logger.Any("filters", *filters),
	)

	log.Info("checking category for existance...")

	if filters.Category != "" {

		_, err := s.strg.Category().GetCategoryID(ctx, filters.Category)
		if err != nil {
			log.Error("failed to get category_id", logger.Error(err))

			return nil, fmt.Errorf("%s: %w", fn, err)
		}
	}

	log.Info("getting all users stats...")

	usersStats, err := s.strg.Expense().GetAllUsersStats(ctx, filters)
	if err != nil {
		log.Error("failed to get users stats", logger.Error(err))

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	categoryExpensesResp := make([]*models.CategoryExpensesResp, 0)
	if len(usersStats) == 0 {
		log.Info("empty result with these filters")

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

func (s *service) GetUserExpenses(ctx context.Context, userID string, filters *models.Filters, verbose bool) ([]*models.UserExpensesResp, error) {
	const fn = "service.GetUserExpenses"

	log := logger.With(
		s.log,
		logger.String("request_id", middleware.GetReqID(ctx)),
		logger.Any("filters", *filters),
	)

	log.Info("checking category for existance...")

	if filters.Category != "" {

		_, err := s.strg.Category().GetCategoryID(ctx, filters.Category)
		if err != nil {
			log.Error("failed to get category_id", logger.Error(err))

			return nil, fmt.Errorf("%s: %w", fn, err)
		}
	}

	log.Info("getting user's expenses...")

	userStats, err := s.strg.Expense().GetUserStats(ctx, userID, filters)
	if err != nil {
		log.Error("failed to get user's stats", logger.Error(err))

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	userExpensesResp := make([]*models.UserExpensesResp, 0)
	if len(userStats) == 0 {
		log.Info("empty result with these filters")

		return userExpensesResp, nil
	}

	prevCategory := ""
	var result *models.UserExpensesResp
	for _, item := range userStats {

		if item.Category != prevCategory {
			result = &models.UserExpensesResp{}
			result.Category = item.Category
			userExpensesResp = append(userExpensesResp, result)
		}

		if item.Count == 0 {
			result.UserExpenses = make([]models.UserExpenseExtended, 0)
			continue
		}

		userExpense := models.UserExpenseExtended{
			Description: item.Description,
			Amount:      types.Money(item.Amount) / 100,
			Date:        types.CustomTime(item.Date),
		}

		result.UserExpenses = append(result.UserExpenses, userExpense)
		result.TotalExpenses += types.Money(item.Amount) / 100
		result.ExpensesCount += item.Count

		prevCategory = item.Category
	}

	return userExpensesResp, nil
}
