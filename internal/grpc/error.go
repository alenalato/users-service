package grpc

import (
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func commonErrorToGRPCError(err error) error {
	var commonErr common.Error
	switch {
	case errors.As(err, &commonErr):
		switch err.(common.Error).GetType() {
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
