// Package api имплементация grpc сервера
package api

import (
	"context"
	"errors"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/service"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"

	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
)

// PackageAPI имплементация grpc сервера
type PackageAPI struct {
	pb.UnimplementedLogisticPackageApiServiceServer
	packageService *service.PackageService
}

// CreateV1 grpc ручка создания пакета
func (p *PackageAPI) CreateV1(ctx context.Context, req *pb.CreateRequestV1) (*pb.CreateResponseV1, error) {

	log := slog.With("func", "api.CreateV1")

	if span := trace.SpanContextFromContext(ctx); span.IsSampled() { // вытягиваем span из контекста и пробрасываем в лог
		log = log.With("trace_id", span.TraceID().String())
	}

	if err := req.Validate(); err != nil {
		log.Error("CreateV1 - invalid argument", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pkg := model.Package{}
	pkg.FromProto(req.Value)

	id, err1 := p.packageService.Create(ctx, pkg)
	if err1 != nil {
		log.Error("CreateV1 - failed", slog.String("error", err1.Error()))
		return nil, status.Error(codes.Internal, err1.Error())
	}

	log.Debug("CreateV1 - created", slog.Uint64("id", *id))

	return &pb.CreateResponseV1{PackageId: *id},
		status.New(codes.OK, "").Err()
}

// DeleteV1 grpc ручка удаления пакета
func (p *PackageAPI) DeleteV1(ctx context.Context, req *pb.DeleteV1Request) (*pb.DeleteV1Response, error) {

	log := slog.With("func", "api.DeleteV1")

	if span := trace.SpanContextFromContext(ctx); span.IsSampled() { // вытягиваем span из контекста и пробрасываем в лог
		log = log.With("trace_id", span.TraceID().String())
	}

	if err := req.Validate(); err != nil {
		log.Error("DeleteV1 - invalid argument", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err1 := p.packageService.Delete(ctx, req.PackageId)
	if err1 != nil {
		log.Error("DeleteV1 - failed", slog.String("error", err1.Error()))

		if errors.Is(err1, general.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err1.Error())
		}
		return nil, status.Error(codes.Internal, err1.Error())
	}

	log.Debug("DeleteV1 - deleted", slog.Uint64("id", req.PackageId))
	return &pb.DeleteV1Response{},
		status.New(codes.OK, "").Err()

}

// GetV1 grpc ручка получения пакета
func (p *PackageAPI) GetV1(ctx context.Context, req *pb.GetV1Request) (*pb.GetV1Response, error) {

	log := slog.With("func", "api.GetV1")

	if span := trace.SpanContextFromContext(ctx); span.IsSampled() { // вытягиваем span из контекста и пробрасываем в лог
		log = log.With("trace_id", span.TraceID().String())
	}

	if err := req.Validate(); err != nil {
		log.Error("GetV1 - invalid argument", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pkg, err1 := p.packageService.Get(ctx, req.PackageId)
	if err1 != nil {
		if errors.Is(err1, general.ErrNotFound) {
			log.Debug("not found", slog.Uint64("id", req.PackageId))
			return nil, status.Error(codes.NotFound, "")
		}
		log.Error("failed", slog.String("error", err1.Error()))
		return nil, status.Error(codes.Internal, err1.Error())
	}

	return &pb.GetV1Response{
			Value: pkg.ToProto(),
		},
		status.New(codes.OK, "").Err()
}

// ListV1 grpc ручка получения списка пакетов
func (p *PackageAPI) ListV1(ctx context.Context, req *pb.ListV1Request) (*pb.ListV1Response, error) {

	log := slog.With("func", "api.ListV1")

	if span := trace.SpanContextFromContext(ctx); span.IsSampled() { // вытягиваем span из контекста и пробрасываем в лог
		log = log.With("trace_id", span.TraceID().String())
	}

	if err := req.Validate(); err != nil {
		log.Error("ListV1 - invalid argument", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	packages, err1 := p.packageService.List(ctx, req.Offset, req.Limit)
	if err1 != nil {
		log.Error("ListV1 - failed", slog.String("error", err1.Error()))
		return nil, status.Error(codes.Internal, err1.Error())
	}

	if len(packages) == 0 {
		log.Debug("ListV1 - empty")
		return &pb.ListV1Response{}, status.Error(codes.NotFound, "")
	}

	log.Debug("ListV1 - found", slog.Uint64("count", uint64(len(packages))))

	resp := make([]*pb.Package, len(packages))
	for i := 0; i < len(packages); i++ {
		resp[i] = packages[i].ToProto()
	}

	return &pb.ListV1Response{
			Packages: resp,
		},
		status.New(codes.OK, "").Err()
}

// UpdateV1 grpc ручка изменения пакета
func (p *PackageAPI) UpdateV1(ctx context.Context, req *pb.UpdateV1Request) (*pb.UpdateV1Response, error) {

	log := slog.With("func", "api.UpdateV1")

	if span := trace.SpanContextFromContext(ctx); span.IsSampled() { // вытягиваем span из контекста и пробрасываем в лог
		log = log.With("trace_id", span.TraceID().String())
	}

	if err := req.Validate(); err != nil {
		log.Error("UpdateV1 - invalid argument", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pkg := model.Package{}
	pkg.FromProto(req.Value)

	err1 := p.packageService.Update(ctx, pkg)
	if err1 != nil {
		if errors.Is(err1, general.ErrNotFound) {
			log.Debug("package not found", slog.Uint64("id", pkg.ID))
			return nil, status.Error(codes.NotFound, "")
		}
		log.Error("failed", slog.String("error", err1.Error()))
		return nil, status.Error(codes.Internal, err1.Error())
	}

	log.Debug("UpdateV1 - updated", slog.Any("package", pkg))

	return &pb.UpdateV1Response{},
		status.New(codes.OK, "").Err()
}

// NewPackageAPI returns api of logistic-package-api service
func NewPackageAPI(packageService *service.PackageService) *PackageAPI {
	return &PackageAPI{
		packageService: packageService,
	}
}
