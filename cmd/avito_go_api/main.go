package main

import (
	"avito_go_api/cmd/internal/config"
	"avito_go_api/cmd/internal/http-server/handlers/add_segment"
	"avito_go_api/cmd/internal/http-server/handlers/delete_segment"
	"avito_go_api/cmd/internal/http-server/handlers/get_segments"
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
	//fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("Starting...", slog.String("env", cfg.Env))
	log.Debug("DEBUGGING@!!!")

	storage, err := postgresql.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	//err = storage.CreateSegment("SEGMENT_32722")
	//if err != nil {
	//	log.Error("failed to add segment", sl.Err(err))
	//	os.Exit(1)
	//
	//}
	//
	err = storage.DeleteSegment("SEGMENT_27")
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage
	//
	//err = storage.DeleteSegment("SEGMENT_1")
	//if err != nil {
	//	log.Error("failed to delete segment", sl.Err(err))
	//	os.Exit(1)
	//
	//}
	//err = storage.ReassignSegments([]string{"SEGMENT_2", "SEGMENT_272", "SEGMENT_2722"}, []string{}, 1)
	//if err != nil {
	//	log.Error("failed to init storage", sl.Err(err))
	//	os.Exit(1)
	//}
	//segments, err := storage.GetSegments(1)
	//if err != nil {
	//	log.Error("failed to get segments", sl.Err(err))
	//	os.Exit(1)
	//}
	//fmt.Print(segments)
	// TODO: init router: chi, chi render

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/url", add_segment.New(log, storage))
	router.Delete("/url", delete_segment.New(log, storage))
	router.Put("/url", reassign_segments.New(log, storage))
	router.Get("/url", get_segments.New(log, storage))
	// TODO: run server
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
