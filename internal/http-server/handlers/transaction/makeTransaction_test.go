package transaction_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"infotecsTest/internal/http-server/handlers/transaction"
	"infotecsTest/internal/http-server/handlers/transaction/mocks"
	"infotecsTest/internal/lib/api/response"
	"infotecsTest/internal/storage"
	_ "infotecsTest/internal/storage"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendHandler(t *testing.T) {
	testLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cases := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedResp response.Response
		mockSetup    func(*mocks.TransactionMaker)
	}{
		{
			name: "Успешный перевод",
			requestBody: `{
				"from": "addr1",
				"to": "addr2",
				"amount": 100.0
			}`,
			expectedCode: http.StatusOK,
			expectedResp: response.Response{
				Status: response.StatusOK,
				Code:   http.StatusOK,
				Data:   "Платеж прошел успешно",
			},
			mockSetup: func(m *mocks.TransactionMaker) {
				m.On("AddTransaction", "addr1", "addr2", 100.0).
					Return(nil).
					Once()
			},
		},
		{
			name:         "Пустое тело запроса",
			requestBody:  ``,
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusBadRequest,
				Error:  "Требуется JSON-объект",
			},
		},
		{
			name: "Невалидный JSON",
			requestBody: `{
				"from": "addr1",
				"to": "addr2",
				"amount": "сто"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusBadRequest,
				Error:  "Некорректный JSON-объект",
			},
		},
		{
			name: "Пустой отправитель",
			requestBody: `{
				"from": "",
				"to": "addr2",
				"amount": 100.0
			}`,
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusBadRequest,
				Error:  "Кошелек не найден",
			},
			mockSetup: func(m *mocks.TransactionMaker) {
				m.On("AddTransaction", "", "addr2", 100.0).
					Return(storage.ErrWalletNotFound).
					Once()
			},
		},
		{
			name: "Пустой получатель",
			requestBody: `{
				"from": "addr1",
				"to": "",
				"amount": 100.0
			}`,
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusBadRequest,
				Error:  "Кошелек не найден",
			},
			mockSetup: func(m *mocks.TransactionMaker) {
				m.On("AddTransaction", "addr1", "", 100.0).
					Return(storage.ErrWalletNotFound).
					Once()
			},
		},
		{
			name: "Отрицательная сумма",
			requestBody: `{
				"from": "addr1",
				"to": "addr2",
				"amount": -100.0
			}`,
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusBadRequest,
				Error:  "Сумма перевода должна быть больше нуля",
			},
			mockSetup: func(m *mocks.TransactionMaker) {
				m.On("AddTransaction", "addr1", "addr2", -100.0).
					Return(storage.ErrIncorrectAmount).
					Once()
			},
		},
		{
			name: "Одинаковые адреса",
			requestBody: `{
				"from": "addr1",
				"to": "addr1",
				"amount": 100.0
			}`,
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusBadRequest,
				Error:  "Адреса одинаковые",
			},
			mockSetup: func(m *mocks.TransactionMaker) {
				m.On("AddTransaction", "addr1", "addr1", 100.0).
					Return(storage.ErrAddressesEqual).
					Once()
			},
		},
		{
			name: "Недостаточно средств",
			requestBody: `{
				"from": "addr1",
				"to": "addr2",
				"amount": 100.0
			}`,
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusBadRequest,
				Error:  "Недостаточно средств",
			},
			mockSetup: func(m *mocks.TransactionMaker) {
				m.On("AddTransaction", "addr1", "addr2", 100.0).
					Return(storage.ErrInsufficientFunds).
					Once()
			},
		},
		{
			name: "Внутренняя ошибка",
			requestBody: `{
				"from": "addr1",
				"to": "addr2",
				"amount": 100.0
			}`,
			expectedCode: http.StatusInternalServerError,
			expectedResp: response.Response{
				Status: response.StatusError,
				Code:   http.StatusInternalServerError,
				Error:  "Внутренняя ошибка",
			},
			mockSetup: func(m *mocks.TransactionMaker) {
				m.On("AddTransaction", "addr1", "addr2", 100.0).
					Return(errors.New("Внутренняя ошибка")).
					Once()
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockTransactionMaker := mocks.NewTransactionMaker(t)

			if tc.mockSetup != nil {
				tc.mockSetup(mockTransactionMaker)
			}

			handler := transaction.Send(testLogger, mockTransactionMaker)

			req, err := http.NewRequest(
				http.MethodPost,
				"/api/transactions",
				bytes.NewReader([]byte(tc.requestBody)),
			)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

			var resp response.Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResp.Status, resp.Status)
			require.Equal(t, tc.expectedResp.Code, resp.Code)
			require.Equal(t, tc.expectedResp.Error, resp.Error)
			require.Equal(t, tc.expectedResp.Data, resp.Data)

			if tc.mockSetup != nil {
				mockTransactionMaker.AssertExpectations(t)
			}
		})
	}
}
