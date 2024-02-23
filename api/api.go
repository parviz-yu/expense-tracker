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
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/tracker", func(r chi.Router) {
		r.Post("/", h.AddExpense)
		r.Post("/categories", h.AddNewCategory)
		r.Get("/stats/categories", h.UsersCategoriesStats)
		r.Get("/stats/users/{id}", h.UserStats)
	})

	return r
}
