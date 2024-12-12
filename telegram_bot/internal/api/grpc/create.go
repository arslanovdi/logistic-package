package grpc

import (
	"context"
	"log/slog"

	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create вызывает gRPC функцию CreateV1
func (c *Client) Create(ctx context.Context, pkg *model.Package) (*uint64, error) {
	log := slog.With("func", "GrpcClient.Create")

	response, err := c.send.CreateV1(
		ctx,
		&pb.CreateRequestV1{
			Value: pkg.ToProto(),
		})
	if err != nil {
		log.Error("fail to create package", slog.String("error", err.Error())) // Logging here, dont return it

		switch status.Code(err) { // Return static error
		case codes.InvalidArgument:
			return nil, general.ErrInvalidArgument
		case codes.DeadlineExceeded:
			return nil, general.ErrDeadline
		case codes.Internal:
			return nil, general.ErrInternal
		default:
			return nil, general.ErrGrpcError
		}
	}

	return &response.PackageId, nil
}
