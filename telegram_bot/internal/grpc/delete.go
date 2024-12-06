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

// Delete вызывает gRPC функцию DeleteV1
func (client *Client) Delete(ctx context.Context, id uint64) error {

	log := slog.With("func", "GrpcClient.Delete")

	_, err := client.send.DeleteV1(
		ctx,
		&pb.DeleteV1Request{
			PackageId: id,
		})

	if err != nil {
		log.Error("fail to delete package", slog.String("error", err.Error()))

		if status.Code(err) == codes.NotFound {
			return model.ErrNotFound
		}
		return fmt.Errorf("fail to delete package error: %s", status.Code(err).String())
	}
	return nil
}
