package main

import (
	"log/slog"
	"net/http"

	// "fmt"
	"os"
	"project/internal/config"
	"project/internal/http-server/handlers/url/out"
	"project/internal/http-server/handlers/url/save"
	"project/internal/http-server/middleware/logger"
	"project/internal/lib/logger/handlers/slogpretty"
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
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")
	log.Warn("warn messages are enabled")
	log.Info("info messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init db", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/", out.TakeAllUrls(log, storage))
	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.HTTP.Address))

	if err := http.ListenAndServe(cfg.HTTP.Address, router); err != nil {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
		// log = slog.New(
		// 	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// )
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)

}
