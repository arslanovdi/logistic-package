package model

import "errors"

var (
	// ErrNotFound - запись не найдена
	ErrNotFound           = errors.New("not found")
	ErrNoDeliveredMessage = errors.New("kafka produce stopped without delivered message")
	ErrProducerClosed     = errors.New("kafka producer closed")
)
