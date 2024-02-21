package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/expense-tracker/api/handlers"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
)

func SetUpRouter(h *handlers.Handler, log logger.LoggerI) *chi.Mux {
	r := chi.NewRouter()

	r.Use(handlers.NewMWLogger(log))
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	r.Route("/api/v1/tracker", func(r chi.Router) {
		r.Post("/", h.AddExpense)
	})

	return r
}
