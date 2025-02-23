package transaction

import (
	"errors"
	"github.com/go-chi/render"
	"infotecsTest/internal/lib/api/response"
	"infotecsTest/internal/lib/logger/sl"
	"infotecsTest/internal/models"
	"infotecsTest/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.3 --name=TransactionsReceiver --dir=. --output=./mocks --filename=mock_TransactionsReceiver
type TransactionsReceiver interface {
	GetNTransactions(N int) ([]models.Transaction, error)
}

func GetLast(log *slog.Logger, receiver TransactionsReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.transaction.GetLast"

		log := log.With("op", op)

		var query = r.URL.Query().Get("count")
		var N, err = strconv.Atoi(query)
		if err != nil {
			log.Error("unable to convert count to number", sl.Err(err))
			render.JSON(w, r, response.Error(r, "Некорректное значение count", http.StatusBadRequest))
			return
		}

		txs, err := receiver.GetNTransactions(N)
		if err != nil {
			if errors.Is(err, storage.ErrInvalidRequest) {
				log.Error("invalid request", sl.Err(err))
				render.JSON(w, r, response.Error(r, err.Error(), http.StatusBadRequest))
				return
			}
			log.Error("unable to get transactions", sl.Err(err))
			render.JSON(w, r, response.Error(r, "Внутренняя ошибка", http.StatusInternalServerError))
			return
		}

		render.JSON(w, r, response.Success(txs))
	}
}
