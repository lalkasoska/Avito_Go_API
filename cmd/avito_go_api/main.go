package main

import (
	"avito_go_api/cmd/internal/config"
	"avito_go_api/cmd/internal/http-server/handlers/add_segment"
	"avito_go_api/cmd/internal/http-server/handlers/delete_segment"
	"avito_go_api/cmd/internal/http-server/handlers/get_segments"
	"avito_go_api/cmd/internal/http-server/handlers/get_user_history"
	"avito_go_api/cmd/internal/http-server/handlers/reassign_segments"
	mwLogger "avito_go_api/cmd/internal/http-server/middleware/logger"
	"avito_go_api/cmd/internal/lib/logger/sl"
	"avito_go_api/cmd/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Starting...", slog.String("env", cfg.Env))

	storage, err := postgresql.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/", add_segment.New(log, storage))
	router.Delete("/", delete_segment.New(log, storage))
	router.Put("/", reassign_segments.New(log, storage))
	router.Get("/", get_segments.New(log, storage))
	router.Get("/get_history", get_user_history.New(log, storage, *cfg))
	router.Get("/report", func(w http.ResponseWriter, r *http.Request) {
		filePath := "latest_segment_history_report.csv"
		http.ServeFile(w, r, filePath)

	})

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log

}
