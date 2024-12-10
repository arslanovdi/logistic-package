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

// Get вызывает gRPC функцию GetV1
func (client *Client) Get(ctx context.Context, id uint64) (*model.Package, error) {

	log := slog.With("func", "GrpcClient.Get")

	response, err := client.send.GetV1(
		ctx,
		&pb.GetV1Request{
			PackageId: id,
		})

	if err != nil {
		log.Error("fail to get package", slog.String("error", err.Error()))

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

	pkg := model.Package{}
	pkg.FromProto(response.Value)

	return &pkg, nil
}
