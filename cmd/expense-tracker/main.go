package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/parviz-yu/expense-tracker/api"
	"github.com/parviz-yu/expense-tracker/api/handlers"
	"github.com/parviz-yu/expense-tracker/internal/config"
	"github.com/parviz-yu/expense-tracker/internal/service"
	"github.com/parviz-yu/expense-tracker/internal/storage/postgres"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)

	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		log.Error("failed to init storage", logger.Error(err))
		os.Exit(1)
	}
	defer strg.CloseDB()

	svc := service.NewService(log, strg)
	hand := handlers.NewHandler(log, svc)
	router := api.SetUpRouter(hand, log)

	listenAddr := fmt.Sprintf(":%s", cfg.HTTPServer.Port)
	log.Info("starting server...", logger.String("address", listenAddr))
	srv := http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", logger.Error(err))
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", logger.Error(err))
		return
	}

	log.Info("server stopped")
}
