package grpc

import (
	"context"
	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// Create вызывает gRPC функцию CreateV1
func (client *Client) Create(ctx context.Context, pkg *model.Package) (*uint64, error) {

	log := slog.With("func", "GrpcClient.Create")

	response, err := client.send.CreateV1(
		ctx,
		&pb.CreateRequestV1{
			Value: pkg.ToProto(),
		})

	if err != nil {
		log.Error("fail to create package", slog.String("error", err.Error()))

		switch status.Code(err) {
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
