package grpc

import (
	"context"
	"fmt"
	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// List вызывает gRPC функцию ListV1
func (client *Client) List(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error) {

	log := slog.With("func", "GrpcClient.List")

	response, err1 := client.send.ListV1(
		ctx,
		&pb.ListV1Request{
			Offset: offset,
			Limit:  limit,
		})

	if err1 != nil {
		log.Error("fail to list packages", slog.String("error", err1.Error()))

		if status.Code(err1) == codes.NotFound {
			return nil, model.ErrNotFound
		}

		return nil, fmt.Errorf("fail to list packages error: %s", status.Code(err1).String())
	}

	packages := make([]model.Package, len(response.Packages))
	for i := 0; i < len(response.Packages); i++ {
		packages[i].FromProto(response.Packages[i])
	}

	return packages, nil
}
