package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/expense-tracker/internal/models"
	"github.com/parviz-yu/expense-tracker/internal/service"
	"github.com/parviz-yu/expense-tracker/pkg/errs"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
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
