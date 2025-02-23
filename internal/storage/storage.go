// Package storage определяет ошибки уровня хранилища данных.
package storage

import "errors"

var (
	// ErrWalletNotFound возвращается при отсутствии кошелька в БД.
	ErrWalletNotFound = errors.New("Кошелек не найден")

	// ErrInsufficientFunds возникает при попытке списать сумму больше баланса.
	ErrInsufficientFunds = errors.New("Недостаточно средств")

	// ErrIncorrectAmount указывает на недопустимую сумму перевода (<= 0).
	ErrIncorrectAmount = errors.New("Сумма перевода должна быть больше нуля")

	// ErrInvalidRequest возвращается при некорректных параметрах запроса.
	ErrInvalidRequest = errors.New("Count должен быть больше 0")

	// ErrAddressesEqual возникает при совпадении адресов отправителя и получателя.
	ErrAddressesEqual = errors.New("Адреса одинаковые")
)
