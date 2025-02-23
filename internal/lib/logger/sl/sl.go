// Package sl содержит вспомогательные функции для работы с slog.
package sl

import "log/slog"

// Err создает slog-атрибут для ошибки с ключом "error".
// Упрощает логирование ошибок в структурированном формате.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
