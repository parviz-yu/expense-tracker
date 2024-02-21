package service

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/expense-tracker/internal/models"
	"github.com/parviz-yu/expense-tracker/internal/storage"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
)

type ServiceI interface {
	AddExpense(ctx context.Context, expenseReq *models.ExpenseReq) error
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
