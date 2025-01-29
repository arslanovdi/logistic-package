// Package general - общие ошибки
package general

import "errors"

var (
	// ErrNotFound - запись не найдена
	ErrNotFound = errors.New("not found")
	// ErrNoDeliveredMessage - Завершение работы kafka продюсера без подтверждения доставки сообщения.
	ErrNoDeliveredMessage = errors.New("kafka produce stopped without delivered message")
	// ErrProducerClosed - Попытка отправки сообщения после завершения работы kafka продюсера
	ErrProducerClosed = errors.New("kafka producer closed")
)
