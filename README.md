# Payment System API

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Go-Chi](https://img.shields.io/badge/Router-Go--Chi-6C31D5)](https://go-chi.io/)
[![Testify](https://img.shields.io/badge/Testing-Testify-2CA5E0)](https://github.com/stretchr/testify)

Веб-приложение для управления платежами с REST API, реализованное на Go.

## 🌟 Особенности
- Разработано с применением https://go.dev/doc/effective_go
- Чистая архитектура на основе https://github.com/golang-standards/project-layout
- Детальное логирование через slog https://pkg.go.dev/log/slog
- Полное покрытие тестами обработчиков (handlers)
- Мокирование зависимостей через mockery https://github.com/vektra/mockery
- Обработка ошибок с понятными кодами ответов

## 🛠 Технологии
**Базовый стек:**
- **Go 1.23** - ядро системы
- **Go-Chi** - роутинг и middleware
- **Slog** - структурированное логирование

**Тестирование:**
- **Testify** - assertions и моки
- **Mockery** - генерация мок-интерфейсов
- **httptest** - тестирование HTTP-ендпоинтов

**Инфраструктура:**
- SQLite - хранение данных (легковесная БД)

## 📚 Документация API

### Основные эндпоинты:
| Метод | Путь | Описание |
|-------|------|-----------|
| `GET` | `/api/wallet/{address}/balance` | Получение баланса кошелька |
| `GET` | `api/transactions?count=n` |Получение последних n транзакций|
| `POST` | `/api/send` | Создание новой транзакции |

Пример GET запросa /api/transactions?count=3

Ответ приложения
```JSON
{
    "status": "OK",
    "code": 200,
    "data": [
        {
            "from": "449c9ef0-d2e9-4ec3-97b0-20919d8fac26",
            "to": "e3675892-1de2-4718-8c5a-847a10dd103c",
            "amount": 3.5,
            "time": "05:35:41 24-02-2025"
        },
        {
            "from": "943fc531-2479-4f20-a695-21c5f44191fb",
            "to": "5c1d8064-7f48-4664-9b35-01af789e0179",
            "amount": 52,
            "time": "05:35:13 24-02-2025"
        },
        {
            "from": "53e45c72-53a3-4688-9140-d002f2dd41d3",
            "to": "b4cd8e8f-c2fa-433c-8b07-d2ebfe61468a",
            "amount": 5,
            "time": "05:31:38 24-02-2025"
        }
    ]
}
```
## 🚀 Запуск
Клонировать репозиторий:

```bash
git clone https://github.com/tkchh/payment-system.git
```

Запустить сервер:

```bash
go run cmd/payment-system/main.go
```
