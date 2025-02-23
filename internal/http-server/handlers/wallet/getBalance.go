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

//go:generate go run github.com/vektra/mockery/v2@v2.52.3 --name=BalanceReceiver --dir=. --output=./mocks --filename=mock_BalanceReceiver
type BalanceReceiver interface {
	GetWalletBalance(address string) (models.Wallet, error)
}

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
