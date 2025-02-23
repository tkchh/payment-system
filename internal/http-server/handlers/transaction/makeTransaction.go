package transaction

import (
	"errors"
	"github.com/go-chi/render"
	"infotecsTest/internal/lib/api/response"
	"infotecsTest/internal/lib/logger/sl"
	"infotecsTest/internal/models"
	"infotecsTest/internal/storage"
	"io"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.3 --name=TransactionMaker --dir=. --output=./mocks --filename=mock_TransactionMaker
type TransactionMaker interface {
	AddTransaction(from, to string, amount float64) error
}

func Send(log *slog.Logger, maker TransactionMaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.transaction.Send"

		log := log.With("op", op)

		var req models.Transaction

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, response.Error(r, "Требуется JSON-объект", http.StatusBadRequest))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error(r, "Некорректный JSON-объект", http.StatusBadRequest))
			return
		}

		if err = maker.AddTransaction(req.From, req.To, req.Amount); err != nil {
			log.Error("failed to make transaction", sl.Err(err))
			switch {
			case errors.Is(err, storage.ErrWalletNotFound):
				render.JSON(w, r, response.Error(r, err.Error(), http.StatusBadRequest))
			case errors.Is(err, storage.ErrIncorrectAmount):
				render.JSON(w, r, response.Error(r, err.Error(), http.StatusBadRequest))
			case errors.Is(err, storage.ErrInsufficientFunds):
				render.JSON(w, r, response.Error(r, err.Error(), http.StatusBadRequest))
			case errors.Is(err, storage.ErrAddressesEqual):
				render.JSON(w, r, response.Error(r, err.Error(), http.StatusBadRequest))
			default:
				render.JSON(w, r, response.Error(r, "Внутренняя ошибка", http.StatusInternalServerError))
			}
			return
		}
		render.JSON(w, r, response.Success("Платеж прошел успешно"))
		return
	}
}
