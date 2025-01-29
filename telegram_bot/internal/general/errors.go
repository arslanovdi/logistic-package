// Package general errors
package general

import "errors"

var (
	// ErrEndOfList конец полученного от сервера списка пакетов
	ErrEndOfList = errors.New("end of list")
	// ErrNotFound пакет не найден
	ErrNotFound = errors.New("not found")

	// ErrGrpcError ошибка gRPC
	ErrGrpcError = errors.New("grpc error")

	// ErrDeadline превышен срок действия запроса
	ErrDeadline = errors.New("request deadline exceeded")

	// ErrInvalidArgument некорректные данные
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrInternal внутренняя ошибка
	ErrInternal = errors.New("internal error in logistic-package-api")
)

/*
switch status.Code(err) {
case codes.OK:
case codes.Canceled:
case codes.Unknown:
case codes.InvalidArgument:
case codes.DeadlineExceeded:
case codes.NotFound:
case codes.AlreadyExists:
case codes.PermissionDenied:
case codes.ResourceExhausted:
case codes.FailedPrecondition:
case codes.Aborted:
case codes.OutOfRange:
case codes.Unimplemented:
case codes.Internal:
case codes.Unavailable:
case codes.DataLoss:
case codes.Unauthenticated:
default:
			return general.ErrGrpcError
}
*/
