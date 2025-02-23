package wallet_test

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"infotecsTest/internal/http-server/handlers/wallet"
	"infotecsTest/internal/http-server/handlers/wallet/mocks"
	"infotecsTest/internal/lib/api/response"
	"infotecsTest/internal/models"
	"infotecsTest/internal/storage"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBalanceHandler(t *testing.T) {
	testLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cases := []struct {
		name         string
		address      string
		expectedCode int
		expectedResp response.Response
		mockSetup    func(receiver *mocks.BalanceReceiver)
	}{
		{
			name:         "Успешный запрос",
			address:      "addr1",
			expectedCode: http.StatusOK,
			expectedResp: response.Response{
				Status: response.StatusOK,
				Data:   models.Wallet{Address: "addr1", Balance: 100},
			},
			mockSetup: func(m *mocks.BalanceReceiver) {
				m.On("GetWalletBalance", "addr1").Return(models.Wallet{Address: "addr1", Balance: 100}, nil).Once()
			},
		},
		{
			name:         "Кошелек не найден",
			address:      "not_found_address",
			expectedCode: http.StatusBadRequest,
			expectedResp: response.Response{
				Status: response.StatusError,
				Error:  storage.ErrWalletNotFound.Error(),
			},
			mockSetup: func(m *mocks.BalanceReceiver) {
				m.On("GetWalletBalance", "not_found_address").Return(models.Wallet{}, storage.ErrWalletNotFound).Once()
			},
		},
		{
			name:         "Внутренняя ошибка",
			address:      "addr1",
			expectedCode: http.StatusInternalServerError,
			expectedResp: response.Response{
				Status: response.StatusError,
				Error:  "Внутренняя ошибка",
			},
			mockSetup: func(m *mocks.BalanceReceiver) {
				m.On("GetWalletBalance", "addr1").Return(models.Wallet{}, errors.New("unexpected error")).Once()
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockBalanceReceiver := mocks.NewBalanceReceiver(t)

			if tc.mockSetup != nil {
				tc.mockSetup(mockBalanceReceiver)
			}

			handler := wallet.GetBalance(testLogger, mockBalanceReceiver)

			router := chi.NewRouter()
			router.Get("/{address}", handler)

			req, err := http.NewRequest(
				http.MethodGet,
				"/"+tc.address,
				nil,
			)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

			var resp response.Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResp.Status, resp.Status)
			require.Equal(t, tc.expectedResp.Error, resp.Error)

			if tc.expectedCode == http.StatusOK {
				jsonData, err := json.Marshal(resp.Data)
				require.NoError(t, err)

				var w models.Wallet
				err = json.Unmarshal(jsonData, &w)
				require.NoError(t, err)
				require.Equal(t, tc.expectedResp.Data.(models.Wallet), w)
			}

			if tc.mockSetup != nil {
				mockBalanceReceiver.AssertExpectations(t)
			}
		})
	}
}
