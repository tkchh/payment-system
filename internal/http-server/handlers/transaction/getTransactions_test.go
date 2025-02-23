package transaction_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"infotecsTest/internal/http-server/handlers/transaction"
	"infotecsTest/internal/http-server/handlers/transaction/mocks"
	"infotecsTest/internal/lib/api/response"
	"infotecsTest/internal/models"
	"infotecsTest/internal/storage"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLastHandler(t *testing.T) {
	testLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cases := []struct {
		name         string
		countParam   string
		expectedCode int
		expectedResp response.Response
		mockSetup    func(receiver *mocks.TransactionsReceiver)
	}{
		{
			name:         "Успешный запрос",
			countParam:   "5",
			expectedCode: http.StatusOK,
			expectedResp: response.Response{
				Status: response.StatusOK,
				Data:   []models.Transaction{},
			},
			mockSetup: func(m *mocks.TransactionsReceiver) {
				m.On("GetNTransactions", 5).Return([]models.Transaction{}, nil).Once()
			},
		},
		{
			name:         "Нулевое значение",
			countParam:   "0",
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Error:  storage.ErrInvalidRequest.Error(),
			},
			mockSetup: func(m *mocks.TransactionsReceiver) {
				m.On("GetNTransactions", 0).Return(nil, storage.ErrInvalidRequest).Once()
			},
		},
		{
			name:         "Отрицательное значение",
			countParam:   "-52",
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Error:  storage.ErrInvalidRequest.Error(),
			},
			mockSetup: func(m *mocks.TransactionsReceiver) {
				m.On("GetNTransactions", -52).Return(nil, storage.ErrInvalidRequest).Once()
			},
		},
		{
			name:         "Параметр не число",
			countParam:   "infotecs))",
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Error:  "Некорректное значение count",
			},
		},
		{
			name:         "Внутренняя ошибка",
			countParam:   "5",
			expectedCode: http.StatusInternalServerError,
			expectedResp: response.Response{
				Status: response.StatusError,
				Error:  "Внутренняя ошибка",
			},
			mockSetup: func(m *mocks.TransactionsReceiver) {
				m.On("GetNTransactions", 5).Return(nil, errors.New("unexpected error")).Once()
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockTransactionReceiver := mocks.NewTransactionsReceiver(t)

			if tc.mockSetup != nil {
				tc.mockSetup(mockTransactionReceiver)
			}

			handler := transaction.GetLast(testLogger, mockTransactionReceiver)

			req, err := http.NewRequest(
				http.MethodGet,
				"api/transactions?count="+tc.countParam,
				nil,
			)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

			var resp response.Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResp.Status, resp.Status)
			require.Equal(t, tc.expectedResp.Error, resp.Error)

			if tc.expectedCode == http.StatusOK {
				jsonData, err := json.Marshal(resp.Data)
				require.NoError(t, err)

				var transactions []models.Transaction
				err = json.Unmarshal(jsonData, &transactions)
				require.NoError(t, err)
			}

			if tc.mockSetup != nil {
				mockTransactionReceiver.AssertExpectations(t)
			}
		})
	}
}
