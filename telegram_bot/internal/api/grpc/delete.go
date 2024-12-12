package grpc

import (
	"context"
	"log/slog"

	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Delete вызывает gRPC функцию DeleteV1
func (c *Client) Delete(ctx context.Context, id uint64) error {
	log := slog.With("func", "GrpcClient.Delete")

	_, err := c.send.DeleteV1(
		ctx,
		&pb.DeleteV1Request{
			PackageId: id,
		})
	if err != nil {
		log.Error("fail to delete package", slog.String("error", err.Error())) // Logging here, dont return it

		switch status.Code(err) { // Return static error
		case codes.InvalidArgument:
			return general.ErrInvalidArgument
		case codes.DeadlineExceeded:
			return general.ErrDeadline
		case codes.NotFound:
			return general.ErrNotFound
		case codes.Internal:
			return general.ErrInternal
		default:
			return general.ErrGrpcError
		}
	}
	return nil
}
