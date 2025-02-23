// Package models содержит структуры данных приложения.
package models

// Wallet представляет данные кошелька пользователя.
type Wallet struct {
	Address string  `json:"address"` // Уникальный адрес кошелька
	Balance float64 `json:"balance"` // Баланс кошелька
}
