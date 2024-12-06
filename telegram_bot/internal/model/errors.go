package model

import "errors"

var (
	// ErrEndOfList конец полученного от сервера списка пакетов
	ErrEndOfList = errors.New("end of list")
	// ErrNotFound пакет не найден
	ErrNotFound = errors.New("not found")
)
