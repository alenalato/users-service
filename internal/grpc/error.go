package grpc

import (
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// commonErrorToGRPCError converts a common.Error to a gRPC Status
func commonErrorToGRPCError(err error) error {
	var errCommon common.Error
	switch {
	case errors.As(err, &errCommon):
		switch err.(common.Error).Type() {
		case common.ErrTypeNotFound:
			return status.New(codes.NotFound, err.Error()).Err()
		case common.ErrTypeAlreadyExists:
			return status.New(codes.AlreadyExists, err.Error()).Err()
		case common.ErrTypeInvalidArgument:
			return status.New(codes.InvalidArgument, err.Error()).Err()
		case common.ErrTypeInternal:
			return status.New(codes.Internal, err.Error()).Err()
		default:
			return status.New(codes.Internal, err.Error()).Err()
		}
	default:
		return err
	}
}
