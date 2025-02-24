// Package main является точкой входа приложения.
// Содержит конфигурацию, инициализацию компонентов и запуск сервера.
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

// Константы окружений для настройки логгера
const (
	envLocal = "local" // Локальное окружение
	envDev   = "dev"   // Разработочное окружение
	envProd  = "prod"  // Продакшн окружение
)

// main инициализирует и запускает приложение:
// 1. Загружает конфигурацию
// 2. Настраивает логгер
// 3. Инициализирует хранилище
// 4. Настраивает роутер и обработчики
// 5. Запускает HTTP-сервер
// 6. Обрабатывает сигналы завершения
func main() {
	// Загрузка конфигурации приложения
	cfg := config.MustLoad()

	// Инициализация логгера в зависимости от окружения
	logger := setupLogger(cfg.Env)

	// Подключение к хранилищу SQLite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to initialize storage", sl.Err(err))
	}
	defer func() {
		// Гарантированное закрытие соединения с БД при завершении
		if err = storage.Close(); err != nil {
			logger.Error("unable to close database: ", sl.Err(err))
		}
	}()

	// Настройка роутера
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // Добавляет ID к каждому запросу
	router.Use(mwLogger.New(logger)) // Логирование запросов
	router.Use(middleware.Recoverer) // Восстановление после паник

	// Регистрация обработчиков маршрутов
	router.Get("/api/transactions", transaction.GetLast(logger, storage))
	router.Get("/api/wallet/{address}/balance", wallet.GetBalance(logger, storage))
	router.Post("/api/send", transaction.Send(logger, storage))

	// Канал для обработки сигналов завершения
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Конфигурация HTTP-сервера
	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		if err = server.ListenAndServe(); err != nil {
			logger.Error("failed to serve server")
		}
	}()

	logger.Info("server started")

	// Ожидание сигнала завершения
	<-done
	logger.Info("stopping server")
}

// setupLogger инициализирует логгер в зависимости от окружения
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		// Текстовый логгер для разработки
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		// JSON-логгер с debug уровнем
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		// JSON-логгер для продакшна
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
