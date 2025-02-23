// Package config предоставляет функциональность загрузки и валидации конфигурации приложения.
// Поддерживает чтение из файлов (YAML) и переменных окружения.
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

// Config представляет основную конфигурацию приложения.
// Содержит настройки среды выполнения, хранилища и HTTP-сервера.
type Config struct {
	Env         string               `yaml:"env" env-default:"development"`    // Окружение приложения (dev/prod)
	StoragePath string               `yaml:"storage_path" env-required:"true"` // Путь к файлу хранилища данных
	HTTPServer  `yaml:"http_server"` // Настройки HTTP-сервера
}

// HTTPServer содержит конфигурационные параметры HTTP-сервера.
type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8080"` // Адрес сервера (host:port)
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`             // Таймаут обработки запросов
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`       // Таймаут бездействующих соединений
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`   // Таймаут graceful shutdown
}

// MustLoad загружает конфигурацию из файла и переменных окружения.
// Завершает выполнение приложения с фатальной ошибкой в случае:
// - Не указан путь к конфигурации (CONFIG_PATH)
// - Ошибки чтения/парсинга конфигурационного файла
//
// Возвращает:
//   - *Config: указатель на загруженную конфигурацию
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Переменная окружения CONFIG_PATH не установлена")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Не удалось открыть конфиг файл: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Не удалось прочитать конфиг файл: %s", err)
	}

	return &cfg
}
