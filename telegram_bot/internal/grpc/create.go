package grpc

import (
	"context"
	"fmt"
	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/model"
)

// Create вызывает gRPC функцию CreateV1
func (client *Client) Create(ctx context.Context, pkg model.Package) (*uint64, error) {

	response, err := client.send.CreateV1(
		ctx,
		&pb.CreateRequestV1{
			Value: pkg.ToProto(),
		})

	if err != nil {
		return nil, fmt.Errorf("grpc.Create: %w", err)
	}

	return &response.PackageId, nil
}
