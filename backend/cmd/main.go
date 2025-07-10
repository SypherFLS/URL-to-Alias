package main

import (
	"log/slog"
	_ "net/http"
	// "fmt"
	"os"
	"project/internal/config"
	"project/internal/lib/logger/sl"
	"project/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting backend", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init db", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()
	
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	

	// TODO: init logger : slog
	// TODO: init db : gorm
	// TODO: init router : chi, chirender
	// TODO: init storage : sqlIte

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
