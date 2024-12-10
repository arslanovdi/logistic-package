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

// Update вызывает gRPC функцию UpdateV1
func (client *Client) Update(ctx context.Context, pkg *model.Package) error {

	log := slog.With("func", "GrpcClient.Update")

	_, err := client.send.UpdateV1(
		ctx,
		&pb.UpdateV1Request{
			Value: pkg.ToProto(),
		})

	if err != nil {
		log.Error("fail to update package", slog.String("error", err.Error()))

		switch status.Code(err) {
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
