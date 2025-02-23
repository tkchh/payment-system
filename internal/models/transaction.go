// Package models содержит структуры данных приложения.
package models

// Transaction описывает денежный перевод между кошельками.
type Transaction struct {
	From   string  `json:"from"`           // Адрес отправителя
	To     string  `json:"to"`             // Адрес получателя
	Amount float64 `json:"amount"`         // Сумма перевода
	Time   string  `json:"time,omitempty"` // Время транзакции (опционально)
}
