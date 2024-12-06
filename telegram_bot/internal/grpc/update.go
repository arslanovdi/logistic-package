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

// Update вызывает gRPC функцию UpdateV1
func (client *Client) Update(ctx context.Context, pkg model.Package) error {

	log := slog.With("func", "GrpcClient.Update")

	_, err := client.send.UpdateV1(
		ctx,
		&pb.UpdateV1Request{
			Value: pkg.ToProto(),
		})

	if err != nil {
		log.Error("fail to update package", slog.String("error", err.Error()))

		if status.Code(err) == codes.NotFound {
			return model.ErrNotFound
		}
		return fmt.Errorf("fail to update package error: %s", status.Code(err).String())
	}

	return nil
}
