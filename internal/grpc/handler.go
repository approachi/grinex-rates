package grpc

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/approachi/grinex-rates/gen/rates"
	"github.com/approachi/grinex-rates/internal/observability/metrics"
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
	start := time.Now()
	rate, err := h.svc.GetRates(ctx, req.Method, int(req.PositionFrom), int(req.PositionTo))
	metrics.GetRatesDuration.WithLabelValues(req.Method).Observe(time.Since(start).Seconds())

	if err != nil {
		metrics.GetRatesTotal.WithLabelValues(req.Method, "error").Inc()
		return nil, toGRPCError(err)
	}
	metrics.GetRatesTotal.WithLabelValues(req.Method, "ok").Inc()

	return &pb.GetRatesResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: timestamppb.New(rate.CreatedAt),
	}, nil
}

func (h *Handler) HealthCheck(_ context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Ok: true}, nil
}
