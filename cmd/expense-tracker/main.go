package main

import (
	"context"
	"os"

	"github.com/parviz-yu/expense-tracker/internal/config"
	"github.com/parviz-yu/expense-tracker/internal/storage/postgres"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)
	log.Info("INFO MESSAGE!!!")

	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		log.Error("failed to init storage", logger.Error(err))
		os.Exit(1)
	}

}
