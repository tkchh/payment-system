package storage

import "errors"

var (
	ErrWalletNotFound    = errors.New("Кошелек не найден")
	ErrInsufficientFunds = errors.New("Недостаточно средств")
	ErrIncorrectAmount   = errors.New("Сумма перевода должна быть больше нуля")
	ErrInvalidRequest    = errors.New("Count должен быть больше 0")
	ErrAddressesEqual    = errors.New("Адреса одинаковые")
)
