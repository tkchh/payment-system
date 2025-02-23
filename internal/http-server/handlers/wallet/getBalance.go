// Package wallet содержит обработчики HTTP-запросов для работы с кошельками.
package wallet

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"infotecsTest/internal/lib/api/response"
	"infotecsTest/internal/lib/logger/sl"
	"infotecsTest/internal/models"
	"infotecsTest/internal/storage"
	"log/slog"
	"net/http"
)

// BalanceReceiver определяет интерфейс для получения баланса из хранилища.
// Генерирует моки через go:generate.
type BalanceReceiver interface {
	GetWalletBalance(address string) (models.Wallet, error)
}

// GetBalance создает HTTP-обработчик для получения баланса кошелька.
// Извлекает адрес из URL-параметров, обрабатывает ошибки хранилища,
// возвращает баланс в формате JSON или соответствующие HTTP-ошибки.
func GetBalance(log *slog.Logger, receiver BalanceReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.wallet.GetBalance"

		log := log.With("op", op)

		address := chi.URLParam(r, "address")

		wallet, err := receiver.GetWalletBalance(address)
		if err != nil {
			if errors.Is(err, storage.ErrWalletNotFound) {
				log.Error("wallet not found", sl.Err(err))
				render.JSON(w, r, response.Error(r, err.Error(), http.StatusBadRequest))
				return
			}
			log.Error("unable to get balance", sl.Err(err))
			render.JSON(w, r, response.Error(r, "Внутренняя ошибка", http.StatusInternalServerError))
			return
		}

		render.JSON(w, r, response.Success(wallet))
	}
}
