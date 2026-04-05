package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/approachi/grinex-rates/internal/domain"
)

func toGRPCError(err error) error {
	switch err {
	case domain.ErrUnknownMethod, domain.ErrEmptyEntries, domain.ErrOutOfRange, domain.ErrInvalidRange:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
