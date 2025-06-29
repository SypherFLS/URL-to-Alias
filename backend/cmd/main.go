package main

import (
    "log/slog"
    // "fmt"  
    "os"  
	"project/internal/config"
)

const (
    envLocal = "local"
    envDev = "dev"
    envProd = "prod"
)

func main() {   

	cfg := config.MustLoad()

    log := setupLogger(cfg.Env)

    log.Info("starting backend", slog.String("env", cfg.Env))
    // TODO: init logger : slog
    // TODO: init db : gorm
    // TODO: init router : chi, chirender
    // TODO: init storage : sqlIte
}


func setupLogger (env string) *slog.Logger  {
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