package main

import (
	"github.com/parviz-yu/expense-tracker/internal/config"
	"github.com/parviz-yu/expense-tracker/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)
}
