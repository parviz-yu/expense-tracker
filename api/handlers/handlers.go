package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/expense-tracker/internal/models"
	"github.com/parviz-yu/expense-tracker/internal/service"
	"github.com/parviz-yu/expense-tracker/pkg/errs"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
	"github.com/parviz-yu/expense-tracker/pkg/utils"
)

type Handler struct {
	log logger.LoggerI
	svc service.ServiceI
}

func NewHandler(log logger.LoggerI, svc service.ServiceI) *Handler {
	return &Handler{
		log: log,
		svc: svc,
	}
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.AddExpense"

	log := logger.With(
		h.log,
		logger.String("fn", fn),
		logger.String("request_id", middleware.GetReqID(r.Context())),
	)

	req := models.ExpenseReq{}
	jsonDec := json.NewDecoder(r.Body)
	jsonDec.DisallowUnknownFields()
	if err := jsonDec.Decode(&req); err != nil {
		log.Error(err.Error())

		Error(w, r, http.StatusBadRequest, errs.ErrInvalidRequestBody)
		return
	}

	if !req.Amount.IsPositive() {
		log.Warn(errs.ErrNegativeAmount.Error(), logger.Any("req", req))

		Error(w, r, http.StatusBadRequest, errs.ErrNegativeAmount)
		return
	}

	err := h.svc.AddExpense(r.Context(), &req)
	if errors.Is(err, errs.ErrCategoryNotExists) {
		log.Warn(err.Error(), logger.Any("req", req))

		Error(w, r, http.StatusNotFound, errs.ErrCategoryNotExists)
		return
	}
	if err != nil {
		log.Error(err.Error(), logger.Any("req", req))

		Error(w, r, http.StatusInternalServerError, errs.ErrInternalServerError)
		return
	}

	Respond(w, r, http.StatusOK, nil)
}

func (h *Handler) UsersCategoriesStats(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.GetAllUsersStats"

	log := logger.With(
		h.log,
		logger.String("fn", fn),
		logger.String("request_id", middleware.GetReqID(r.Context())),
	)

	category := URLQueryParam(r, "category")
	start := URLQueryParam(r, "start_date")
	end := URLQueryParam(r, "end_date")
	min := URLQueryParam(r, "min_amount")
	max := URLQueryParam(r, "max_amount")

	startDate, endDate, err := utils.VerifyTimes(start, end)
	if err != nil {
		log.Error(err.Error(), logger.String("start_date", start), logger.String("end_date", end))

		if errors.Is(err, errs.ErrInvalidDateRange) {
			Error(w, r, http.StatusBadRequest, errs.ErrInvalidDateRange)
			return
		}

		Error(w, r, http.StatusBadRequest, errs.ErrInvalidDateFormat)
		return
	}

	minAmount, maxAmount, err := utils.VerifyMinMax(min, max)
	if err != nil {
		log.Error(err.Error(), logger.String("min_amount", min), logger.String("max_amount", max))

		Error(w, r, http.StatusBadRequest, errs.ErrInvalidMinMaxAmounts)
		return
	}

	filters := &models.Filters{
		Category:  category,
		MinAmount: int(minAmount * 100),
		MaxAmount: int(maxAmount * 100),
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := h.svc.GetCategoriesExpenses(r.Context(), filters)
	if errors.Is(err, errs.ErrCategoryNotExists) {
		log.Warn(err.Error(), logger.Any("filter", filters))

		Error(w, r, http.StatusNotFound, errs.ErrCategoryNotExists)
		return
	}
	if err != nil {
		log.Error(err.Error(), logger.Any("filter", filters))

		Error(w, r, http.StatusInternalServerError, errs.ErrInternalServerError)
		return
	}

	Respond(w, r, http.StatusOK, map[string]interface{}{"result": resp})
}

func (h *Handler) UserStats(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.GetAllUsersStats"

	log := logger.With(
		h.log,
		logger.String("fn", fn),
		logger.String("request_id", middleware.GetReqID(r.Context())),
	)

	category := URLQueryParam(r, "category")
	start := URLQueryParam(r, "start_date")
	end := URLQueryParam(r, "end_date")
	min := URLQueryParam(r, "min_amount")
	max := URLQueryParam(r, "max_amount")
	userID := strings.TrimSpace(chi.URLParam(r, "id"))

	if userID == "" {
		log.Error(errs.ErrEmptyUserIDParam.Error())

		Error(w, r, http.StatusBadRequest, errs.ErrEmptyUserIDParam)
		return
	}

	startDate, endDate, err := utils.VerifyTimes(start, end)
	if err != nil {
		log.Error(err.Error(), logger.String("start_date", start), logger.String("end_date", end))

		if errors.Is(err, errs.ErrInvalidDateRange) {
			Error(w, r, http.StatusBadRequest, errs.ErrInvalidDateRange)
			return
		}

		Error(w, r, http.StatusBadRequest, errs.ErrInvalidDateFormat)
		return
	}

	minAmount, maxAmount, err := utils.VerifyMinMax(min, max)
	if err != nil {
		log.Error(err.Error(), logger.String("min_amount", min), logger.String("max_amount", max))

		Error(w, r, http.StatusBadRequest, errs.ErrInvalidMinMaxAmounts)
		return
	}

	filters := &models.Filters{
		Category:  category,
		MinAmount: int(minAmount * 100),
		MaxAmount: int(maxAmount * 100),
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := h.svc.GetUserExpenses(r.Context(), userID, filters)
	if errors.Is(err, errs.ErrCategoryNotExists) {
		log.Warn(err.Error(), logger.Any("filter", filters))

		Error(w, r, http.StatusNotFound, errs.ErrCategoryNotExists)
		return
	}
	if err != nil {
		log.Error(err.Error(), logger.Any("filter", filters))

		Error(w, r, http.StatusInternalServerError, errs.ErrInternalServerError)
		return
	}

	Respond(w, r, http.StatusOK, map[string]interface{}{"result": resp})
}
