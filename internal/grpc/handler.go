package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/approachi/grinex-rates/gen/rates"
	"github.com/approachi/grinex-rates/internal/service"
)

type Handler struct {
	pb.UnimplementedRatesServiceServer
	svc service.RatesService
}

func NewHandler(svc service.RatesService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	rate, err := h.svc.GetRates(ctx, req.Method, int(req.PositionFrom), int(req.PositionTo))
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pb.GetRatesResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: timestamppb.New(rate.CreatedAt),
	}, nil
}

func (h *Handler) HealthCheck(_ context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Ok: true}, nil
}

