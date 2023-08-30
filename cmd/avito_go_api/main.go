package main

import (
	"avito_go_api/cmd/internal/config"
	"avito_go_api/cmd/internal/lib/logger/sl"
	"avito_go_api/cmd/internal/storage/postgresql"
	"log/slog"
	"os"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()
	//fmt.Println(cfg)

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("Starting...", slog.String("env", cfg.Env))
	log.Debug("DEBUGGING@!!!")

	// TODO: init storage: postgresql
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

	// TODO: run server
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
