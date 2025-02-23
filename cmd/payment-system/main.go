package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"infotecsTest/internal/config"
	"infotecsTest/internal/http-server/handlers/transaction"
	"infotecsTest/internal/http-server/handlers/wallet"
	mwLogger "infotecsTest/internal/http-server/middleware/logger"
	"infotecsTest/internal/lib/logger/sl"
	"infotecsTest/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to initialize storage", sl.Err(err))
	}
	defer func() {
		if err = storage.Close(); err != nil {
			logger.Error("unable to close database: ", sl.Err(err))
		}
	}() //closing storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)

	router.Get("/api/transactions", transaction.GetLast(logger, storage))
	router.Get("/api/wallet/{address}/balance", wallet.GetBalance(logger, storage))
	router.Post("/api/send", transaction.Send(logger, storage))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil {
			logger.Error("failed to serve server")
		}
	}()

	logger.Info("server started")

	<-done
	logger.Info("stopping server")
}

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
