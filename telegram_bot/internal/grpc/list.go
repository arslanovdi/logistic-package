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

// List вызывает gRPC функцию ListV1
func (client *Client) List(ctx context.Context, offset, limit uint64) ([]model.Package, error) {

	log := slog.With("func", "GrpcClient.List")

	response, err := client.send.ListV1(
		ctx,
		&pb.ListV1Request{
			Offset: offset,
			Limit:  limit,
		})

	if err != nil {
		log.Error("fail to list packages", slog.String("error", err.Error()))

		switch status.Code(err) {
		case codes.InvalidArgument:
			return nil, general.ErrInvalidArgument
		case codes.DeadlineExceeded:
			return nil, general.ErrDeadline
		case codes.NotFound:
			return nil, general.ErrNotFound
		case codes.Internal:
			return nil, general.ErrInternal
		default:
			return nil, general.ErrGrpcError
		}
	}

	packages := make([]model.Package, len(response.Packages))
	for i := 0; i < len(response.Packages); i++ {
		packages[i].FromProto(response.Packages[i])
	}

	return packages, nil
}
